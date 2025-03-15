package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SupportedMethod string

const (
	MethodDelete SupportedMethod = http.MethodDelete
	MethodGet    SupportedMethod = http.MethodGet
	MethodPatch  SupportedMethod = http.MethodPatch
	MethodPost   SupportedMethod = http.MethodPost
	MethodPut    SupportedMethod = http.MethodPut
)

type RequestConfig struct {
	Method  SupportedMethod
	Headers map[string]string
	URL     *url.URL
	Body    []byte
}

type HttpClient struct {
	baseURL       string
	beforeRequest func(*RequestConfig)
	client        *http.Client
}

func NewHttpClient(baseURL string, beforeRequest func(*RequestConfig)) *HttpClient {
	return &HttpClient{
		baseURL:       baseURL,
		beforeRequest: beforeRequest,
		client:        &http.Client{},
	}
}

func (c *HttpClient) Post(ctx context.Context, path string, body interface{}, opts map[string]string) (interface{}, error) {
	return c.fetch(ctx, path, MethodPost, body, opts)
}

func (c *HttpClient) Patch(ctx context.Context, path string, body interface{}, opts map[string]string) (interface{}, error) {
	return c.fetch(ctx, path, MethodPatch, body, opts)
}

func (c *HttpClient) Put(ctx context.Context, path string, body interface{}, opts map[string]string) (interface{}, error) {
	return c.fetch(ctx, path, MethodPut, body, opts)
}

func (c *HttpClient) Delete(ctx context.Context, path string, opts map[string]string) (interface{}, error) {
	return c.fetch(ctx, path, MethodDelete, nil, opts)
}

func (c *HttpClient) Get(ctx context.Context, path string, opts map[string]string) (interface{}, error) {
	return c.fetch(ctx, path, MethodGet, nil, opts)
}

func (c *HttpClient) fetch(ctx context.Context, path string, method SupportedMethod, body interface{}, opts map[string]string) (interface{}, error) {
	config, err := c.makeConfig(path, method, body, opts)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, string(config.Method), config.URL.String(), bytes.NewBuffer(config.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if len(respBody) == 0 {
			return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
		}
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d, response: %s", resp.StatusCode, string(respBody)))

	}

	if len(respBody) == 0 {
		return nil, nil
	}

	var result interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *HttpClient) FetchRawResponse(ctx context.Context, path string, method SupportedMethod, body interface{}, opts map[string]string) (*http.Response, error) {
	config, err := c.makeConfig(path, method, body, opts)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, string(config.Method), config.URL.String(), bytes.NewBuffer(config.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range config.Headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}

func (c *HttpClient) makeConfig(path string, method SupportedMethod, body interface{}, opts map[string]string) (*RequestConfig, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, err
	}
	u, err = u.Parse(path)
	if err != nil {
		return nil, err
	}

	var bodyBytes []byte
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	for k, v := range opts {
		headers[k] = v
	}

	config := &RequestConfig{
		Method:  method,
		Headers: headers,
		URL:     u,
		Body:    bodyBytes,
	}

	if c.beforeRequest != nil {
		c.beforeRequest(config)
	}

	return config, nil
}
