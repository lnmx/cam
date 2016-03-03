package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (a *app) runHttp() (err error) {
	cfg := a.cfg.Server

	http.Handle("/", a)

	log.Println("http server listening on", cfg.Addr)

	return http.ListenAndServe(cfg.Addr, nil)
}

func (a *app) ServeHTTP(rsp http.ResponseWriter, req *http.Request) {
	err := a.dispatch(rsp, req)

	if err != nil {
		http.Error(rsp, err.Error(), 500)
		return
	}

	return
}

func (a *app) dispatch(rsp http.ResponseWriter, req *http.Request) (err error) {
	verb := req.Method
	path := req.URL.Path

	log.Println(verb, path)

	for id, _ := range a.feeds {
		if verb == "GET" && path == fmt.Sprintf("/%s.jpeg", id) {
			return a.getJpeg(id, rsp, req)
		}

		if verb == "GET" && path == fmt.Sprintf("/%s.mjpeg", id) {
			return a.getMjpeg(id, rsp, req)
		}
	}

	http.NotFound(rsp, req)

	return
}

func (a *app) getJpeg(id string, rsp http.ResponseWriter, req *http.Request) (err error) {
	if buf, ok := a.getFrame(id); ok {
		rsp.Header().Set("Content-Type", "image/jpeg")
		rsp.Write(buf)
		return
	}

	http.NotFound(rsp, req)

	return
}

func (a *app) getMjpeg(id string, rsp http.ResponseWriter, req *http.Request) (err error) {
	boundary := fmt.Sprintf("--frame--%d----", time.Now().UnixNano())

	req.ProtoMinor = 0

	rsp.Header().Set("Transfer-Encoding", "identity")
	rsp.Header().Set("Content-Type", "multipart/x-mixed-replace;boundary="+boundary)
	rsp.WriteHeader(http.StatusOK)

	sub := a.subscribe(id)

	defer a.unsubscribe(id, sub)

	for _ = range sub {
		data, ok := a.getFrame(id)

		if !ok {
			break
		}

		fmt.Fprintf(rsp, "--%s\r\n", boundary)
		fmt.Fprintf(rsp, "Content-Type: image/jpeg\r\n")
		fmt.Fprintf(rsp, "Content-Length: %d\r\n", len(data))
		fmt.Fprintf(rsp, "\r\n")

		rsp.Write(data)

		_, err = fmt.Fprintf(rsp, "\r\n\r\n")

		if err != nil {
			break
		}
	}

	return nil
}
