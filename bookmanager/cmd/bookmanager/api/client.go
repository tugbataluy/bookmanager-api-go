package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type APIClient struct {
	baseURL string
	verbose bool
}

func NewAPIClient(baseURL string, verbose bool) *APIClient {
	return &APIClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		verbose: verbose,
	}
}

func (c *APIClient) Get(endpoint string, queryParams map[string]string) ([]byte, error) {
    baseURL, err := url.Parse(c.baseURL)
    if err != nil {
        return nil, fmt.Errorf("invalid base URL: %w", err)
    }

    endpointURL, err := url.Parse(strings.TrimPrefix(endpoint, "/"))
    if err != nil {
        return nil, fmt.Errorf("invalid endpoint: %w", err)
    }

    fullURL := baseURL.ResolveReference(endpointURL)
    if len(queryParams) > 0 {
        q := fullURL.Query()
        for k, v := range queryParams {
            q.Add(k, v)
        }
        fullURL.RawQuery = q.Encode()
    }

    if c.verbose {
        fmt.Printf("GET %s\n", fullURL.String())
    }

    resp, err := http.Get(fullURL.String())
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    if resp.StatusCode >= 400 {
        return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
    }

    return body, nil
}

func (c *APIClient) Post(endpoint string, data interface{}) ([]byte, error) {
	return c.sendRequest("POST", endpoint, data)
}

func (c *APIClient) Put(endpoint string, data interface{}) ([]byte, error) {
	return c.sendRequest("PUT", endpoint, data)
}

func (c *APIClient) Patch(endpoint string, data interface{}) ([]byte, error) {
	return c.sendRequest("PATCH", endpoint, data)
}

func (c *APIClient) Delete(endpoint string) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(endpoint, "/"))

	if c.verbose {
		fmt.Printf("DELETE %s\n", url)
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *APIClient) sendRequest(method, endpoint string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, strings.TrimPrefix(endpoint, "/"))

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request data: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	if c.verbose {
		fmt.Printf("%s %s\n", method, url)
		if body != nil {
			fmt.Printf("Request Body: %s\n", body)
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *APIClient) BuildQueryParams(where, groupBy, orderBy string, limit, offset int) map[string]string {
	params := make(map[string]string)
	
	if where != "" {
		params["where"] = where
	}
	if groupBy != "" {
		params["group_by"] = groupBy
	}
	if orderBy != "" {
		params["order_by"] = orderBy
	}
	if limit > 0 {
		params["limit"] = strconv.Itoa(limit)
	}
	if offset > 0 {
		params["offset"] = strconv.Itoa(offset)
	}
	
	return params
}


