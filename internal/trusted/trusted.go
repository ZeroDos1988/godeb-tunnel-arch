package trusted

import (
	"context"
	"net/http"
	"time"
)

const Header = "X-GoDeb-Trusted"

// IsTrusted returns true when GET url responds 200 with X-GoDeb-Trusted: 1|true.
func IsTrusted(ctx context.Context, url string) bool {
	if url == "" {
		return false
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}
	client := &http.Client{Timeout: 3 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false
	}
	v := res.Header.Get(Header)
	return v == "1" || v == "true"
}
