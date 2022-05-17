package httpclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var userAgents = []string{

	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:56.0) Gecko/20100101 Firefox/56.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Safari/604.1.38",
}

var mu = new(sync.Mutex)

//Build Default Headers sent with every request
func BuildDefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":      "goshod by @jayateerthag", //RandomUA(),
		"Accept":          "text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8",
		"Accept-Language": "en-US,en;q=0.5",
		//	"Accept-Encoding": "gzip,deflate",
		"DNT": "1",
		//"Connection": "close",
	}
}

//Get request with params
func GetRequestP(url string, params map[string]string, headers map[string]string, timeoutP int) (*http.Response, error) {
	timeout := timeoutP
	if timeout <= 0 {
		timeout = 5
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	var client = new(http.Client)

	client = &http.Client{Transport: transport, Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true

	//params adding and encoding
	q := req.URL.Query()
	for name, value := range params {
		q.Add(name, value)
	}
	req.URL.RawQuery = q.Encode()

	//Set Default Headers with Ramdom UA
	tmpHeaders := BuildDefaultHeaders()
	for key, value := range tmpHeaders {
		req.Header.Set(key, value)
	}

	mu.Lock()
	if headers != nil {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	mu.Unlock()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
func GetRequest(uri string, headers map[string]string, timeoutP int, proxy string) (*http.Response, error) {
	h := CloneHeaders(headers)
	pu, err := url.Parse(proxy)

	if proxy != "" {
		if err != nil {
			fmt.Println("Invalid proxy...")
			return nil, nil
		}
	}
	timeout := timeoutP
	if timeout <= 0 {
		timeout = 5
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	var client = new(http.Client)
	if proxy != "" {
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyURL(pu),
		}
	}
	client = &http.Client{Transport: transport, Timeout: time.Duration(timeout) * time.Second}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Close = true
	//Set Default Headers with Ramdom UA
	tmpHeaders := BuildDefaultHeaders()
	for key, value := range tmpHeaders {
		req.Header.Set(key, value)
	}
	mu.Lock()
	if h != nil {
		for key, value := range h {
			req.Header.Set(key, value)
		}
	}
	mu.Unlock()
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

//Post Request supports normal and json body too
func PostRequest(uri string, data string, headers map[string]string, timeoutP int, proxy string) (*http.Response, error) {
	h := CloneHeaders(headers)
	pu, err := url.Parse(proxy)

	if proxy != "" {
		if err != nil {
			fmt.Println("Invalid proxy...")
			return nil, nil
		}
	}
	var req *http.Request
	var resp *http.Response

	timeout := timeoutP
	if timeout <= 0 {
		timeout = 5
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	var client = new(http.Client)
	if proxy != "" {
		transport = &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyURL(pu),
		}
	}
	client = &http.Client{Transport: transport, Timeout: time.Duration(timeout) * time.Second}
	req, err = http.NewRequest("POST", uri, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Close = true

	req.Header.Set("Content-Type", "application/x-www-form-uriencoded")

	//Set Default Headers with Ramdom UA
	tmpHeaders := BuildDefaultHeaders()
	for key, value := range tmpHeaders {
		req.Header.Set(key, value)
	}

	mu.Lock()
	if h != nil {
		for key, value := range h {
			req.Header.Set(key, value)
		}
	}
	mu.Unlock()
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func RandomUA() string {
	rand.Seed(time.Now().Unix())
	choice := rand.Intn(len(userAgents))
	return userAgents[choice]
}

//default json escapes html characters such as <> so custom marshal needed
func JsonMarshal(params interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&params)
	if err != nil {
		log.Println(err)
	}
	return buffer.Bytes(), err
}

//default json escapes html characters such as <> so custom marshal needed
func JsonMarshalIndent(params interface{}, prefix string, indent string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	enc.SetIndent(prefix, indent)
	err := enc.Encode(&params)
	if err != nil {
		log.Println(err)
	}
	return buffer.Bytes(), err
}

func CloneHeaders(m map[string]string) map[string]string {
	cp := make(map[string]string)
	if m == nil {
		return nil
	}
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
