package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/balajiss36/cache-proxy/cache"
)

type Proxy struct {
	Context context.Context
	URL     string `json:"url"`
	Port    string `json:"port"`
	// this acts as a interface between cache and proxy to run the cache functions in proxy server
	Cache cache.CacheService
}

var customTransport = http.DefaultTransport

func (p *Proxy) StartServer() error {
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", p.Port),
		Handler: http.HandlerFunc(p.handleRequest),
	}

	log.Printf("Starting proxy server on %s", p.Port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting proxy server: ", err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	if err := server.Shutdown(p.Context); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	return nil
}

func (p *Proxy) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	uri := fmt.Sprintf("%s%s", p.URL, r.URL.Path)
	log.Println("URI: & path ", uri, r.URL.Path)
	val, err := p.Cache.Get(path)
	if err != nil {
		log.Printf("Error while getting cache for key %s err %s", path, err)
	}
	if val != nil {
		log.Println("Cache hit")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = io.Copy(w, bytes.NewReader(val))
		if err != nil {
			http.Error(w, "Error copying response", http.StatusInternalServerError)
			return
		}
		return
	}

	proxyReq, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	log.Println("Proxying request to ", uri)
	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, bytes.NewReader(res))
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	err = p.Cache.Set(path, res)
	log.Println("Setting cache for key", path)
	if err != nil {
		log.Println("Error while setting cache")
	}
}
