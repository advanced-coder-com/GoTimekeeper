package integration_test_helper

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func CreateTask(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
	projectIdIndex uint64,
	taskName string,
) {
	taskBody := map[string]interface{}{
		"name":       taskName,
		"project_id": testVars.ProjectID[projectIdIndex],
	}
	projectResp := DoPostAuth(t, client, server.URL+"/api/tasks/create", taskBody, testVars.AuthToken)
	if projectResp.StatusCode != http.StatusCreated && projectResp.StatusCode != http.StatusOK {
		t.Fatalf("task creation failed: status %d", projectResp.StatusCode)
	}
	var taskData struct {
		ID   *uint64 `json:"id"`
		Name string  `json:"name"`
	}
	DecodeJSON(t, projectResp.Body, &taskData)

	if taskData.ID == nil {
		t.Fatal("invalid task ID returned")
	}
	testVars.TaskID = append(testVars.TaskID, *taskData.ID)
}

func StartTask(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
	taskId uint64,
) {
	url := server.URL + "/api/tasks/start/" + strconv.FormatUint(taskId, 10)

	taskStartResp := DoGetAuth(t, client, url, testVars.AuthToken)
	if taskStartResp.StatusCode != http.StatusOK {
		t.Fatalf("start task failed: status %d", taskStartResp.StatusCode)
	}
}

func StopTask(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
	taskId uint64,
) {
	url := server.URL + "/api/tasks/stop/" + strconv.FormatUint(taskId, 10)

	taskStartResp := DoGetAuth(t, client, url, testVars.AuthToken)
	if taskStartResp.StatusCode != http.StatusOK {
		t.Fatalf("stop task failed: status %d", taskStartResp.StatusCode)
	}
}
