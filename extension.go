package pagerduty

import (
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Extension is a third-party extension to PagerDuty's API.
type Extension struct {
	APIObject
	Name             string    `json:"name,omitempty"`
	EndpointURL      string    `json:"endpoint_url,omitempty"`
	ExtensionSchema  APIObject `json:"extension_schema,omitempty"`
	ExtensionObjects APIObject `json:"extension_objects,omitempty"`
}

// ListExtensionOptions are the options available when calling the ListExtensions API endpoint.
type ListExtensionOptions struct {
	APIListObject
	ExtensionObjects []APIObject `url:"extension_objects,omitempty,brackets"`
	Query            string      `query:"query,omitempty"`
	ExtensionSchema  []APIObject `url:"extension_schema,omitempty,brackets"`
	Includes         []string    `url:"include,omitempty,brackets"`
}

// ListExtensionResponse is the response when calling the ListAddons API endpoint.
type ListExtensionResponse struct {
	APIListObject
	Extensions []Extension `json:"extensions"`
}

// ListExtensions lists all of the add-ons installed on your account.
func (c *Client) ListExtensions(o ListExtensionOptions) (*ListExtensionResponse, error) {
	v, err := query.Values(o)
	if err != nil {
		return nil, err
	}
	resp, err := c.get("/extensions?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result ListExtensionResponse
	return &result, c.decodeJSON(resp, &result)
}

// CreateExtension updates an existing extension.
func (c *Client) CreateExtension(a Extension) (*Extension, error) {
	v := make(map[string]Extension)
	v["extension"] = a
	resp, err := c.post("/extensions", v, nil)
	if err != nil {
		return nil, err
	}
	return getExtensionFromResponse(c, resp)
}

// DeleteExtension deletes an extension from your account.
func (c *Client) DeleteExtension(id string) error {
	_, err := c.delete("/extensions/" + id)
	return err
}

// GetExtension gets details about an existing extension.
func (c *Client) GetExtension(id string) (*Extension, error) {
	resp, err := c.get("/extensions/" + id)
	if err != nil {
		return nil, err
	}
	return getExtensionFromResponse(c, resp)
}

// UpdateExtension updates an existing extension.
func (c *Client) UpdateExtension(id string, a Extension) (*Extension, error) {
	v := make(map[string]Extension)
	v["extension"] = a
	resp, err := c.put("/extensions/"+id, v, nil)
	if err != nil {
		return nil, err
	}
	return getExtensionFromResponse(c, resp)
}

func getExtensionFromResponse(c *Client, resp *http.Response) (*Extension, error) {
	var result map[string]Extension
	if err := c.decodeJSON(resp, &result); err != nil {
		return nil, err
	}
	a, ok := result["extension"]
	if !ok {
		return nil, fmt.Errorf("JSON response does not have 'addon' field")
	}
	return &a, nil
}
