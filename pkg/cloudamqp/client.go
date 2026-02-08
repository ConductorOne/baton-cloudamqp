package cloudamqp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultBaseURL = "https://customer.cloudamqp.com/api"

type Client struct {
	httpClient *http.Client
	Password   string
	baseURL    string
}

type UsersResponse = []User

func NewClient(httpClient *http.Client, password string, baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		httpClient: httpClient,
		Password:   password,
		baseURL:    baseURL,
	}
}

func (c *Client) usersBaseURL() string {
	return c.baseURL + "/team"
}

func (c *Client) userBaseURL() string {
	return c.baseURL + "/team/%s"
}

// GetUsers returns all users under the team account.
func (c *Client) GetUsers(ctx context.Context) ([]User, error) {
	var usersResponse UsersResponse

	err := c.get(
		ctx,
		c.usersBaseURL(),
		&usersResponse,
	)

	if err != nil {
		return nil, err
	}

	return usersResponse, nil
}

func NewUpdateUserRolePayload(role string) url.Values {
	payload := url.Values{}

	payload.Set("role", role)

	return payload
}

// UpdateUserRole updates role of provided user.
func (c *Client) UpdateUserRole(ctx context.Context, userId string, role string) error {
	err := c.put(
		ctx,
		fmt.Sprintf(c.userBaseURL(), userId),
		NewUpdateUserRolePayload(role),
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) get(ctx context.Context, urlAddress string, resourceResponse interface{}) error {
	return c.doRequest(ctx, urlAddress, http.MethodGet, nil, resourceResponse)
}

func (c *Client) put(ctx context.Context, urlAddress string, data url.Values, resourceResponse interface{}) error {
	return c.doRequest(ctx, urlAddress, http.MethodPut, data, resourceResponse)
}

func (c *Client) doRequest(
	ctx context.Context,
	urlAddress string,
	method string,
	data url.Values,
	resourceResponse interface{},
) error {
	var body strings.Reader

	if data != nil {
		encodedData := data.Encode()
		bodyReader := strings.NewReader(encodedData)
		body = *bodyReader
	}

	req, err := http.NewRequestWithContext(ctx, method, urlAddress, &body)
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
		return status.Error(codes.Code(uint32(rawResponse.StatusCode)), "Request failed") //nolint:gosec // StatusCode is always valid HTTP code
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
