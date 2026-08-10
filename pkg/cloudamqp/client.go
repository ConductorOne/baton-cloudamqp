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
	//nolint:gosec,nolintlint // G117: legitimate field name, not a credential
	Password string
	baseURL  string
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

func (c *Client) teamInviteURL() string {
	return c.baseURL + "/team/invite"
}

func (c *Client) teamRemoveURL() string {
	return c.baseURL + "/team/remove"
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

// GetUserByEmail returns the team member with the given email, or a NotFound
// status error if no member matches. CloudAMQP has no per-user lookup endpoint,
// so this filters the full team list.
func (c *Client) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	users, err := c.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range users {
		if strings.EqualFold(users[i].Email, email) {
			return &users[i], nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "baton-cloudamqp: no team member with email %s", email)
}

// GetUserByID returns the team member with the given id, or a NotFound status
// error if no member matches.
func (c *Client) GetUserByID(ctx context.Context, id string) (*User, error) {
	users, err := c.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	for i := range users {
		if users[i].Id == id {
			return &users[i], nil
		}
	}

	return nil, status.Errorf(codes.NotFound, "baton-cloudamqp: no team member with id %s", id)
}

// InviteUser invites a new team member by email with the given role and
// optional instance tags. Each tag is sent as a separate form field so the
// server receives an array (e.g. tags[]=prod&tags[]=staging). The invitee must
// accept the emailed invitation before they appear as a team member.
func (c *Client) InviteUser(ctx context.Context, email string, role string, tags []string) error {
	payload := url.Values{}
	payload.Set("email", email)
	payload.Set("role", role)
	for _, t := range tags {
		payload.Add("tags[]", t)
	}
	return c.post(ctx, c.teamInviteURL(), payload, nil)
}

// RemoveUser removes a team member by email. CloudAMQP uses POST /team/remove
// (not an HTTP DELETE) and has no soft-disable, so this is a hard removal.
func (c *Client) RemoveUser(ctx context.Context, email string) error {
	payload := url.Values{}
	payload.Set("email", email)
	return c.post(ctx, c.teamRemoveURL(), payload, nil)
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

func (c *Client) post(ctx context.Context, urlAddress string, data url.Values, resourceResponse interface{}) error {
	return c.doRequest(ctx, urlAddress, http.MethodPost, data, resourceResponse)
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

	rawResponse, err := c.httpClient.Do(req) //nolint:gosec,nolintlint // G704: URL constructed from trusted config
	if err != nil {
		return err
	}

	defer rawResponse.Body.Close()

	if rawResponse.StatusCode >= 300 {
		return newAPIError(rawResponse.StatusCode, rawResponse.Body)
	}

	// Provisioning calls (invite/remove/update-role) pass a nil response target
	// and may return an empty body on success — skip decoding in that case.
	if resourceResponse == nil {
		return nil
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
