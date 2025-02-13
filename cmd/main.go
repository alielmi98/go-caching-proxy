package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/alielmi98/go-caching-proxy/pkg/cache"
	"github.com/alielmi98/go-caching-proxy/pkg/proxy"
)

func main() {
	port := flag.String("port", "8080", "Port to run the caching proxy server on")
	origin := flag.String("origin", "", "Origin URL to forward requests to")
	clearCache := flag.Bool("clear-cache", false, "Clear the cache before starting the server")

	flag.Parse()

	cache := cache.NewCache()

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

	fmt.Printf("Starting caching proxy server on port %s, forwarding to %s\n", *port, *origin)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
