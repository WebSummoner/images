package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	listen      = ":4444"
	target      = "http://localhost:5555"
	waitTimeout = 30 * time.Second
	gracePeriod = 30 * time.Second
	browserName = "safari"
	// baked into the image as WEBKIT_VERSION
	browserVersion = os.Getenv("WEBKIT_VERSION")
	// WebKit plays nothing without a user gesture
	enableAudio = os.Getenv("ENABLE_AUDIO") == "true"
)

// autoclick grants the page the gesture WebKit demands before playback.
func autoclick(sessionId string) {
	if !enableAudio || sessionId == "" {
		return
	}
	go func() {
		time.Sleep(500 * time.Millisecond)
		body := `{"actions":[{"type":"pointer","id":"prism","parameters":{"pointerType":"mouse"},"actions":[{"type":"pointerMove","x":1,"y":1},{"type":"pointerDown","button":0},{"type":"pointerUp","button":0}]}]}`
		req, err := http.NewRequest(http.MethodPost, target+path.Join("/session", sessionId, "actions"), strings.NewReader(body))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}()
}

func wait(ctx context.Context, target string) (*url.URL, error) {
	for {
		r, err := http.NewRequest(http.MethodHead, target, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("new %s request to %s: %v", http.MethodHead, target, err)
		}
		resp, err := http.DefaultClient.Do(r.WithContext(ctx))
		if resp != nil {
			resp.Body.Close()
		}
		if err != nil {
			if err, ok := err.(*url.Error); ok {
				switch err.Err {
				case context.Canceled, context.DeadlineExceeded:
					return nil, err
				default:
					<-time.After(100 * time.Millisecond)
					continue
				}
			}
		}
		return r.URL, nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	e := make(chan error)
	go func() {
		stop := make(chan os.Signal)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		select {
		case err := <-e:
			log.Fatalf("server: %v", err)
		case <-stop:
			cancel()
		}
	}()
	waitCtx, waitCancel := context.WithTimeout(ctx, waitTimeout)
	defer waitCancel()
	u, err := wait(waitCtx, target)
	if err != nil {
		log.Fatal(fmt.Errorf("wait target: %v", err))
	}
	server := &http.Server{
		Addr: listen,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var value map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&value)
			if err != nil && err != io.EOF {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if r.Method == http.MethodPost {
				fragments := strings.Split(strings.TrimPrefix(r.URL.Path, "/wd/hub"), "/")
				if len(fragments) == 4 && fragments[1] == "session" && fragments[3] == "url" {
					autoclick(fragments[2])
				}
			}
			if err == nil {
				if _, ok := value["desiredCapabilities"]; ok {
					delete(value, "desiredCapabilities")
				}
				if o, ok := value["capabilities"]; ok {
					if w3cCapabilities, ok := o.(map[string]interface{}); ok {
						for _, match := range []string{"alwaysMatch", "firstMatch"} {
							delete(w3cCapabilities, match)
						}
					}
				}
				body, err := json.Marshal(value)
				if err != nil {
					log.Printf("[ERROR] marshalling capabilities: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				r.Body = io.NopCloser(bytes.NewReader(body))
				r.ContentLength = int64(len(body))
			}
			(&httputil.ReverseProxy{
				Director: func(r *http.Request) {
					r.URL.Scheme, r.URL.Host, r.URL.Path = u.Scheme, u.Host, path.Join(u.Path, r.URL.Path)
				},
				ModifyResponse: func(resp *http.Response) error {
					if resp.StatusCode != http.StatusOK {
						return nil
					}
					var values map[string]interface{}
					defer resp.Body.Close()
					err := json.NewDecoder(resp.Body).Decode(&values)
					if err != nil {
						return fmt.Errorf("decode json response: %v", err)
					}
					if o, ok := values["value"]; ok {
						if value, ok := o.(map[string]interface{}); ok {
							if capabilities, ok := value["capabilities"]; ok {
								if caps, ok := capabilities.(map[string]interface{}); ok {
									caps["browserName"] = browserName
									caps["browserVersion"] = browserVersion
									delete(caps, "platformName")
								}
							}
							if id, ok := value["sessionId"].(string); ok && id != "" {
								autoclick(id)
							}
						}
					}
					buf, err := json.Marshal(&values)
					if err != nil {
						return fmt.Errorf("encode json response: %v", err)
					}
					resp.Header.Del("Server")
					resp.Header.Del("Content-Length")
					resp.ContentLength = int64(len(buf))
					resp.Body = io.NopCloser(bytes.NewReader(buf))
					return nil
				},
			}).ServeHTTP(w, r)
		}),
	}
	go func() {
		e <- server.ListenAndServe()
	}()
	<-ctx.Done()
	shCtx, shCancel := context.WithTimeout(context.Background(), gracePeriod)
	defer shCancel()
	if err := server.Shutdown(shCtx); err != nil {
		log.Fatalf("graceful shutdown: %v]", err)
	}
}
