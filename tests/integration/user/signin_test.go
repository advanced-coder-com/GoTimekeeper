package user_test

import (
	"fmt"
	"github.com/advanced-coder-com/go-timekeeper/internal/service"
	helper "github.com/advanced-coder-com/go-timekeeper/tests/integration/helper"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"github.com/gin-gonic/gin"
)

func TestSignInSuccess(t *testing.T) {
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
	testingVariables.Password = "password"

	if ok, _ := helper.SignUp(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign up user. Email: %s", testingVariables.Email)
	}

	if ok, _ := helper.SignIn(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign in user. Email: %s", testingVariables.Email)
	}
	t.Logf("✅ Successfully signed in user. Email: %s, token: %s", testingVariables.Email, testingVariables.AuthToken)
}

func TestSignInFailsWithIncorrectPassword(t *testing.T) {
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
	testingVariables.Password = "password"

	if ok, _ := helper.SignUp(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign up user. Email: %s", testingVariables.Email)
	}

	testingVariables.Password = "wrong-password"
	ok, resp := helper.SignIn(t, &client, server, testingVariables)
	if ok {
		t.Fatalf("❌ Sign in should have failed with incorrect password. Email: %s", testingVariables.Email)
	}

	var errorMessage = helper.ErrorMessage
	helper.DecodeJSON(t, resp.Body, &errorMessage)
	if errorMessage.ErrorMessage == service.ErrUserSignInFailed.Error() {
		t.Logf("✅ Sign in failed as expected with wrong password. Error: %s", errorMessage.ErrorMessage)
	} else {
		t.Fatalf("❌ Unexpected error message on failed sign in. Email: %s, got: %s", testingVariables.Email, errorMessage.ErrorMessage)
	}
}
