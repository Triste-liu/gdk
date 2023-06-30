package writer

import (
	"bytes"
	"net/http"
)

type HttpWriter struct {
	Url    string
	Method string
}

func (w *HttpWriter) Write(p []byte) (n int, err error) {
	body := bytes.NewReader(p)
	newRequest, err := http.NewRequest(w.Method, w.Url, body)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(newRequest)
	if err != nil {
		return
	}
	defer res.Body.Close()
	return len(p), nil
}
