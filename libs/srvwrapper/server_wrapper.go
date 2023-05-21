package srvwrapper

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Wrapper[Req, Res any] struct {
	fn func(req Req) (Res, error)
}

func New[Req, Res any](fn func(req Req) (Res, error)) *Wrapper[Req, Res] {
	return &Wrapper[Req, Res]{
		fn: fn,
	}
}

func (w *Wrapper[Req, Res]) ServeHTTP(resWriter http.ResponseWriter, httpReq *http.Request) {

	var req Req

	err := json.NewDecoder(httpReq.Body).Decode(&req)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "parse request", err)
		return
	}

	resp, err := w.fn(req)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "exec handler", err)
		return
	}

	rawData, err := json.Marshal(&resp)
	if err != nil {
		resWriter.WriteHeader(http.StatusInternalServerError)
		writeErrorText(resWriter, "decode response", err)
		return
	}

	_, _ = resWriter.Write(rawData)
}

func writeErrorText(w http.ResponseWriter, text string, err error) {
	buf := bytes.NewBufferString(text)

	buf.WriteString(": ")
	buf.WriteString(err.Error())
	buf.WriteByte('\n')

	w.Write(buf.Bytes())
}
