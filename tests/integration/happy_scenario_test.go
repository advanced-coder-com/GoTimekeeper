package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"github.com/gin-gonic/gin"
)

type testingContext struct {
	Email     string
	Password  string
	AuthToken string
	ProjectID []uint64
	TaskID    []uint64
}

const ENV = "../../.env.test"

func initConfig() {
	viper.SetConfigFile(ENV)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env.test file found or error reading it: %v", err)
	}
}

func TestHappyScenario(t *testing.T) {
	_ = os.Setenv("APP_ENV_FILE", ".env.test")
	initConfig()
	fmt.Println(viper.GetString("DB_HOST"))
	db.Init()

	engine := gin.Default()
	gin.SetMode(gin.TestMode)
	router.SetupRoutes(engine)

	server := httptest.NewServer(engine)
	defer server.Close()

	client := http.Client{}

	testingVariables := &testingContext{}
	testingVariables.Email = "user" + uuid.NewString() + "@example.com"
	testingVariables.Password = "password"
	// 1. SignUp
	testSignUp(t, &client, server, testingVariables)
	t.Logf("✅ Successfully created user with email %s", testingVariables.Email)

	// 2. SignIn
	testSignIn(t, &client, server, testingVariables)
	t.Logf("✅ Successfully Sign in user with email %s token %s", testingVariables.Email, testingVariables.AuthToken)

	// 3. Create projects
	testCreateProject(t, &client, server, testingVariables, "My First Project")
	t.Logf("✅ Successfully created project with ID %d", testingVariables.ProjectID[0])

	testCreateProject(t, &client, server, testingVariables, "My Second Project")
	t.Logf("✅ Successfully created project with ID %d", testingVariables.ProjectID[1])

	// Create tasks
	testCreateTask(t, &client, server, testingVariables, 0, "First Task of First Project")
	t.Logf("✅ Successfully created first task in first project with ID %d", testingVariables.TaskID[0])

	testCreateTask(t, &client, server, testingVariables, 0, "Second Task of First Project")
	t.Logf("✅ Successfully created second task in first project with ID %d", testingVariables.TaskID[1])

	testCreateTask(t, &client, server, testingVariables, 1, "First Task of Second Project")
	t.Logf("✅ Successfully created first task in Second project with ID %d", testingVariables.TaskID[2])

	testCreateTask(t, &client, server, testingVariables, 1, "Second Task of Second Project")
	t.Logf("✅ Successfully created second task in Second project with ID %d", testingVariables.TaskID[3])

	testStartTask(t, &client, server, testingVariables, testingVariables.TaskID[0])
	t.Logf("✅ Successfully started first task in first project with ID %d", testingVariables.TaskID[0])

	testStopTask(t, &client, server, testingVariables, testingVariables.TaskID[0])
	t.Logf("✅ Successfully stopped first task in first project with ID %d", testingVariables.TaskID[0])

	testStartTask(t, &client, server, testingVariables, testingVariables.TaskID[3])
	t.Logf("✅ Successfully started second task in Second project with ID %d", testingVariables.TaskID[3])

	testStopTask(t, &client, server, testingVariables, testingVariables.TaskID[3])
	t.Logf("✅ Successfully stopped second task in Second project with ID %d", testingVariables.TaskID[3])
}

func testSignUp(t *testing.T, client *http.Client, server *httptest.Server, testVars *testingContext) {
	signupBody := map[string]string{
		"email":    testVars.Email,
		"password": testVars.Password,
	}
	signupResp := doPost(t, client, server.URL+"/api/user/signup", signupBody)
	if signupResp.StatusCode != http.StatusCreated && signupResp.StatusCode != http.StatusOK {
		t.Fatalf("signup failed: status %d", signupResp.StatusCode)
	}
}

func testSignIn(t *testing.T, client *http.Client, server *httptest.Server, testVars *testingContext) {
	signinBody := map[string]string{
		"email":    testVars.Email,
		"password": testVars.Password,
	}
	signinResp := doPost(t, client, server.URL+"/api/user/signin", signinBody)
	var signinData struct {
		Token string `json:"token"`
	}
	decodeJSON(t, signinResp.Body, &signinData)

	if signinData.Token == "" {
		t.Fatal("no token received after signin")
	}
	testVars.AuthToken = signinData.Token
}

func testCreateProject(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *testingContext,
	projectName string,
) {
	projectBody := map[string]string{
		"name": projectName,
	}
	projectResp := doPostAuth(t, client, server.URL+"/api/projects/create", projectBody, testVars.AuthToken)
	if projectResp.StatusCode != http.StatusCreated && projectResp.StatusCode != http.StatusOK {
		t.Fatalf("project creation failed: status %d", projectResp.StatusCode)
	}
	var projectData struct {
		ID   *uint64 `json:"id"`
		Name string  `json:"name"`
	}
	decodeJSON(t, projectResp.Body, &projectData)

	if projectData.ID == nil {
		t.Fatal("invalid project ID returned")
	}
	testVars.ProjectID = append(testVars.ProjectID, *projectData.ID)
}

func testCreateTask(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *testingContext,
	projectIdIndex uint64,
	taskName string,
) {
	taskBody := map[string]interface{}{
		"name":       taskName,
		"project_id": testVars.ProjectID[projectIdIndex],
	}
	projectResp := doPostAuth(t, client, server.URL+"/api/tasks/create", taskBody, testVars.AuthToken)
	if projectResp.StatusCode != http.StatusCreated && projectResp.StatusCode != http.StatusOK {
		t.Fatalf("task creation failed: status %d", projectResp.StatusCode)
	}
	var taskData struct {
		ID   *uint64 `json:"id"`
		Name string  `json:"name"`
	}
	decodeJSON(t, projectResp.Body, &taskData)

	if taskData.ID == nil {
		t.Fatal("invalid task ID returned")
	}
	testVars.TaskID = append(testVars.TaskID, *taskData.ID)
}

func testStartTask(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *testingContext,
	taskId uint64,
) {
	url := server.URL + "/api/tasks/start/" + strconv.FormatUint(taskId, 10)

	taskStartResp := doGetAuth(t, client, url, testVars.AuthToken)
	if taskStartResp.StatusCode != http.StatusOK {
		t.Fatalf("start task failed: status %d", taskStartResp.StatusCode)
	}
}

func testStopTask(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *testingContext,
	taskId uint64,
) {
	url := server.URL + "/api/tasks/stop/" + strconv.FormatUint(taskId, 10)

	taskStartResp := doGetAuth(t, client, url, testVars.AuthToken)
	if taskStartResp.StatusCode != http.StatusOK {
		t.Fatalf("stop task failed: status %d", taskStartResp.StatusCode)
	}
}

func doPost(t *testing.T, client *http.Client, url string, body any) *http.Response {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	return resp
}

func doPostAuth(t *testing.T, client *http.Client, url string, body any, token string) *http.Response {
	buf, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		t.Fatalf("failed to create auth request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("auth HTTP request failed: %v", err)
	}
	return resp
}

func doGet(t *testing.T, client *http.Client, url string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	return resp
}

func doGetAuth(t *testing.T, client *http.Client, url string, token string) *http.Response {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create GET request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("authenticated GET request failed: %v", err)
	}
	return resp
}

func decodeJSON(t *testing.T, r io.Reader, v any) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}
