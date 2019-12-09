package common

import (
	"crypto/tls"
	"encoding/json"

	"io"
	"log"
	"net/http"
	"time"
)

func HttpClient(method string, url string, body io.Reader, target interface{}) (status int, err error) {
	//tr := &http.Transport{DisableKeepAlives: true}
	//caCertPool := x509.NewCertPool()
	transport := &http.Transport{
		MaxIdleConnsPerHost:   150,
		IdleConnTimeout:       10 * time.Second,
		DisableCompression:    true,
		//MaxIdleConns:          15,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 600 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
		ResponseHeaderTimeout: 60 * time.Second,
		ForceAttemptHTTP2:     true,

		//DialContext:			(&net.Dialer{ Timeout:   30 * time.Second, KeepAlive: 30 * time.Second, DualStack: true, }).DialContext,
	}
	//transport := &http2.Transport{
	//	AllowHTTP: true,
	//	DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
	//		return net.Dial(network, addr)
	//	},
	//}

	client := &http.Client{Transport: transport, Timeout: 20 * time.Second}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println("request error", err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth("admin", "Harbor12345")
	res, err := client.Do(request)
	if err != nil {
		log.Println("Http Request Error:", err.Error())
		return
	}

	defer res.Body.Close()
	return res.StatusCode, json.NewDecoder(res.Body).Decode(target)
}
