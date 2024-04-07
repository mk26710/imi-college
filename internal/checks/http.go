package checks

import "net/http"

func IsJson(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}
