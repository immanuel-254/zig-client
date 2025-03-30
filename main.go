package main

import (
	"C"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
}

type Client struct {
	Url    string
	Key    string
	Route  string
	Method string
	Input  map[string]string
	Output map[string]any
}

func (c *Client) Send() {
	jsonData, err := json.Marshal(c.Input)
	if err != nil {
		c.Output["error"] = err.Error()
		return
	}

	if len(c.Route) > 0 {
		c.Route = strings.TrimPrefix(c.Route, "/")
	}

	req, err := http.NewRequest(c.Method, fmt.Sprintf("%s/%s", c.Url, c.Route), bytes.NewBuffer(jsonData))
	if err != nil {
		c.Output["error"] = err.Error()
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Key != "" {
		req.Header.Set("auth", c.Key)
	}

	resp, err := client.Do(req)
	if err != nil {
		c.Output["error"] = err.Error()
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Output["error"] = err.Error()
		return
	}

	c.Output["status"] = resp.StatusCode
	c.Output["body"] = string(body)
}

//export Request
func Request(url, key, route, method, input *C.char) *C.char {
	var request_input map[string]string

	if C.GoString(input) != "" {
		err := json.Unmarshal([]byte(C.GoString(input)), &request_input)
		if err != nil {
			return C.CString(err.Error())
		}
	}

	client := Client{
		Url:    C.GoString(url),
		Key:    C.GoString(key),
		Route:  C.GoString(route),
		Method: C.GoString(method),
		Input:  request_input,
		Output: make(map[string]any),
	}
	client.Send()

	// Serialize the Output map to JSON
	jsonOutput, err := json.Marshal(client.Output)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "failed to serialize output: %s"}`, err.Error()))
	}

	// Return the JSON string as a C-compatible string
	return C.CString(string(jsonOutput))
}

func main() {}

// go build -trimpath -ldflags="-s -w" -buildmode=c-archive -o include/client.a
