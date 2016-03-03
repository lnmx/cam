package main

import (
	"cam/config"
	"cam/ipcam"
	"log"
)

func main() {
	var app app

	err := app.run()

	if err != nil {
		log.Fatal(err)
	}
}

type app struct {
	cfg   config.App
	feeds map[string]*feed
}

func (a *app) configure() (err error) {
	a.cfg, err = config.Read("config.json")

	if err != nil {
		return
	}

	a.feeds = make(map[string]*feed)

	for _, src_cfg := range a.cfg.Sources {
		id := src_cfg.Name

		feed := feed{id: id}
		feed.src = ipcam.New(src_cfg, feed.sink)

		a.feeds[id] = &feed
	}

	return
}

func (a *app) run() (err error) {
	err = a.configure()

	if err != nil {
		return err
	}

	for _, feed := range a.feeds {
		go feed.run()
	}

	err = a.runHttp()

	return
}

func (a *app) getFrame(id string) (buf []byte, ok bool) {

	if feed, ok := a.feeds[id]; ok {
		return feed.getFrame()
	}

	return nil, false
}

func (a *app) subscribe(id string) (c chan struct{}) {
	if feed, ok := a.feeds[id]; ok {
		return feed.sub()
	}

	return nil
}

func (a *app) unsubscribe(id string, c chan struct{}) {
	if feed, ok := a.feeds[id]; ok {
		feed.unsub(c)
	}
}
