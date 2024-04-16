package checks

import "net/http"

func IsJson(r *http.Request) bool {
	if r.Body == nil {
		return false
	}

	return r.Header.Get("Content-Type") == "application/json"
}
