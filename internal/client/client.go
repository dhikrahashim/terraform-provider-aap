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

// ==================== PROJECT ====================

type Project struct {
	ID                    int    `json:"id,omitempty"`
	Name                  string `json:"name"`
	Description           string `json:"description,omitempty"`
	Organization          int    `json:"organization"`
	ScmType               string `json:"scm_type"`
	ScmUrl                string `json:"scm_url,omitempty"`
	ScmBranch             string `json:"scm_branch,omitempty"`
	ScmCredential         int    `json:"credential,omitempty"`
	ScmClean              bool   `json:"scm_clean,omitempty"`
	ScmDeleteOnUpdate     bool   `json:"scm_delete_on_update,omitempty"`
	ScmUpdateOnLaunch     bool   `json:"scm_update_on_launch,omitempty"`
	ScmUpdateCacheTimeout int    `json:"scm_update_cache_timeout,omitempty"`
	LocalPath             string `json:"local_path,omitempty"`
}

func (c *Client) GetProject(id int) (*Project, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/projects/%d/", id), nil)
	if err != nil {
		return nil, err
	}
	var p Project
	err = json.Unmarshal(resp, &p)
	return &p, err
}

func (c *Client) CreateProject(p *Project) (*Project, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/projects/", p)
	if err != nil {
		return nil, err
	}
	var newP Project
	err = json.Unmarshal(resp, &newP)
	return &newP, err
}

func (c *Client) UpdateProject(p *Project) (*Project, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/projects/%d/", p.ID), p)
	if err != nil {
		return nil, err
	}
	var updated Project
	err = json.Unmarshal(resp, &updated)
	return &updated, err
}

func (c *Client) DeleteProject(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/projects/%d/", id), nil)
	return err
}

// ==================== CREDENTIAL ====================

type CredentialInputs struct {
	Username         string `json:"username,omitempty"`
	Password         string `json:"password,omitempty"`
	SSHKeyData       string `json:"ssh_key_data,omitempty"`
	SSHPublicKeyData string `json:"ssh_public_key_data,omitempty"`
	SSHKeyUnlock     string `json:"ssh_key_unlock,omitempty"`
	BecomeMethod     string `json:"become_method,omitempty"`
	BecomeUsername   string `json:"become_username,omitempty"`
	BecomePassword   string `json:"become_password,omitempty"`
}

type Credential struct {
	ID             int              `json:"id,omitempty"`
	Name           string           `json:"name"`
	Description    string           `json:"description,omitempty"`
	Organization   int              `json:"organization"`
	CredentialType int              `json:"credential_type"`
	Inputs         CredentialInputs `json:"inputs,omitempty"`
}

func (c *Client) GetCredential(id int) (*Credential, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/credentials/%d/", id), nil)
	if err != nil {
		return nil, err
	}
	var cred Credential
	err = json.Unmarshal(resp, &cred)
	return &cred, err
}

func (c *Client) CreateCredential(cred *Credential) (*Credential, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/credentials/", cred)
	if err != nil {
		return nil, err
	}
	var newCred Credential
	err = json.Unmarshal(resp, &newCred)
	return &newCred, err
}

func (c *Client) UpdateCredential(cred *Credential) (*Credential, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/credentials/%d/", cred.ID), cred)
	if err != nil {
		return nil, err
	}
	var updated Credential
	err = json.Unmarshal(resp, &updated)
	return &updated, err
}

func (c *Client) DeleteCredential(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/credentials/%d/", id), nil)
	return err
}

// ==================== INVENTORY SOURCE ====================

type InventorySource struct {
	ID                 int    `json:"id,omitempty"`
	Name               string `json:"name"`
	Description        string `json:"description,omitempty"`
	Inventory          int    `json:"inventory"`
	Source             string `json:"source"`
	SourcePath         string `json:"source_path,omitempty"`
	SourceVars         string `json:"source_vars,omitempty"`
	Credential         int    `json:"credential,omitempty"`
	SourceProject      int    `json:"source_project,omitempty"`
	UpdateOnLaunch     bool   `json:"update_on_launch,omitempty"`
	UpdateCacheTimeout int    `json:"update_cache_timeout,omitempty"`
	Overwrite          bool   `json:"overwrite,omitempty"`
	OverwriteVars      bool   `json:"overwrite_vars,omitempty"`
}

func (c *Client) GetInventorySource(id int) (*InventorySource, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/inventory_sources/%d/", id), nil)
	if err != nil {
		return nil, err
	}
	var is InventorySource
	err = json.Unmarshal(resp, &is)
	return &is, err
}

func (c *Client) CreateInventorySource(is *InventorySource) (*InventorySource, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/inventory_sources/", is)
	if err != nil {
		return nil, err
	}
	var newIS InventorySource
	err = json.Unmarshal(resp, &newIS)
	return &newIS, err
}

func (c *Client) UpdateInventorySource(is *InventorySource) (*InventorySource, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/inventory_sources/%d/", is.ID), is)
	if err != nil {
		return nil, err
	}
	var updated InventorySource
	err = json.Unmarshal(resp, &updated)
	return &updated, err
}

func (c *Client) DeleteInventorySource(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/inventory_sources/%d/", id), nil)
	return err
}

// ==================== CREDENTIAL TYPE ====================

type CredentialType struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Kind        string `json:"kind"`
	Inputs      string `json:"inputs,omitempty"`
	Injectors   string `json:"injectors,omitempty"`
}

func (c *Client) GetCredentialType(id int) (*CredentialType, error) {
	resp, err := c.doRequest("GET", fmt.Sprintf("/api/controller/v2/credential_types/%d/", id), nil)
	if err != nil {
		return nil, err
	}
	var ct CredentialType
	err = json.Unmarshal(resp, &ct)
	return &ct, err
}

func (c *Client) CreateCredentialType(ct *CredentialType) (*CredentialType, error) {
	resp, err := c.doRequest("POST", "/api/controller/v2/credential_types/", ct)
	if err != nil {
		return nil, err
	}
	var newCT CredentialType
	err = json.Unmarshal(resp, &newCT)
	return &newCT, err
}

func (c *Client) UpdateCredentialType(ct *CredentialType) (*CredentialType, error) {
	resp, err := c.doRequest("PATCH", fmt.Sprintf("/api/controller/v2/credential_types/%d/", ct.ID), ct)
	if err != nil {
		return nil, err
	}
	var updated CredentialType
	err = json.Unmarshal(resp, &updated)
	return &updated, err
}

func (c *Client) DeleteCredentialType(id int) error {
	_, err := c.doRequest("DELETE", fmt.Sprintf("/api/controller/v2/credential_types/%d/", id), nil)
	return err
}
