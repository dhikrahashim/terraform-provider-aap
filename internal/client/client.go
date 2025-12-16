package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	Host     string
	Username string
	Password string
	Token    string
	Insecure bool
	HTTP     *http.Client
}

type Organization struct {
	ID               int    `json:"id,omitempty"`
	Name             string `json:"name"`
	Description      string `json:"description,omitempty"`
	MaxHosts         int    `json:"max_hosts,omitempty"`
	CustomVirtualEnv string `json:"custom_virtualenv,omitempty"`
}

func NewClient(host, username, password, token string, insecure bool) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	c := &Client{
		Host:     strings.TrimRight(host, "/"),
		Username: username,
		Password: password,
		Token:    token,
		Insecure: insecure,
		HTTP: &http.Client{
			Timeout:   30 * time.Second,
			Transport: tr,
		},
	}

	return c
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.Host, path), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	} else if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	return io.ReadAll(resp.Body)
}

// GetOrganization retrieves an organization by ID
func (c *Client) GetOrganization(id int) (*Organization, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/organizations/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	var org Organization
	err = json.Unmarshal(resp, &org)
	return &org, err
}

// CreateOrganization creates a new organization
func (c *Client) CreateOrganization(org *Organization) (*Organization, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/organizations/", org)
	if err != nil {
		return nil, err
	}

	var newOrg Organization
	err = json.Unmarshal(resp, &newOrg)
	return &newOrg, err
}

// UpdateOrganization updates an existing organization
func (c *Client) UpdateOrganization(org *Organization) (*Organization, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/organizations/%d/", org.ID), org)
	if err != nil {
		return nil, err
	}

	var updatedOrg Organization
	err = json.Unmarshal(resp, &updatedOrg)
	return &updatedOrg, err
}

// DeleteOrganization deletes an organization
func (c *Client) DeleteOrganization(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/organizations/%d/", id), nil)
	return err
}

type Inventory struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	Organization int    `json:"organization"`
	Kind         string `json:"kind,omitempty"`
	HostFilter   string `json:"host_filter,omitempty"`
	Variables    string `json:"variables,omitempty"`
}

// GetInventory retrieves an inventory by ID
func (c *Client) GetInventory(id int) (*Inventory, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/inventories/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	var inv Inventory
	err = json.Unmarshal(resp, &inv)
	return &inv, err
}

// CreateInventory creates a new inventory
func (c *Client) CreateInventory(inv *Inventory) (*Inventory, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/inventories/", inv)
	if err != nil {
		return nil, err
	}

	var newInv Inventory
	err = json.Unmarshal(resp, &newInv)
	return &newInv, err
}

// UpdateInventory updates an existing inventory
func (c *Client) UpdateInventory(inv *Inventory) (*Inventory, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/inventories/%d/", inv.ID), inv)
	if err != nil {
		return nil, err
	}

	var updatedInv Inventory
	err = json.Unmarshal(resp, &updatedInv)
	return &updatedInv, err
}

// DeleteInventory deletes an inventory
func (c *Client) DeleteInventory(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/inventories/%d/", id), nil)
	return err
}

type JobTemplate struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	JobType      string `json:"job_type"`
	Inventory    int    `json:"inventory"`
	Project      int    `json:"project"`
	Playbook     string `json:"playbook"`
	Forks        int    `json:"forks,omitempty"`
	Limit        string `json:"limit,omitempty"`
	Verbosity    int    `json:"verbosity,omitempty"`
	ExtraVars    string `json:"extra_vars,omitempty"`
}

// GetJobTemplate retrieves a job template by ID
func (c *Client) GetJobTemplate(id int) (*JobTemplate, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/job_templates/%d/", id), nil)
	if err != nil {
		return nil, err
	}

	var jt JobTemplate
	err = json.Unmarshal(resp, &jt)
	return &jt, err
}

// CreateJobTemplate creates a new job template
func (c *Client) CreateJobTemplate(jt *JobTemplate) (*JobTemplate, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/job_templates/", jt)
	if err != nil {
		return nil, err
	}

	var newJt JobTemplate
	err = json.Unmarshal(resp, &newJt)
	return &newJt, err
}

// UpdateJobTemplate updates an existing job template
func (c *Client) UpdateJobTemplate(jt *JobTemplate) (*JobTemplate, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/job_templates/%d/", jt.ID), jt)
	if err != nil {
		return nil, err
	}

	var updatedJt JobTemplate
	err = json.Unmarshal(resp, &updatedJt)
	return &updatedJt, err
}

// DeleteJobTemplate deletes a job template
func (c *Client) DeleteJobTemplate(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/job_templates/%d/", id), nil)
	return err
}
