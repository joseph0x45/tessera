package tessera

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthResponse struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
}

type TesseraClient struct {
	TesseraServerURL string
	TesseraAppID     string
	httpClient       *http.Client
}

func Client(tesseraServerURL, tesseraAppID string) *TesseraClient {
	return &TesseraClient{
		TesseraServerURL: tesseraServerURL,
		TesseraAppID:     tesseraAppID,
		httpClient:       &http.Client{Timeout: time.Second * 10},
	}
}

func getErr(res *http.Response) error {
	defer res.Body.Close()
	errorMap := map[string]string{}
	json.NewDecoder(res.Body).Decode(&errorMap)
	errorMessage := errorMap["error"]
	return ErrMap[errorMessage]
}

func getAuthResponse(res *http.Response) *AuthResponse {
	defer res.Body.Close()
	authResponse := &AuthResponse{}
	json.NewDecoder(res.Body).Decode(authResponse)
	return authResponse
}

func (c *TesseraClient) Register(username, password string) (*AuthResponse, error) {
	requestURL := fmt.Sprintf("%s/api/users/register", c.TesseraServerURL)
	body, err := json.Marshal(map[string]any{
		"app_id":   c.TesseraAppID,
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case http.StatusInternalServerError:
		return nil, ErrServer
	case http.StatusBadRequest:
		return nil, getErr(res)
	case http.StatusConflict:
		return nil, getErr(res)
	case http.StatusCreated:
		return getAuthResponse(res), nil
	default:
		return nil, fmt.Errorf("Tessera: Unexpected status %d", res.StatusCode)
	}
}

func (c *TesseraClient) Login(username, password string) (*AuthResponse, error) {
	requestURL := fmt.Sprintf("%s/api/users/login", c.TesseraServerURL)
	body, err := json.Marshal(map[string]any{
		"app_id":   c.TesseraAppID,
		"username": username,
		"password": password,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case http.StatusInternalServerError:
		return nil, ErrServer
	case http.StatusBadRequest:
		return nil, getErr(res)
	case http.StatusOK:
		return getAuthResponse(res), nil
	default:
		return nil, fmt.Errorf("Tessera: Unexpected status %d", res.StatusCode)
	}
}

func (c *TesseraClient) Delete(username string) error {
	requestURL := fmt.Sprintf("%s/api/users", c.TesseraServerURL)
	body, err := json.Marshal(map[string]any{
		"app_id":   c.TesseraAppID,
		"username": username,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, requestURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case http.StatusInternalServerError:
		return ErrServer
	case http.StatusBadRequest:
		return getErr(res)
	case http.StatusOK:
		return nil
	default:
		return fmt.Errorf("Tessera: Unexpected status %d", res.StatusCode)
	}
}
