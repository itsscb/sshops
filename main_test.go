package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Run tests
	exitCode := m.Run()

	// Exit with the appropriate code
	os.Exit(exitCode)
}

func TestEnableHandler(t *testing.T) {
	// Set up the test environment
	os.Setenv("AUTH_PASSWORD", "P@55w0rd")
	req, err := http.NewRequest("POST", "/enable", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer P@55w0rd")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(enableHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse := "Enabled\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q but got %q", expectedResponse, rr.Body.String())
	}
}

func TestDisableHandler(t *testing.T) {
	// Set up the test environment
	os.Setenv("AUTH_PASSWORD", "P@55w0rd")
	req, err := http.NewRequest("POST", "/disable", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer P@55w0rd")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(disableHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse := "Disabled\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q but got %q", expectedResponse, rr.Body.String())
	}
}

func TestStateHandler(t *testing.T) {
	// Set up the test environment
	os.Setenv("AUTH_PASSWORD", "P@55w0rd")
	req, err := http.NewRequest("GET", "/state", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer P@55w0rd")

	// Mock the exec.Command execution to return "0" as output
	_ = func(command string, args ...string) *exec.Cmd {
		return exec.Command("echo", "0")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(stateHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	// Check the response body
	expectedResponse := "0\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q but got %q", expectedResponse, rr.Body.String())
	}

	// Reset the execCommand to its original value
	// execCommand = exec.Command
}


func TestUnauthorizedRequest(t *testing.T) {
	// Set up the test environment
	os.Setenv("AUTH_PASSWORD", "P@55w0rd")
	req, err := http.NewRequest("POST", "/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(enableHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, rr.Code)
	}

	// Check the response body
	expectedResponse := "Unauthorized\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q but got %q", expectedResponse, rr.Body.String())
	}
}

func TestMissingAuthorizationHeader(t *testing.T) {
	// Set up the test environment
	os.Setenv("AUTH_PASSWORD", "P@55w0rd")
	req, err := http.NewRequest("POST", "/enable", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(enableHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d but got %d", http.StatusUnauthorized, rr.Code)
	}

	// Check the response body
	expectedResponse := "Unauthorized\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response body %q but got %q", expectedResponse, rr.Body.String())
	}
}

func TestServerBehavior(t *testing.T) {
	// Set up the test environment
	os.Setenv("AUTH_PASSWORD", "P@55w0rd")

	// Test enable endpoint
	reqEnable, _ := http.NewRequest("POST", "/enable", nil)
	reqEnable.Header.Set("Authorization", "Bearer P@55w0rd")
	rrEnable := httptest.NewRecorder()
	handlerEnable := http.HandlerFunc(enableHandler)
	handlerEnable.ServeHTTP(rrEnable, reqEnable)
	assert.Equal(t, http.StatusOK, rrEnable.Code)
	assert.Equal(t, "Enabled\n", rrEnable.Body.String())

	// Test disable endpoint
	reqDisable, _ := http.NewRequest("POST", "/disable", nil)
	reqDisable.Header.Set("Authorization", "Bearer P@55w0rd")
	rrDisable := httptest.NewRecorder()
	handlerDisable := http.HandlerFunc(disableHandler)
	handlerDisable.ServeHTTP(rrDisable, reqDisable)
	assert.Equal(t, http.StatusOK, rrDisable.Code)
	assert.Equal(t, "Disabled\n", rrDisable.Body.String())

	// Test state endpoint
	reqState, _ := http.NewRequest("GET", "/state", nil)
	reqState.Header.Set("Authorization", "Bearer P@55w0rd")
	rrState := httptest.NewRecorder()
	handlerState := http.HandlerFunc(stateHandler)
	handlerState.ServeHTTP(rrState, reqState)
	assert.Equal(t, http.StatusOK, rrState.Code)
	assert.Contains(t, rrState.Body.String(), "0")

	// Clean up the test environment
	os.Unsetenv("AUTH_PASSWORD")
}