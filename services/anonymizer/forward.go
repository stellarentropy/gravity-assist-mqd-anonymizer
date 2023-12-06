package anonymizer

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	config "github.com/stellarentropy/gravity-assist-mqd-anonymizer/config/anonymizer"
)

var client *http.Client

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	client = &http.Client{
		Timeout:   60 * time.Second,
		Transport: t,
	}
}

func ForwardRequest(w http.ResponseWriter, r *http.Request) {
	defer func() { _ = r.Body.Close() }()

	agentUrl := url.URL{}

	agentUrl.Scheme = config.Anonymizer.AgentSchema
	agentUrl.Host = net.JoinHostPort(config.Anonymizer.AgentAddress, fmt.Sprintf("%d", config.Anonymizer.AgentPort))

	// Parse the input endpoint into a URL structure
	agentUrl.Path = r.URL.Path

	// Create a new HTTP request with the given method, URL, and body
	req, err := http.NewRequest(r.Method, agentUrl.String(), AnonymizePayload(r.Body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the headers for the request
	for k, v := range r.Header {
		req.Header.Set(k, v[0])
	}
	// Attach the context to the request
	req = req.WithContext(r.Context())

	// Send the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}

	w.WriteHeader(resp.StatusCode)

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
