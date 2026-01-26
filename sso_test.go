package garmin

import (
	"net/http"
	"net/url"
	"testing"
)

func TestExtractCSRF(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid CSRF token",
			html:    `<input name="_csrf" value="abc123def">`,
			want:    "abc123def",
			wantErr: false,
		},
		{
			name:    "CSRF with extra whitespace",
			html:    `<input name="_csrf"   value="token-with-spaces">`,
			want:    "token-with-spaces",
			wantErr: false,
		},
		{
			name:    "CSRF in larger HTML",
			html:    `<html><head></head><body><form><input type="hidden" name="_csrf" value="xyz789"></form></body></html>`,
			want:    "xyz789",
			wantErr: false,
		},
		{
			name:    "no CSRF token",
			html:    `<html><head></head><body></body></html>`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractCSRF(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractCSRF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractCSRF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name:    "simple title",
			html:    `<html><head><title>Success</title></head></html>`,
			want:    "Success",
			wantErr: false,
		},
		{
			name:    "MFA title",
			html:    `<html><head><title>GARMIN > MFA Challenge</title></head></html>`,
			want:    "GARMIN > MFA Challenge",
			wantErr: false,
		},
		{
			name:    "no title",
			html:    `<html><head></head><body></body></html>`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTitle(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTicket(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		want    string
		wantErr bool
	}{
		{
			name:    "valid ticket",
			html:    `<a href="https://sso.garmin.com/sso/embed?ticket=ST-123-abc">Continue</a>`,
			want:    "ST-123-abc",
			wantErr: false,
		},
		{
			name:    "ticket in script",
			html:    `<script>window.location="https://sso.garmin.com/sso/embed?ticket=ST-456-def";</script>`,
			want:    "ST-456-def",
			wantErr: false,
		},
		{
			name:    "no ticket",
			html:    `<html><body>No ticket here</body></html>`,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractTicket(tt.html)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractTicket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOAuth1Signer_PercentEncode(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"abc", "abc"},
		{"ABC", "ABC"},
		{"123", "123"},
		{"abc-._~", "abc-._~"},
		{"hello world", "hello%20world"},
		{"a=b&c=d", "a%3Db%26c%3Dd"},
		{"100%", "100%25"},
		{"https://example.com/path", "https%3A%2F%2Fexample.com%2Fpath"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := percentEncode(tt.input)
			if got != tt.want {
				t.Errorf("percentEncode(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestOAuth1Signer_EncodeParams(t *testing.T) {
	params := url.Values{
		"oauth_consumer_key": {"key123"},
		"oauth_nonce":        {"nonce456"},
		"b":                  {"2"},
		"a":                  {"1"},
	}

	result := encodeParams(params)

	// Parameters should be sorted alphabetically
	expected := "a=1&b=2&oauth_consumer_key=key123&oauth_nonce=nonce456"
	if result != expected {
		t.Errorf("encodeParams() = %q, want %q", result, expected)
	}
}

func TestOAuth1Signer_GenerateSignature(t *testing.T) {
	// Test case from OAuth 1.0 spec examples
	signer := &OAuth1Signer{
		ConsumerKey:    "dpf43f3p2l4k3l03",
		ConsumerSecret: "kd94hf93k423kf44",
		Token:          "nnch734d00sl2jdk",
		TokenSecret:    "pfkkdhi9sl3r4s00",
	}

	reqURL, _ := url.Parse("http://photos.example.net/photos")

	params := url.Values{
		"oauth_consumer_key":     {"dpf43f3p2l4k3l03"},
		"oauth_token":            {"nnch734d00sl2jdk"},
		"oauth_signature_method": {"HMAC-SHA1"},
		"oauth_timestamp":        {"1191242096"},
		"oauth_nonce":            {"kllo9940pd9333jh"},
		"oauth_version":          {"1.0"},
		"file":                   {"vacation.jpg"},
		"size":                   {"original"},
	}

	signature := signer.generateSignature("GET", reqURL, params)

	// The expected signature is from the OAuth 1.0 spec
	expected := "tR3+Ty81lMeYAr/Fid0kMTYa/WM="
	if signature != expected {
		t.Errorf("generateSignature() = %q, want %q", signature, expected)
	}
}

func TestOAuth1Signer_BuildAuthHeader(t *testing.T) {
	signer := &OAuth1Signer{}

	params := map[string]string{
		"oauth_consumer_key":     "key123",
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_signature":        "sig%3D",
	}

	header := signer.buildAuthHeader(params)

	// Check header format
	if !startsWith(header, "OAuth ") {
		t.Errorf("buildAuthHeader() should start with 'OAuth ', got %q", header)
	}

	// Verify all params are present and properly formatted
	if !contains(header, `oauth_consumer_key="key123"`) {
		t.Error("buildAuthHeader() missing oauth_consumer_key")
	}
	if !contains(header, `oauth_signature_method="HMAC-SHA1"`) {
		t.Error("buildAuthHeader() missing oauth_signature_method")
	}
}

func TestOAuth1Signer_Sign(t *testing.T) {
	signer := &OAuth1Signer{
		ConsumerKey:    "consumer-key",
		ConsumerSecret: "consumer-secret",
		Token:          "token",
		TokenSecret:    "token-secret",
	}

	req, _ := http.NewRequest(http.MethodGet, "https://api.example.com/resource?param=value", http.NoBody)
	signer.Sign(req)

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		t.Error("Sign() did not set Authorization header")
	}

	if !startsWith(authHeader, "OAuth ") {
		t.Errorf("Sign() Authorization header should start with 'OAuth ', got %q", authHeader)
	}

	// Verify required OAuth params are present
	requiredParams := []string{
		"oauth_consumer_key",
		"oauth_signature_method",
		"oauth_timestamp",
		"oauth_nonce",
		"oauth_version",
		"oauth_token",
		"oauth_signature",
	}

	for _, param := range requiredParams {
		if !contains(authHeader, param) {
			t.Errorf("Sign() Authorization header missing %s", param)
		}
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce1 := generateNonce()
	nonce2 := generateNonce()

	if nonce1 == "" {
		t.Error("generateNonce() returned empty string")
	}

	if len(nonce1) != 32 { // 16 bytes = 32 hex chars
		t.Errorf("generateNonce() length = %d, want 32", len(nonce1))
	}

	if nonce1 == nonce2 {
		t.Error("generateNonce() returned same value twice")
	}
}

// Helper functions for tests
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s != "" && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
