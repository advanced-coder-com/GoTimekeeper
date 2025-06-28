package user_test

import (
	"context"
	"fmt"
	helper "github.com/advanced-coder-com/go-timekeeper/tests/integration/helper"
	"testing"

	"github.com/advanced-coder-com/go-timekeeper/internal/repository"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/advanced-coder-com/go-timekeeper/internal/db"
	"github.com/advanced-coder-com/go-timekeeper/internal/router"
	"github.com/gin-gonic/gin"
)

func TestSignUpSuccess(t *testing.T) {
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

	userRepo := repository.NewUserRepository()
	user, err := userRepo.GetByEmail(context.Background(), testingVariables.Email)
	if err != nil {
		t.Fatalf("❌ Could not find created user by email: %s. Error: %v", testingVariables.Email, err)
	}
	t.Logf("✅ Successfully created user. Email: %s, ID: %s", user.Email, user.ID)
}

func TestSignUpFailsWithDuplicateEmail(t *testing.T) {
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

	ok, resp := helper.SignUp(t, &client, server, testingVariables)
	if ok {
		t.Fatalf("❌ Duplicate user creation should have failed. Email: %s", testingVariables.Email)
	}

	var errorMessage = helper.ErrorMessage
	helper.DecodeJSON(t, resp.Body, &errorMessage)
	errMessage := fmt.Sprintf("User with email %s already exists", testingVariables.Email)
	if errorMessage.ErrorMessage == errMessage {
		t.Logf("✅ Sign up failed as expected with duplicate email. Error: %s", errorMessage.ErrorMessage)
	} else {
		t.Fatalf("❌ Unexpected error message on duplicate sign up. Email: %s, got: %s", testingVariables.Email, errorMessage.ErrorMessage)
	}
}
