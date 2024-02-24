package php

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func WebRequestToRequest(req *WebRequest) (*http.Request, error) {
	decodedBody := []byte(req.Body)

	path := req.Url

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	serverAddress := req.Schema + req.Headers["Host"]
	path = serverAddress + path

	httpRequest, err := http.NewRequest(
		strings.ToUpper(req.Method),
		path,
		bytes.NewReader(decodedBody),
	)

	if err != nil {
		fmt.Printf("Could not convert request %s:%s to http.Request\n", req.Method, req.Url)
		log.Println(err)
		return nil, err
	}

	httpRequest.RemoteAddr = req.Headers["X-REMOTE-ADDR"]

	for h := range req.Headers {
		httpRequest.Header.Add(h, req.Headers[h])
	}

	httpRequest.RequestURI = httpRequest.URL.RequestURI()

	return httpRequest, nil
}

func ResponseToWebResponse(res *http.Response) *WebResponse {
	webResponse := &WebResponse{
		Headers: make(map[string]string),
	}

	for h := range res.Header {
		webResponse.Headers[h] = res.Header.Get(h)
	}

	webResponse.Status = res.StatusCode

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(res.Body)
	if err != nil {
		webResponse.Status = http.StatusInternalServerError
		buf.Reset()
		buf.WriteString(err.Error())
	}
	defer buf.Reset()

	webResponse.Body = buf.String()

	return webResponse
}

type WebRequest struct {
	Method  string             `json:"method"`
	Url     string             `json:"url"`
	Headers map[string]string  `json:"headers"`
	Body    string             `json:"body"`
	Files   map[string]WebFile `json:"files"`
	Schema  string             `json:"schema"`
}

type WebFile struct {
	Name     string `json:"name"`
	FullPath string `json:"full_path"`
	Type     string `json:"type"`
	TmpName  string `json:"tmp_name"`
	Error    int    `json:"error"`
	Size     int    `json:"size"`
}

type WebResponse struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
	Status  int               `json:"status"`
}

func (r *WebResponse) String() string {
	response, _ := json.Marshal(r)
	return string(response)
}

func Start(name string, handler http.Handler) {
	Export(name, map[string]interface{}{
		"run": func(requestJson string) string {
			var err error
			var request WebRequest
			var req *http.Request
			response := &WebResponse{}
			err = json.Unmarshal([]byte(requestJson), &request)
			if err != nil {
				response.Status = http.StatusBadRequest
				response.Body = fmt.Sprintf(`{"error": "Could not unmarshal request: %s"}`, err)

				return response.String()
			}

			req, err = WebRequestToRequest(&request)
			if err != nil {
				return fmt.Sprintf(`{"error": "Could not convert request: %s"}`, err)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			return ResponseToWebResponse(w.Result()).String()
		},
	})
}
