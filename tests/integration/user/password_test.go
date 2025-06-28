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

func TestChangePasswordSuccess(t *testing.T) {
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

	newPassword := "new_P@ssw0rd"
	if ok, resp := helper.ChangePassword(t, &client, server, testingVariables, newPassword); !ok {
		var errorMessage = helper.ErrorMessage
		helper.DecodeJSON(t, resp.Body, &errorMessage)
		t.Fatalf("❌ Password change failed. Email: %s, error: %s", testingVariables.Email, errorMessage.ErrorMessage)
	}

	testingVariables.Password = newPassword
	if ok, _ := helper.SignIn(t, &client, server, testingVariables); !ok {
		t.Fatalf("❌ Failed to sign in with new password. Email: %s", testingVariables.Email)
	}
	t.Logf("✅ Successfully changed password and signed in. Email: %s", testingVariables.Email)
}

func TestChangePasswordFailsWithInvalidOldPassword(t *testing.T) {
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

	testingVariables.Password = "wrong_old_P@ssw0rd"
	ok, resp := helper.ChangePassword(t, &client, server, testingVariables, "new_P@ssw0rd")
	if ok {
		t.Fatalf("❌ Password change should have failed with invalid old password. Email: %s", testingVariables.Email)
	}

	var errorMessage = helper.ErrorMessage
	helper.DecodeJSON(t, resp.Body, &errorMessage)
	if errorMessage.ErrorMessage == service.ErrUserChangePasswordFailed.Error() {
		t.Logf("✅ Password change failed as expected with invalid old password. Error: %s", errorMessage.ErrorMessage)
	} else {
		t.Fatalf("❌ Unexpected error message. Email: %s, got: %s", testingVariables.Email, errorMessage.ErrorMessage)
	}
}

func TestChangePasswordFailsWithSameOldAndNewPassword(t *testing.T) {
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

	ok, resp := helper.ChangePassword(t, &client, server, testingVariables, testingVariables.Password)
	if ok {
		t.Fatalf("❌ Password change should have failed with same old and new password. Email: %s", testingVariables.Email)
	}

	var errorMessage = helper.ErrorMessage
	helper.DecodeJSON(t, resp.Body, &errorMessage)
	if errorMessage.ErrorMessage == "Old password must not be same as a new one" {
		t.Logf("✅ Password change failed as expected with same old and new password. Error: %s", errorMessage.ErrorMessage)
	} else {
		t.Fatalf("❌ Unexpected error message. Email: %s, got: %s", testingVariables.Email, errorMessage.ErrorMessage)
	}
}
