package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nanmenkaimak/github-gist/internal/gist/config"
	"io"
	"net/http"
	"time"
)

type UserTransport struct {
	config config.UserTransport
	client *http.Client
}

func NewTransport(config config.UserTransport, client *http.Client) *UserTransport {
	return &UserTransport{
		config: config,
		client: client,
	}
}

type GetUserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	RoleID    int       `json:"role_id"`
}

func (ut *UserTransport) GetUser(ctx context.Context, username string) (*GetUserResponse, error) {
	var response *GetUserResponse

	responseBody, err := ut.makeRequest(
		ctx, http.MethodGet, fmt.Sprintf("/api/user/%s", username), ut.config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to makeRequest err: %w", err)
	}

	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshall response err: %w", err)
	}
	return response, nil
}

func (ut *UserTransport) makeRequest(ctx context.Context, httpMethod string, endpoint string, timeout time.Duration) (b []byte, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	requestURl := ut.config.Host + endpoint

	req, err := http.NewRequestWithContext(ctx, httpMethod, requestURl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to NewRequestWithContext err: %v", err)
	}

	res, err := ut.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client making http request err: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body err: %w", err)
	}
	return body, nil
}
