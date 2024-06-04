package writer

import (
	"fmt"
	"net/http"
	"time"
)

func SetCacheControlSWR(w http.ResponseWriter, maxAge time.Duration, swr time.Duration) {
	header := fmt.Sprintf("max-age=%.0f, stale-while-revalidate=%.0f", maxAge.Seconds(), swr.Seconds())
	w.Header().Set("Cache-Control", header)
}
