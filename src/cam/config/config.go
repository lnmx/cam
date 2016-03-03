package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type App struct {
	Sources []Source
	Server  Server
}

func (cfg *App) validate() (err error) {
	err = cfg.Server.validate()

	if err != nil {
		return err
	}

	for i, _ := range cfg.Sources {
		err = cfg.Sources[i].validate()

		if err != nil {
			return err
		}
	}

	return
}

type Server struct {
	Addr string
}

func (cfg *Server) validate() (err error) {
	if cfg.Addr == "" {
		cfg.Addr = "0.0.0.0:8680"
	}

	return
}

type Source struct {
	Name    string
	URL     string
	User    string
	Pass    string
	Refresh float32
}

func (cfg *Source) validate() (err error) {

	if cfg.Name == "" {
		return fmt.Errorf("\"name\" is required for source")
	}

	if cfg.URL == "" {
		return fmt.Errorf("\"url\" is required for source")
	}

	if cfg.Refresh == 0 {
		cfg.Refresh = 1
	}

	return
}

func Read(path string) (app App, err error) {
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return
	}

	err = json.Unmarshal(buf, &app)

	if err != nil {
		err = fmt.Errorf("json error in %s: %v", path, err)
		return
	}

	err = app.validate()

	if err != nil {
		err = fmt.Errorf("validation error in %s: %v", path, err)
	}

	return
}
