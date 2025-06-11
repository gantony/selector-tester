package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"embed"

	"github.com/sirupsen/logrus"
)

//go:embed tool.html tool.js
var staticFiles embed.FS

func main() {

	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	fmt.Printf("Go to http://127.0.0.1:%v\n", *port)

	http.HandleFunc("/proxy", proxy)
	http.HandleFunc("/", serveTool)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

// proxy takes a POST request with a JSON body taht contains a "url" and "payload" fields.
// It then makes a POST request to the specified URL with the provided payload.
// The response from the proxied request is returned to the original client.
func proxy(w http.ResponseWriter, req *http.Request) {
	addCorsHeader(w)
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	type ProxyRequest struct {
		URL     string `json:"url"`
		Payload string `json:"payload"`
		Token   string `json:"token"`
	}

	var proxyReq ProxyRequest
	if err := json.NewDecoder(req.Body).Decode(&proxyReq); err != nil {
		logrus.Errorf("Failed to decode JSON body: %v", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if proxyReq.URL == "" {
		logrus.Error("Missing 'url' field in request")
		http.Error(w, "Missing 'url' field", http.StatusBadRequest)
		return
	}
	if proxyReq.Token == "" {
		logrus.Error("Missing 'token' field in request")
		http.Error(w, "Missing 'token' field", http.StatusBadRequest)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest("POST", proxyReq.URL, strings.NewReader(string(proxyReq.Payload)))
	if err != nil {
		http.Error(w, "Failed to create proxied request: "+err.Error(), http.StatusBadGateway)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	if proxyReq.Token != "" {
		request.Header.Set("Authorization", "Bearer "+proxyReq.Token)
	}

	proxyResp, err := client.Do(request)
	if err != nil {
		http.Error(w, "Failed to proxy request: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer proxyResp.Body.Close()

	// Copy status code
	w.WriteHeader(proxyResp.StatusCode)
	// Copy headers (except for hop-by-hop headers)
	for k, vv := range proxyResp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// Copy body
	_, _ = io.Copy(w, proxyResp.Body)
}

func serveTool(w http.ResponseWriter, req *http.Request) {
	addCorsHeader(w)
	path := req.URL.Path
	if path == "/" || path == "/tool.html" {
		w.Header().Set("Content-Type", "text/html")
		data, err := staticFiles.ReadFile("tool.html")
		if err != nil {
			http.Error(w, "tool.html not found", http.StatusNotFound)
			return
		}
		w.Write(data)
		return
	}
	if path == "/tool.js" {
		w.Header().Set("Content-Type", "application/javascript")
		data, err := staticFiles.ReadFile("tool.js")
		if err != nil {
			http.Error(w, "tool.js not found", http.StatusNotFound)
			return
		}
		w.Write(data)
		return
	}
	http.NotFound(w, req)
}

func addCorsHeader(res http.ResponseWriter) {
	headers := res.Header()
	headers.Add("Access-Control-Allow-Origin", "*")
	headers.Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
	headers.Add("Access-Control-Allow-Methods", "GET, POST,OPTIONS")
}
