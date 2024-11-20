package proxy

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Proxy struct {
	Context context.Context
	URL     string `json:"url"`
	Port    string `json:"port"`
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
	uri := fmt.Sprintf("%s%s", p.URL, r.URL.Path)
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

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Error copying response", http.StatusInternalServerError)
		return
	}
}
