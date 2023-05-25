package cloudamqp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const BaseURL = "https://customer.cloudamqp.com/api"
const UsersBaseURL = BaseURL + "/team"

type Client struct {
	httpClient *http.Client
	Password   string
}

type UsersResponse = []User

func NewClient(httpClient *http.Client, password string) *Client {
	return &Client{
		httpClient: httpClient,
		Password:   password,
	}
}

// GetUsers returns all users under the team account.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var usersResponse UsersResponse

	err := c.doRequest(
		ctx,
		UsersBaseURL,
		&usersResponse,
	)

	if err != nil {
		return nil, err
	}

	return usersResponse, nil
}

func (c *Client) doRequest(
	ctx context.Context,
	urlAddress string,
	resourceResponse interface{},
) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlAddress, nil)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", constructAuth(c.Password))

	rawResponse, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return status.Error(codes.Code(rawResponse.StatusCode), "Request failed")
	}

	if err := json.NewDecoder(rawResponse.Body).Decode(&resourceResponse); err != nil {
		return err
	}

	return nil
}

func constructAuth(pass string) string {
	credentials := fmt.Sprintf("%s:%s", "", pass)
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	return "Basic " + encodedCredentials
}
