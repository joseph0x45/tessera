package tessera

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserInfo struct {
	ID       string         `json:"id"`
	Username string         `json:"username"`
	Password string         `json:"password"`
	Metadata map[string]any `json:"metadata"` //NOT_IMPLEMENTED_YET
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

func getSessionID(res *http.Response) string {
	defer res.Body.Close()
	result := map[string]string{}
	json.NewDecoder(res.Body).Decode(&result)
	return result["session_id"]
}

func (c *TesseraClient) Register(username, password string) (string, error) {
	requestURL := fmt.Sprintf("%s/api/users/register", c.TesseraServerURL)
	body, err := json.Marshal(map[string]any{
		"app_id":   c.TesseraAppID,
		"username": username,
		"password": password,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	switch res.StatusCode {
	case http.StatusInternalServerError:
		return "", ErrServer
	case http.StatusBadRequest:
		return "", getErr(res)
	case http.StatusConflict:
		return "", getErr(res)
	case http.StatusCreated:
		return getSessionID(res), nil
	default:
		return "", fmt.Errorf("Tessera: Unexpected status %d", res.StatusCode)
	}
}

func (c *TesseraClient) Login(username, password string) (string, error) {
	requestURL := fmt.Sprintf("%s/api/users/login", c.TesseraServerURL)
	body, err := json.Marshal(map[string]any{
		"app_id":   c.TesseraAppID,
		"username": username,
		"password": password,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	switch res.StatusCode {
	case http.StatusInternalServerError:
		return "", ErrServer
	case http.StatusBadRequest:
		return "", getErr(res)
	case http.StatusOK:
		return getSessionID(res), nil
	default:
		return "", fmt.Errorf("Tessera: Unexpected status %d", res.StatusCode)
	}
}

func (c *TesseraClient) GetUserInfo(username string) (*UserInfo, error) {
	return nil, nil
}

func (c *TesseraClient) GetSessionUserInfo(sessionID string) (*UserInfo, error)

func (c *TesseraClient) Delete(username string) error {
	return nil
}

func (c *TesseraClient) DeleteSessionUser() error {
	return nil
}
