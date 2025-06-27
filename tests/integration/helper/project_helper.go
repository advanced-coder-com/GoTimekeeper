package integration_test_helper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func CreateProject(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
	projectName string,
) {
	projectBody := map[string]string{
		"name": projectName,
	}
	projectResp := DoPostAuth(t, client, server.URL+"/api/projects/create", projectBody, testVars.AuthToken)
	if projectResp.StatusCode != http.StatusCreated && projectResp.StatusCode != http.StatusOK {
		t.Fatalf("project creation failed: status %d", projectResp.StatusCode)
	}
	var projectData struct {
		ID   *uint64 `json:"id"`
		Name string  `json:"name"`
	}
	DecodeJSON(t, projectResp.Body, &projectData)

	if projectData.ID == nil {
		t.Fatal("invalid project ID returned")
	}
	testVars.ProjectID = append(testVars.ProjectID, *projectData.ID)
}
