package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		HTTPClient: http.DefaultClient,
	}
}

// buckets

func (c *Client) CreateBucket(name string) error {
	req, _ := http.NewRequest(http.MethodPut, c.BaseURL+"/buckets"+name, nil)
	return c.doCheck(req)
}

func (c *Client) DeleteBucket(name string) error {
	req, _ := http.NewRequest(http.MethodDelete, c.BaseURL+"/buckets/"+name, nil)
	return c.doCheck(req)
}

func (c *Client) ListBuckets() ([]string, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/buckets")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("list buckets failed: %s", resp.Status)
	}
	var buckets []string
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

// objects

func (c *Client) PutObject(bucket, key string, data io.Reader) error {
	url := fmt.Sprintf("%s/buckets/%s/object/%s", c.BaseURL, bucket, key)
	req, _ := http.NewRequest(http.MethodPut, url, data)
	return c.doCheck(req)
}

func (c *Client) GetObject(bucket, key string) (io.ReadCloser, error) {
	url := fmt.Sprintf("%s/buckets/%s/object/%s", c.BaseURL, bucket, key)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("get object failed: %s", resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) DeleteObject(bucket, key string) error {
	url := fmt.Sprintf("%s/buckets/%s/object/%s", c.BaseURL, bucket, key)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return c.doCheck(req)
}

func (c *Client) ListObjects(bucket string) ([]string, error) {
	url := fmt.Sprintf("%s/buckets/%s/objects", c.BaseURL, bucket)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var objs []string
	if err := json.NewDecoder(resp.Body).Decode(&objs); err != nil {
		return nil, err
	}
	return objs, nil
}

// helpers

// run a request and error on non 2XX
func (c *Client) doCheck(req *http.Request) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%s: %s", resp.Status, body)
	}
	return nil
}
