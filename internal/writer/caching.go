package writer

import (
	"fmt"
	"net/http"
	"time"
)

// The stale-while-revalidate response directive indicates that the cache
// could reuse a stale response while it revalidates it to a cache.
//
// Revalidation will make the cache be fresh again, so it appears to clients
// that it was always fresh during that period — effectively hiding the latency
// penalty of revalidation from them.
//
// If no request happened during that period, the cache became stale and the
// next request will revalidate normally.
func SetCacheControlSWR(w http.ResponseWriter, maxAge time.Duration, swr time.Duration) {
	header := fmt.Sprintf("max-age=%.0f, stale-while-revalidate=%.0f", maxAge.Seconds(), swr.Seconds())
	w.Header().Set("Cache-Control", header)
}
