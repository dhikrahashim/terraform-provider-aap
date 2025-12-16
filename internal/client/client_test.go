package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetOrganization(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected method GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/controller/v2/organizations/1/" {
			t.Errorf("Expected path /api/controller/v2/organizations/1/, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Organization{
			ID:   1,
			Name: "Test Org",
			MaxHosts: 10,
		})
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "pass", "", true)
	org, err := c.GetOrganization(1)
	if err != nil {
		t.Fatalf("GetOrganization failed: %s", err)
	}

	if org.Name != "Test Org" {
		t.Errorf("Expected name 'Test Org', got %s", org.Name)
	}
	if org.MaxHosts != 10 {
		t.Errorf("Expected MaxHosts 10, got %d", org.MaxHosts)
	}
}

func TestCreateJobTemplate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected method POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/controller/v2/job_templates/" {
			t.Errorf("Expected path /api/controller/v2/job_templates/, got %s", r.URL.Path)
		}

		var reqBody JobTemplate
		json.NewDecoder(r.Body).Decode(&reqBody)
		if reqBody.Name != "New Job" {
			t.Errorf("Expected name 'New Job', got %s", reqBody.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(JobTemplate{
			ID:      2,
			Name:    reqBody.Name,
			JobType: reqBody.JobType,
		})
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "user", "pass", "", true)
	jt, err := c.CreateJobTemplate(&JobTemplate{
		Name:    "New Job",
		JobType: "run",
	})
	if err != nil {
		t.Fatalf("CreateJobTemplate failed: %s", err)
	}

	if jt.ID != 2 {
		t.Errorf("Expected ID 2, got %d", jt.ID)
	}
}
