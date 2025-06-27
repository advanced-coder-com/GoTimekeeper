package integration_test

import (
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/service"
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

func TestUserLoginFailedBadPassword(t *testing.T) {
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
	// 1. SignUp first user
	result, _ := helper.SignUp(t, &client, server, testingVariables)
	if !result {
		t.Fatalf("SignUp user failed. Email %s", testingVariables.Email)
	}

	testingVariables.Password = "wrong-password"

	// 2. SignIn
	result, response := helper.SignIn(t, &client, server, testingVariables)
	if result == true {
		t.Fatalf("Login user with INVALID password was successfull!!")
	}
	errorMessage := helper.ErrorMessage
	helper.DecodeJSON(t, response.Body, &errorMessage)
	if errorMessage.ErrorMessage == service.ErrUserSignInFailed.Error() {
		t.Logf("✅ LogIn wih bad password test passed. Error message: %s", errorMessage.ErrorMessage)
	} else {
		t.Fatalf(
			"LogIn wih bad password test failed. Email: %s, password %s, API error message %s",
			testingVariables.Email,
			testingVariables.Password,
			errorMessage.ErrorMessage,
		)
	}
}
