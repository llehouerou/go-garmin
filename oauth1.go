// oauth1.go - OAuth1 signature generation (HMAC-SHA1)
package garmin

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // OAuth1 requires SHA1, this is not a security concern for signing
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

// OAuth1Signer signs HTTP requests using OAuth1.0a HMAC-SHA1
type OAuth1Signer struct {
	ConsumerKey    string
	ConsumerSecret string
	Token          string
	TokenSecret    string
}

// Sign adds OAuth1 authorization header to the request
func (s *OAuth1Signer) Sign(req *http.Request) {
	params := s.buildOAuthParams()

	// Collect all parameters for signature base string
	allParams := url.Values{}
	for k, v := range params {
		allParams.Set(k, v)
	}

	// Add query parameters
	for k, vs := range req.URL.Query() {
		for _, v := range vs {
			allParams.Add(k, v)
		}
	}

	// Note: For POST requests with form data, body params should be included in the signature.
	// Since the body is a reader we can't peek, the caller should ensure body params
	// are included in query string for signature calculation if needed.

	// Generate signature
	signature := s.generateSignature(req.Method, req.URL, allParams)
	params["oauth_signature"] = signature

	// Build Authorization header
	req.Header.Set("Authorization", s.buildAuthHeader(params))
}

// buildOAuthParams creates the base OAuth parameters
func (s *OAuth1Signer) buildOAuthParams() map[string]string {
	params := map[string]string{
		"oauth_consumer_key":     s.ConsumerKey,
		"oauth_signature_method": "HMAC-SHA1",
		"oauth_timestamp":        strconv.FormatInt(time.Now().Unix(), 10),
		"oauth_nonce":            generateNonce(),
		"oauth_version":          "1.0",
	}

	if s.Token != "" {
		params["oauth_token"] = s.Token
	}

	return params
}

// generateSignature creates the HMAC-SHA1 signature
func (s *OAuth1Signer) generateSignature(method string, reqURL *url.URL, params url.Values) string {
	// Build base URL (scheme + host + path, no query string)
	baseURL := fmt.Sprintf("%s://%s%s", reqURL.Scheme, reqURL.Host, reqURL.Path)

	// Sort and encode parameters
	paramString := encodeParams(params)

	// Build signature base string
	baseString := fmt.Sprintf("%s&%s&%s",
		percentEncode(method),
		percentEncode(baseURL),
		percentEncode(paramString),
	)

	// Build signing key
	signingKey := fmt.Sprintf("%s&%s",
		percentEncode(s.ConsumerSecret),
		percentEncode(s.TokenSecret),
	)

	// Generate HMAC-SHA1
	h := hmac.New(sha1.New, []byte(signingKey))
	h.Write([]byte(baseString))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// buildAuthHeader constructs the OAuth Authorization header
func (s *OAuth1Signer) buildAuthHeader(params map[string]string) string {
	// Sort keys for consistent output
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build header parts
	parts := make([]string, 0, len(params))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, percentEncode(params[k])))
	}

	return "OAuth " + strings.Join(parts, ", ")
}

// encodeParams creates the normalized parameter string for signing
func encodeParams(params url.Values) string {
	// Get all keys and sort them
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build parameter string
	pairs := make([]string, 0, len(params))
	for _, k := range keys {
		vs := params[k]
		sort.Strings(vs)
		for _, v := range vs {
			pairs = append(pairs, fmt.Sprintf("%s=%s", percentEncode(k), percentEncode(v)))
		}
	}

	return strings.Join(pairs, "&")
}

// percentEncode performs OAuth-style percent encoding
// (RFC 3986 with specific requirements for OAuth)
func percentEncode(s string) string {
	var result strings.Builder
	for _, c := range []byte(s) {
		if isUnreserved(c) {
			result.WriteByte(c)
		} else {
			result.WriteByte('%')
			result.WriteByte(hexDigit(c >> 4))
			result.WriteByte(hexDigit(c & 0x0F))
		}
	}
	return result.String()
}

// hexDigit returns the uppercase hex digit for a value 0-15
func hexDigit(n byte) byte {
	if n < 10 {
		return '0' + n
	}
	return 'A' + n - 10
}

// isUnreserved checks if a character is unreserved per RFC 3986
func isUnreserved(c byte) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~'
}

// generateNonce creates a random nonce for OAuth requests
func generateNonce() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
