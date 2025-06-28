package user_test

import (
	"fmt"
	"testing"

	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	helper "github.com/advanced-coder-com/go-timekeeper/tests/integration/helper"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"github.com/gin-gonic/gin"
)

func TestDeleteUserSuccess(t *testing.T) {
	_ = os.Setenv("APP_ENV_FILE", ".env.test")
	helper.InitConfig("../../../.env.test")
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
	testingVariables.Password = "P@ssw0rd"

	if ok, _ := helper.SignUp(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign up user. Email: %s", testingVariables.Email)
	}
	if ok, _ := helper.SignIn(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign in user. Email: %s", testingVariables.Email)
	}
	if ok, _ := helper.DeleteUser(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to delete user. Email: %s", testingVariables.Email)
	}
	t.Logf("✅ Successfully deleted user. Email: %s", testingVariables.Email)
}

func TestDeleteAlreadyDeletedUserFails(t *testing.T) {
	_ = os.Setenv("APP_ENV_FILE", ".env.test")
	helper.InitConfig("../../../.env.test")
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
	testingVariables.Password = "P@ssw0rd"

	if ok, _ := helper.SignUp(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign up user. Email: %s", testingVariables.Email)
	}
	if ok, _ := helper.SignIn(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign in user. Email: %s", testingVariables.Email)
	}
	if ok, _ := helper.DeleteUser(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to delete user. Email: %s", testingVariables.Email)
	}

	ok, resp := helper.DeleteUser(t, &client, server, testingVariables)
	if ok {
		t.Fatalf("❌ Deletion of already deleted user should have failed. Email: %s", testingVariables.Email)
	}

	var errorMessage = helper.ErrorMessage
	helper.DecodeJSON(t, resp.Body, &errorMessage)
	if errorMessage.ErrorMessage == service.ErrUserDeleteFailed.Error() {
		t.Logf("✅ Attempt to delete already deleted user failed as expected. Error: %s", errorMessage.ErrorMessage)
	} else {
		t.Fatalf("❌ Unexpected error message when deleting already deleted user. Email: %s, got: %s", testingVariables.Email, errorMessage.ErrorMessage)
	}
}
