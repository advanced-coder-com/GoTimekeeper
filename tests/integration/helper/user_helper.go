package integration_test_helper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func SignIn(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
) (bool, *http.Response) {
	result := true
	signinBody := map[string]string{
		"email":    testVars.Email,
		"password": testVars.Password,
	}
	signinResp := DoPost(t, client, server.URL+"/api/user/signin", signinBody)
	if signinResp.StatusCode != http.StatusOK {
		result = false
		return result, signinResp
	}
	var signinData struct {
		Token string `json:"token"`
	}
	DecodeJSON(t, signinResp.Body, &signinData)

	if signinData.Token == "" {
		result = false
	}
	testVars.AuthToken = signinData.Token
	return result, signinResp
}

func SignUp(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
) (bool, *http.Response) {
	result := true
	signupBody := map[string]string{
		"email":    testVars.Email,
		"password": testVars.Password,
	}
	signupResp := DoPost(t, client, server.URL+"/api/user/signup", signupBody)

	if signupResp.StatusCode != http.StatusCreated && signupResp.StatusCode != http.StatusOK {
		result = false
	}
	return result, signupResp
}

func ChangePassword(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
	newPassword string,
) (bool, *http.Response) {
	result := true
	changePasswordBody := map[string]string{
		"old_password": testVars.Password,
		"new_password": newPassword,
	}
	signupResp := DoPutchAuth(t, client, server.URL+"/api/user/change-password", changePasswordBody, testVars.AuthToken)

	if signupResp.StatusCode != http.StatusOK {
		result = false
	}
	testVars.Password = newPassword
	return result, signupResp
}

func DeleteUser(
	t *testing.T,
	client *http.Client,
	server *httptest.Server,
	testVars *TestingContext,
) (bool, *http.Response) {
	result := true
	signupResp := DoDeleteAuth(t, client, server.URL+"/api/user/delete", nil, testVars.AuthToken)

	if signupResp.StatusCode != http.StatusOK {
		result = false
	}
	return result, signupResp
}
