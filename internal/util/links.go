package util

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// CheckIfValidLink checks an URL if it is
// a valid link. So first, it checks if it
// starts with 'http'. If httpsOnly is set,
// it will be checked if the link starts with
// https. Then, a http request is executed
// to the link. If this fails or responds
// with an status code >= 400, the link is
// invalid and an error will be returned.
// The link is qualified as valid if the
// returned error is nil.
func CheckIfValidLink(url string, httpsOnly bool) error {
	if !strings.HasPrefix(url, "http") {
		return fmt.Errorf("invalud URL format")
	}

	if httpsOnly && !strings.HasPrefix(url, "https") {
		return fmt.Errorf("URL must be https")
	}

	res, err := http.DefaultClient.Get(url)
	if err != nil || res == nil {
		return fmt.Errorf("request to URL failed")
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("ULR request failed with status code %d", res.StatusCode)
	}

	return nil
}

// CheckIfValidShort checks if the short link is
// contained in the reserved string or if the short
// link does not single-result match the allowedRx.
// The short link is qualified as valid if the
// returned error is nil.
func CheckIfValidShort(sl, reserved string, allowedRx *regexp.Regexp) error {
	if strings.Contains(reserved, sl) {
		return fmt.Errorf("short link is reserved")
	}

	if len(allowedRx.FindAllString(sl, -1)) > 1 {
		return fmt.Errorf("unallowed characters")
	}

	return nil
}
