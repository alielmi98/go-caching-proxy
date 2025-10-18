package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alielmi98/go-caching-proxy/pkg/cache"
	"github.com/alielmi98/go-caching-proxy/pkg/proxy"
)

func main() {
	port := flag.String("port", "8080", "Port to run the caching proxy server on")
	origin := flag.String("origin", "", "Origin URL to forward requests to")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache before starting the server")
	maxEntries := flag.Int("max-entries", 1000, "Maximum number of cache entries")
	cleanupInterval := flag.Duration("cleanup-interval", 5*time.Minute, "Interval for periodic cache cleanup")

	flag.Parse()

	cache := cache.NewCacheWithConfig(*maxEntries, *cleanupInterval)
	defer cache.StopCleanup()

	if *origin == "" {
		log.Fatal("Origin URL must be provided")
	}

	if *clearCache {
		cache.ClearCache()
		fmt.Println("Cache cleared")
	}

	p := proxy.NewProxy(*origin, cache)
	http.HandleFunc("/", p.ServeHTTP)

	// Add an endpoint to clear the cache
	http.HandleFunc("/clear-cache", func(w http.ResponseWriter, r *http.Request) {
		cache.ClearCache()
		fmt.Fprintln(w, "Cache cleared")
		log.Println("Cache cleared via /clear-cache endpoint")
	})

	// Add an endpoint to get cache statistics
	http.HandleFunc("/cache-stats", func(w http.ResponseWriter, r *http.Request) {
		current, max := cache.GetStats()
		stats := fmt.Sprintf("Cache Statistics:\nCurrent entries: %d\nMax entries: %d\nUsage: %.1f%%", 
			current, max, float64(current)/float64(max)*100)
		fmt.Fprintln(w, stats)
	})

	fmt.Printf("Starting caching proxy server on port %s, forwarding to %s\n", *port, *origin)
	fmt.Printf("Cache configuration: max entries=%d, cleanup interval=%v\n", *maxEntries, *cleanupInterval)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}