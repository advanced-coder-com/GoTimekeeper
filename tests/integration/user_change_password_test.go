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

func TestUserChangePasswordSuccess(t *testing.T) {
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
	if !result {
		t.Fatalf("SignUp user failed. Email %s", testingVariables.Email)
	}
	// 2. SignIn
	result, _ = helper.SignIn(t, &client, server, testingVariables)
	if result != true {
		t.Fatalf("Sign user failed. Email %s", testingVariables.Email)
	}

	// 3. Change Password
	result, response := helper.ChangePassword(t, &client, server, testingVariables, "new_password")
	if result != true {
		errorMessage := helper.ErrorMessage
		helper.DecodeJSON(t, response.Body, &errorMessage)

		t.Fatalf(
			"Change user password failed. Email %s, new password %s, message %s",
			testingVariables.Email,
			testingVariables.Password,
			errorMessage.ErrorMessage,
		)
	}

	// 4. SignIn with new password
	result, _ = helper.SignIn(t, &client, server, testingVariables)
	if result != true {
		t.Fatalf("Sign user failed. Email %s", testingVariables.Email)
	}
	t.Logf(
		"✅ Successfully Change user password with email %s and new password %s. Token %s",
		testingVariables.Email,
		testingVariables.Password,
		testingVariables.AuthToken,
	)
}
