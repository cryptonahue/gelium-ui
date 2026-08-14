// Temporary smoke entrypoint — not committed.
package main

import (
	"log"
	"net/http"

	"geliumui/internal/app"
)

func main() {
	addr := "100.121.211.121:8787"
	log.Printf("Gelium UI serving on http://%s", addr)
	if err := http.ListenAndServe(addr, app.New()); err != nil {
		log.Fatal(err)
	}
}
