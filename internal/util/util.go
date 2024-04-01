package util

import "net/http"

func HasJsonContentType(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}
