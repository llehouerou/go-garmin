// sso.go - Garmin SSO authentication flow
package garmin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	csrfRE   = regexp.MustCompile(`name="_csrf"\s+value="([^"]+)"`)
	titleRE  = regexp.MustCompile(`<title>([^<]+)</title>`)
	ticketRE = regexp.MustCompile(`embed\?ticket=([^"]+)"`)
)

var (
	ErrCSRFNotFound   = errors.New("garmin: CSRF token not found")
	ErrTitleNotFound  = errors.New("garmin: title not found")
	ErrTicketNotFound = errors.New("garmin: ticket not found")
	ErrLoginFailed    = errors.New("garmin: login failed")
)

func extractCSRF(html string) (string, error) {
	m := csrfRE.FindStringSubmatch(html)
	if m == nil {
		return "", ErrCSRFNotFound
	}
	return m[1], nil
}

func extractTitle(html string) (string, error) {
	m := titleRE.FindStringSubmatch(html)
	if m == nil {
		return "", ErrTitleNotFound
	}
	return m[1], nil
}

func extractTicket(html string) (string, error) {
	m := ticketRE.FindStringSubmatch(html)
	if m == nil {
		return "", ErrTicketNotFound
	}
	return m[1], nil
}

// ssoClient handles the SSO authentication flow
type ssoClient struct {
	httpClient *http.Client
	domain     string
	timeout    time.Duration
}

// newSSOClient creates an SSO client. If baseClient is provided, its transport
// is reused (for VCR testing), otherwise a new client is created.
func newSSOClient(domain string, timeout time.Duration, baseClient *http.Client) (*ssoClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	var httpClient *http.Client
	if baseClient != nil {
		// Reuse the base client's transport (for VCR) but add our own cookie jar
		httpClient = &http.Client{
			Transport: baseClient.Transport,
			Jar:       jar,
			Timeout:   timeout,
		}
	} else {
		httpClient = &http.Client{
			Jar:     jar,
			Timeout: timeout,
		}
	}

	return &ssoClient{
		httpClient: httpClient,
		domain:     domain,
		timeout:    timeout,
	}, nil
}

// ssoLogin performs the full SSO login flow and returns OAuth1 and OAuth2 tokens
func (c *Client) ssoLogin(ctx context.Context, email, password string) error {
	sso, err := newSSOClient(c.opts.Domain, 30*time.Second, c.transport.client)
	if err != nil {
		return err
	}

	// Step 1-4: Perform SSO authentication and get ticket
	ticket, err := sso.authenticate(ctx, email, password, c.opts.MFAHandler)
	if err != nil {
		return err
	}

	// Step 5: Fetch OAuth consumer credentials
	consumer, err := fetchOAuthConsumer(ctx, sso.httpClient)
	if err != nil {
		return fmt.Errorf("failed to fetch OAuth consumer: %w", err)
	}

	// Step 6: Get OAuth1 token using ticket
	oauth1Token, err := sso.getOAuth1Token(ctx, ticket, consumer)
	if err != nil {
		return fmt.Errorf("failed to get OAuth1 token: %w", err)
	}

	// Step 7: Exchange OAuth1 for OAuth2 token
	oauth2Token, err := sso.exchangeOAuth1ForOAuth2(ctx, oauth1Token, consumer)
	if err != nil {
		return fmt.Errorf("failed to exchange for OAuth2 token: %w", err)
	}

	// Update client auth state
	c.auth.OAuth1Token = oauth1Token.Token
	c.auth.OAuth1Secret = oauth1Token.Secret
	c.auth.MFAToken = oauth1Token.MFAToken
	c.auth.OAuth2AccessToken = oauth2Token.AccessToken
	c.auth.OAuth2RefreshToken = oauth2Token.RefreshToken
	c.auth.OAuth2Expiry = oauth2Token.Expiry
	c.auth.OAuth2Scope = oauth2Token.Scope
	c.auth.Domain = c.opts.Domain

	return nil
}

// authenticate performs steps 1-4 of the SSO flow
func (s *ssoClient) authenticate(ctx context.Context, email, password string, mfaHandler func() (string, error)) (string, error) {
	ssoBase := fmt.Sprintf("https://sso.%s/sso", s.domain)
	ssoEmbed := ssoBase + "/embed"

	// Build query params
	embedParams := url.Values{
		"id":          {"gauth-widget"},
		"embedWidget": {"true"},
		"gauthHost":   {ssoBase},
	}

	signinParams := url.Values{
		"id":                              {"gauth-widget"},
		"embedWidget":                     {"true"},
		"gauthHost":                       {ssoEmbed},
		"service":                         {ssoEmbed},
		"source":                          {ssoEmbed},
		"redirectAfterAccountLoginUrl":    {ssoEmbed},
		"redirectAfterAccountCreationUrl": {ssoEmbed},
	}

	// Step 1: Set cookies - GET embed page
	embedURL := ssoEmbed + "?" + embedParams.Encode()
	if _, err := s.doGet(ctx, embedURL, ""); err != nil {
		return "", fmt.Errorf("failed to set cookies: %w", err)
	}

	// Step 2: Get CSRF token - GET signin page
	signinURL := ssoBase + "/signin?" + signinParams.Encode()
	signinHTML, err := s.doGet(ctx, signinURL, embedURL)
	if err != nil {
		return "", fmt.Errorf("failed to get signin page: %w", err)
	}

	csrf, err := extractCSRF(signinHTML)
	if err != nil {
		return "", err
	}

	// Step 3: Submit credentials - POST to signin
	formData := url.Values{
		"username": {email},
		"password": {password},
		"embed":    {"true"},
		"_csrf":    {csrf},
	}

	responseHTML, err := s.doPost(ctx, signinURL, signinURL, formData)
	if err != nil {
		return "", fmt.Errorf("failed to submit credentials: %w", err)
	}

	title, err := extractTitle(responseHTML)
	if err != nil {
		return "", err
	}

	// Step 4: Handle MFA if required
	if strings.Contains(title, "MFA") {
		responseHTML, title, err = s.handleMFA(ctx, responseHTML, ssoBase, signinURL, signinParams, mfaHandler)
		if err != nil {
			return "", err
		}
	}

	// Verify success
	if title != "Success" {
		return "", fmt.Errorf("%w: unexpected title %q", ErrLoginFailed, title)
	}

	// Extract ticket from success page
	ticket, err := extractTicket(responseHTML)
	if err != nil {
		return "", err
	}

	return ticket, nil
}

// handleMFA processes the MFA challenge and returns the updated response HTML and title
func (s *ssoClient) handleMFA(
	ctx context.Context,
	responseHTML, ssoBase, signinURL string,
	signinParams url.Values,
	mfaHandler func() (string, error),
) (newHTML, title string, err error) {
	if mfaHandler == nil {
		return "", "", ErrMFARequired
	}

	// Get CSRF for MFA form
	mfaCSRF, err := extractCSRF(responseHTML)
	if err != nil {
		return "", "", err
	}

	mfaCode, err := mfaHandler()
	if err != nil {
		return "", "", fmt.Errorf("MFA handler failed: %w", err)
	}

	mfaURL := ssoBase + "/verifyMFA/loginEnterMfaCode?" + signinParams.Encode()
	mfaFormData := url.Values{
		"mfa-code": {mfaCode},
		"embed":    {"true"},
		"_csrf":    {mfaCSRF},
		"fromPage": {"setupEnterMfaCode"},
	}

	newHTML, err = s.doPost(ctx, mfaURL, signinURL, mfaFormData)
	if err != nil {
		return "", "", fmt.Errorf("failed to submit MFA code: %w", err)
	}

	title, err = extractTitle(newHTML)
	if err != nil {
		return "", "", err
	}

	return newHTML, title, nil
}

// OAuth1Token represents an OAuth1 token from the preauthorized endpoint
type OAuth1Token struct {
	Token    string
	Secret   string
	MFAToken string
}

// OAuth2Token represents an OAuth2 token from the exchange endpoint
type OAuth2Token struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
	Scope        string
	TokenType    string
}

// getOAuth1Token exchanges the SSO ticket for an OAuth1 token
func (s *ssoClient) getOAuth1Token(ctx context.Context, ticket string, consumer *oauthConsumer) (*OAuth1Token, error) {
	baseURL := fmt.Sprintf("https://connectapi.%s/oauth-service/oauth/", s.domain)
	loginURL := fmt.Sprintf("https://sso.%s/sso/embed", s.domain)

	tokenURL := fmt.Sprintf("%spreauthorized?ticket=%s&login-url=%s&accepts-mfa-tokens=true",
		baseURL, url.QueryEscape(ticket), url.QueryEscape(loginURL))

	// Create OAuth1 signer for this request (no token yet, just consumer credentials)
	signer := &OAuth1Signer{
		ConsumerKey:    consumer.Key,
		ConsumerSecret: consumer.Secret,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tokenURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "com.garmin.android.apps.connectmobile")
	signer.Sign(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OAuth1 token request failed: %s - %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse OAuth1 token from response (URL-encoded format)
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse OAuth1 response: %w", err)
	}

	return &OAuth1Token{
		Token:    values.Get("oauth_token"),
		Secret:   values.Get("oauth_token_secret"),
		MFAToken: values.Get("mfa_token"),
	}, nil
}

// exchangeOAuth1ForOAuth2 exchanges an OAuth1 token for an OAuth2 token
func (s *ssoClient) exchangeOAuth1ForOAuth2(ctx context.Context, oauth1 *OAuth1Token, consumer *oauthConsumer) (*OAuth2Token, error) {
	baseURL := fmt.Sprintf("https://connectapi.%s/oauth-service/oauth/", s.domain)
	exchangeURL := baseURL + "exchange/user/2.0"

	// Create OAuth1 signer with both consumer and token credentials
	signer := &OAuth1Signer{
		ConsumerKey:    consumer.Key,
		ConsumerSecret: consumer.Secret,
		Token:          oauth1.Token,
		TokenSecret:    oauth1.Secret,
	}

	// Build request body
	var body io.Reader
	if oauth1.MFAToken != "" {
		formData := url.Values{"mfa_token": {oauth1.MFAToken}}
		body = strings.NewReader(formData.Encode())
	} else {
		body = http.NoBody
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, exchangeURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "com.garmin.android.apps.connectmobile")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	signer.Sign(req)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OAuth2 exchange failed: %s - %s", resp.Status, string(respBody))
	}

	// Parse JSON response
	var tokenResp struct {
		Scope                 string `json:"scope"`
		JTI                   string `json:"jti"`
		AccessToken           string `json:"access_token"`
		TokenType             string `json:"token_type"`
		RefreshToken          string `json:"refresh_token"`
		ExpiresIn             int64  `json:"expires_in"`
		RefreshTokenExpiresIn int64  `json:"refresh_token_expires_in"`
	}

	if err := readJSON(resp, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse OAuth2 response: %w", err)
	}

	return &OAuth2Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
		Scope:        tokenResp.Scope,
		TokenType:    tokenResp.TokenType,
	}, nil
}

// doGet performs a GET request and returns the response body as string
func (s *ssoClient) doGet(ctx context.Context, reqURL, referer string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, http.NoBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "com.garmin.android.apps.connectmobile")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// doPost performs a POST request with form data and returns the response body
func (s *ssoClient) doPost(ctx context.Context, reqURL, referer string, data url.Values) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "com.garmin.android.apps.connectmobile")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Login authenticates the client with email and password
func (c *Client) Login(ctx context.Context, email, password string) error {
	return c.ssoLogin(ctx, email, password)
}

// readJSON decodes the JSON body from an HTTP response
func readJSON(resp *http.Response, v any) error {
	return json.NewDecoder(resp.Body).Decode(v)
}
