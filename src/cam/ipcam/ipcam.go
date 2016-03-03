package ipcam

import (
	"bytes"
	"cam/config"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Sink func(id string, buf []byte)

func New(cfg config.Source, sink Sink) (src *Source) {
	src = &Source{
		Name: cfg.Name,
		cfg:  cfg,
		sink: sink,
	}

	return
}

type Source struct {
	Name string

	cfg  config.Source
	sink Sink

	req    *http.Request
	client *http.Client
	buf    bytes.Buffer
}

func (src *Source) Run() (err error) {
	log.Println("connecting", src.Name, "to", src.cfg.URL)

	src.req, err = http.NewRequest("GET", src.cfg.URL, nil)

	if err != nil {
		return
	}

	if src.cfg.User != "" {
		src.req.SetBasicAuth(src.cfg.User, src.cfg.Pass)
	}

	src.client = &http.Client{}

	refresh := time.Duration(src.cfg.Refresh * float32(time.Second))

	for {

		time.Sleep(refresh)

		err = src.fetch()

		if err != nil {
			log.Printf("source %s error: %v", src.cfg.Name, err)
			time.Sleep(5 * time.Second)
		}
	}
}

func (src *Source) fetch() (err error) {
	rsp, err := src.client.Do(src.req)

	if err != nil {
		return
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected http status: %s", rsp.Status)
		return
	}

	if ctype := rsp.Header.Get("content-type"); ctype != "image/jpeg" {
		err = fmt.Errorf("expected \"Content-Type: image/jpeg\", received %s", ctype)
		return
	}

	src.buf.Reset()

	_, err = src.buf.ReadFrom(rsp.Body)

	if err != nil {
		return
	}

	src.sink(src.cfg.Name, src.buf.Bytes())

	return
}
