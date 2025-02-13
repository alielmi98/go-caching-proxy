package proxy

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alielmi98/go-caching-proxy/pkg/cache"
)

// Proxy represents the caching proxy server.
type Proxy struct {
	origin string
	cache  *cache.Cache
}

// NewProxy creates a new Proxy instance.
func NewProxy(origin string, cache *cache.Cache) *Proxy {
	return &Proxy{
		origin: origin,
		cache:  cache,
	}
}

// ServeHTTP handles incoming HTTP requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheKey := r.URL.Path
	if cachedResponse, found := p.cache.Get(cacheKey); found {
		w.Header().Set("X-Cache", "HIT")
		w.Write(cachedResponse)
		log.Printf("Served from cache: %s", cacheKey)
		return
	}

	resp, err := http.Get(p.origin + r.URL.Path)
	if err != nil {
		http.Error(w, "Error fetching from origin", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	p.cache.Set(cacheKey, body, 5*time.Minute)
	w.Header().Set("X-Cache", "MISS")
	w.Write(body)
	log.Printf("Fetched from origin and cached: %s", cacheKey)
}
