package tessera_test

import (
	"errors"
	"log"
	"testing"

	"github.com/joseph0x45/tessera"
)

const MAIN_APP_ID = "I5paDtMQOjfdHaO9VnKsc"
const SECONDARY_APP_ID = "5xkeKo7B8NGkk5VPnXxzq"

func TestCreateUserEmptyParameters(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := ""
	client := tessera.Client(serverURL, appID)
	username := ""
	password := ""
	_, err := client.Register(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrRequiredFieldMissing) {
		log.Fatal("Expected ErrRequiredFieldMissing but got", err.Error())
	}
}

func TestCreateUserInvalidAppID(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := "invalid"
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "test"
	_, err := client.Register(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrAppNotFound) {
		log.Fatal("Expected ErrAppNotFound but got", err.Error())
	}
}

func TestCreateUser(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := MAIN_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "test"
	authResponse, err := client.Register(username, password)
	if err != nil {
		log.Fatal("Expected err to be nil but got", err.Error())
	}
	if authResponse == nil {
		log.Fatal("Expected sessionID not to be nil")
	}
}

func TestCreateUserSameUsername(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := MAIN_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "test"
	_, err := client.Register(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil but got")
	}
	if !errors.Is(err, tessera.ErrUserExistsInApp) {
		log.Fatal("Expected ErrUserExistsInApp but got", err.Error())
	}
}

func TestLoginEmptyParameters(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := ""
	client := tessera.Client(serverURL, appID)
	username := ""
	password := ""
	_, err := client.Login(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrRequiredFieldMissing) {
		log.Fatal("Expected ErrRequiredFieldMissing but got", err.Error())
	}
}

func TestLoginInvalidAppID(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := "invalid"
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "test"
	_, err := client.Login(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrAppNotFound) {
		log.Fatal("Expected ErrAppNotFound but got", err.Error())
	}
}

func TestLoginWrongAppID(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := SECONDARY_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "test"
	_, err := client.Login(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrUserNotFound) {
		log.Fatal("Expected ErrUserNotFound but got", err.Error())
	}
}

func TestLoginWrongUsername(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := MAIN_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "teste"
	password := "test"
	_, err := client.Login(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrUserNotFound) {
		log.Fatal("Expected ErrUserNotFound but got", err.Error())
	}
}

func TestLoginWrongPassword(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := MAIN_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "testee"
	_, err := client.Login(username, password)
	if err == nil {
		log.Fatal("Expected err to be non nil")
	}
	if !errors.Is(err, tessera.ErrInvalidPassword) {
		log.Fatal("Expected ErrInvalidPassword but got", err.Error())
	}
}

func TestLogin(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := MAIN_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "test"
	password := "test"
	authResponse, err := client.Login(username, password)
	if err != nil {
		log.Fatal("Expected err to be nil")
	}
	if authResponse == nil {
		log.Fatal("Expected sessionID not to be empty")
	}
}

func TestDeleteUser(t *testing.T) {
	serverURL := "http://localhost:8080"
	appID := MAIN_APP_ID
	client := tessera.Client(serverURL, appID)
	username := "test"
	err := client.Delete(username)
	if err != nil {
		log.Fatal("Expected err to be nil")
	}
}
