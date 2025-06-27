package integration_test

import (
	"fmt"
	helper "github.com/advanced-coder-com/go-timekeeper/tests/integration/helper"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"github.com/gin-gonic/gin"
)

func TestHappyScenario(t *testing.T) {
	_ = os.Setenv("APP_ENV_FILE", ".env.test")
	helper.InitConfig()
	fmt.Println(viper.GetString("DB_HOST"))
	db.Init()

	engine := gin.Default()
	gin.SetMode(gin.TestMode)
	router.SetupRoutes(engine)

	server := httptest.NewServer(engine)
	defer server.Close()

	client := http.Client{}

	testingVariables := &helper.TestingContext{}
	testingVariables.Email = "user" + uuid.NewString() + "@example.com"
	testingVariables.Password = "password"
	// 1. SignUp
	result, _ := helper.SignUp(t, &client, server, testingVariables)
	if result {
		t.Logf("✅ Successfully created user with email %s", testingVariables.Email)
	} else {
		t.Fatalf("SignUp user failed. Email %s", testingVariables.Email)
	}

	// 2. SignIn
	helper.SignIn(t, &client, server, testingVariables)
	t.Logf("✅ Successfully Sign in user with email %s token %s", testingVariables.Email, testingVariables.AuthToken)

	// 3. Create projects
	helper.CreateProject(t, &client, server, testingVariables, "My First Project")
	t.Logf("✅ Successfully created project with ID %d", testingVariables.ProjectID[0])

	helper.CreateProject(t, &client, server, testingVariables, "My Second Project")
	t.Logf("✅ Successfully created project with ID %d", testingVariables.ProjectID[1])

	// Create tasks
	helper.CreateTask(t, &client, server, testingVariables, 0, "First Task of First Project")
	t.Logf("✅ Successfully created first task in first project with ID %d", testingVariables.TaskID[0])

	helper.CreateTask(t, &client, server, testingVariables, 0, "Second Task of First Project")
	t.Logf("✅ Successfully created second task in first project with ID %d", testingVariables.TaskID[1])

	helper.CreateTask(t, &client, server, testingVariables, 1, "First Task of Second Project")
	t.Logf("✅ Successfully created first task in Second project with ID %d", testingVariables.TaskID[2])

	helper.CreateTask(t, &client, server, testingVariables, 1, "Second Task of Second Project")
	t.Logf("✅ Successfully created second task in Second project with ID %d", testingVariables.TaskID[3])

	helper.StartTask(t, &client, server, testingVariables, testingVariables.TaskID[0])
	t.Logf("✅ Successfully started first task in first project with ID %d", testingVariables.TaskID[0])

	helper.StopTask(t, &client, server, testingVariables, testingVariables.TaskID[0])
	t.Logf("✅ Successfully stopped first task in first project with ID %d", testingVariables.TaskID[0])

	helper.StartTask(t, &client, server, testingVariables, testingVariables.TaskID[3])
	t.Logf("✅ Successfully started second task in Second project with ID %d", testingVariables.TaskID[3])

	helper.StopTask(t, &client, server, testingVariables, testingVariables.TaskID[3])
	t.Logf("✅ Successfully stopped second task in Second project with ID %d", testingVariables.TaskID[3])
}
