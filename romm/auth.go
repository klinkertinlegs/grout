package romm

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	gaba "github.com/BrandonKowalski/gabagool/v2/pkg/gabagool"
)

type TokenExchangeRequest struct {
	Code string `json:"code"`
}

type TokenExchangeResponse struct {
	RawToken  string   `json:"raw_token"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	ExpiresAt string   `json:"expires_at"`
}

type CurrentUser struct {
	Username string `json:"username"`
}

func (c *Client) ValidateConnection() error {
	req, err := http.NewRequest("GET", c.baseURL+endpointHeartbeat, nil)
	if err != nil {
		return ClassifyError(fmt.Errorf("failed to create validation request: %w", err))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		classifiedErr := ClassifyError(fmt.Errorf("failed to connect: %w", err))

		shouldTryProtocolSwitch := !errors.Is(classifiedErr, ErrTimeout) &&
			!errors.Is(classifiedErr, ErrConnectionRefused) &&
			!errors.Is(classifiedErr, ErrInvalidHostname)

		if shouldTryProtocolSwitch {
			if protocolErr := c.tryAlternateProtocol(req.URL.Scheme, func(r *http.Response) bool {
				return r.StatusCode >= 200 && r.StatusCode < 300
			}); protocolErr != nil {
				return protocolErr
			}
		}

		return classifiedErr
	}
	defer resp.Body.Close()

	if req.URL.Scheme != resp.Request.URL.Scheme {
		return &ProtocolError{
			RequestedProtocol: req.URL.Scheme,
			CorrectProtocol:   resp.Request.URL.Scheme,
			Err:               ErrWrongProtocol,
		}
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	case resp.StatusCode >= 500:
		logResponseDebug("ValidateConnection: server error", resp)
		return &AuthError{
			StatusCode: resp.StatusCode,
			Message:    "Server error",
			Err:        ErrServerError,
		}
	default:
		logResponseDebug("ValidateConnection: unexpected status", resp)
		if protocolErr := c.tryAlternateProtocol(req.URL.Scheme, func(r *http.Response) bool {
			return r.StatusCode >= 200 && r.StatusCode < 300
		}); protocolErr != nil {
			return protocolErr
		}
		return fmt.Errorf("heartbeat check failed with status: %d", resp.StatusCode)
	}
}

func (c *Client) Login(username, password string) error {
	req, err := http.NewRequest("POST", c.baseURL+endpointLogin, nil)
	if err != nil {
		return ClassifyError(fmt.Errorf("failed to create login request: %w", err))
	}

	req.SetBasicAuth(username, password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ClassifyError(fmt.Errorf("failed to login: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		logResponseDebug("Login: failed", resp)
	}

	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		return nil
	case resp.StatusCode == 401:
		return &AuthError{StatusCode: 401, Message: "Invalid username or password", Err: ErrUnauthorized}
	case resp.StatusCode == 403:
		return &AuthError{StatusCode: 403, Message: "Access forbidden", Err: ErrForbidden}
	case resp.StatusCode >= 500:
		return &AuthError{StatusCode: resp.StatusCode, Message: "Server error", Err: ErrServerError}
	case resp.StatusCode == 405:
		if protocolErr := c.tryAlternateProtocolForLogin(req, username, password); protocolErr != nil {
			return protocolErr
		}
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	default:
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}
}

func ExchangeToken(baseURL string, code string, insecureSkipVerify bool) (*TokenExchangeResponse, error) {
	client := NewClient(baseURL, WithInsecureSkipVerify(insecureSkipVerify))
	var result TokenExchangeResponse
	err := client.doRequest("POST", endpointTokenExchange, nil, TokenExchangeRequest{Code: code}, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ValidateToken() error {
	var platforms []Platform
	return c.doRequest("GET", endpointPlatforms, nil, nil, &platforms)
}

func (c *Client) GetCurrentUser() (CurrentUser, error) {
	var user CurrentUser
	err := c.doRequest("GET", endpointCurrentUser, nil, nil, &user)
	return user, err
}

func switchProtocol(baseURL string) string {
	if len(baseURL) > 8 && baseURL[:8] == "https://" {
		return "http://" + baseURL[8:]
	}
	if len(baseURL) > 7 && baseURL[:7] == "http://" {
		return "https://" + baseURL[7:]
	}
	return baseURL
}

func (c *Client) tryAlternateProtocol(originalScheme string, isSuccess func(resp *http.Response) bool) *ProtocolError {
	switchedURL := switchProtocol(c.baseURL)
	if switchedURL == c.baseURL {
		return nil
	}

	testReq, err := http.NewRequest("GET", switchedURL+endpointHeartbeat, nil)
	if err != nil {
		return nil
	}

	testResp, err := c.httpClient.Do(testReq)
	if err != nil {
		return nil
	}
	defer testResp.Body.Close()

	if isSuccess(testResp) {
		return &ProtocolError{
			RequestedProtocol: originalScheme,
			CorrectProtocol:   testReq.URL.Scheme,
			Err:               ErrWrongProtocol,
		}
	}
	return nil
}

func (c *Client) tryAlternateProtocolForLogin(originalReq *http.Request, username, password string) *ProtocolError {
	switchedURL := switchProtocol(c.baseURL)
	if switchedURL == c.baseURL {
		return nil
	}

	testClient := NewClient(switchedURL, WithTimeout(c.httpClient.Timeout))
	testReq, err := http.NewRequest("POST", switchedURL+endpointLogin, nil)
	if err != nil {
		return nil
	}
	testReq.SetBasicAuth(username, password)

	testResp, err := testClient.httpClient.Do(testReq)
	if err != nil {
		return nil
	}
	defer testResp.Body.Close()

	if testResp.StatusCode != 405 && testResp.StatusCode < 500 {
		return &ProtocolError{
			RequestedProtocol: originalReq.URL.Scheme,
			CorrectProtocol:   testReq.URL.Scheme,
			Err:               ErrWrongProtocol,
		}
	}
	return nil
}

func logResponseDebug(label string, resp *http.Response) {
	logger := gaba.GetLogger()
	body, _ := io.ReadAll(resp.Body)

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	logger.Debug(label,
		"status", resp.StatusCode,
		"url", resp.Request.URL.String(),
		"headers", headers,
		"body", string(body),
	)
}
