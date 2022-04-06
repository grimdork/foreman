package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/francoispqt/gojay"
	"github.com/grimdork/foreman/api"
	"github.com/grimdork/xos"
)

// Config holds client configuration.
type Config struct {
	ServerURL string
	ID        string
	Key       string
}

// UnmarshalJSONObject decodes this config from JSON via gojay.
func (cfg *Config) UnmarshalJSONObject(dec *gojay.Decoder, key string) error {
	switch key {
	case "serverurl":
		return dec.String(&cfg.ServerURL)
	case "id":
		return dec.String(&cfg.ID)
	case "key":
		return dec.String(&cfg.Key)
	}
	return nil
}

// NKeys is required to unmarshal.
func (cfg *Config) NKeys() int {
	return 3
}

// MarshalJSONObject encodes this config to JSON via gojay.
func (cfg *Config) MarshalJSONObject(enc *gojay.Encoder) {
	enc.StringKey("serverurl", cfg.ServerURL)
	enc.StringKey("id", cfg.ID)
	enc.StringKey("key", cfg.Key)
}

// IsNil returns true if this config is nil.
func (cfg *Config) IsNil() bool {
	return cfg == nil
}

// ConfigPath constant.
const (
	ConfigPath = "foreman"
	ConfigFile = "config.json"
)

func loadConfig() (*Config, error) {
	path, err := xos.NewConfig(ConfigPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	f, err := os.Open(filepath.Join(path.Path(), ConfigFile))
	if err != nil {
		return &cfg, nil
	}

	defer f.Close()
	dec := gojay.NewDecoder(f)
	err = dec.DecodeObject(&cfg)
	return &cfg, err
}

func saveConfig(cfg *Config) error {
	path, err := xos.NewConfig(ConfigPath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(path.Path(), 0700)
	if err != nil {
		return err
	}

	buf, err := gojay.MarshalJSONObject(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(path.Path(), ConfigFile), buf, 0600)
}

func (cfg *Config) request(method, ep string, args api.Request) (*http.Response, error) {
	url := fmt.Sprintf("%s/api%s", cfg.ServerURL, ep)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("id", cfg.ID)
	req.Header.Set("key", cfg.Key)
	for k, v := range args {
		req.Header.Set(k, v)
	}

	c := &http.Client{Timeout: time.Second * 10}
	res, err := c.Do(req)
	return res, err
}

// Get is for retrieval.
func (cfg *Config) Get(ep string, args api.Request) ([]byte, error) {
	res, err := cfg.request(http.MethodGet, ep, args)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("couldn't GET: %s", res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		var e api.ErrorResponse
		err = json.Unmarshal(data, &e)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("GET failure: %s", e.Error)
	}

	return data, nil
}

// Post is for creation.
func (cfg *Config) Post(ep string, args api.Request) ([]byte, error) {
	res, err := cfg.request(http.MethodPost, ep, args)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("couldn't POST: %s", res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		var e api.ErrorResponse
		err = json.Unmarshal(data, &e)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("POST failure: %s", e.Error)
	}

	return data, nil
}

// Delete requests.
func (cfg *Config) Delete(ep string, args api.Request) ([]byte, error) {
	res, err := cfg.request(http.MethodDelete, ep, args)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("couldn't DELETE: %s", res.Status)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		var e api.ErrorResponse
		err = json.Unmarshal(data, &e)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("DELETE failure: %s", e.Error)
	}

	return data, nil
}
