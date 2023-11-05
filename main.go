//go:generate go run scripts/templ/main.go
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/unrolled/secure"
	"github.com/unrolled/secure/cspbuilder"

	"github.com/sjansen/bedrock/internal/templates/pages"
)

var subtitles = []string{
	"A New Hope",
	"The adventure begins...",
	"There and Back Again",
}

func main() {
	cspBuilder := cspbuilder.Builder{
		Directives: map[string][]string{
			cspbuilder.DefaultSrc: {"'self'"},
		},
	}
	secureMiddleware := secure.New(secure.Options{
		ContentSecurityPolicy: cspBuilder.MustBuild(),
		ContentTypeNosniff:    true,
		FrameDeny:             true,
		// TODO IsDevelopment:     true,
		// TODO PermissionsPolicy: ...,
	})

	r := chi.NewRouter()
	r.Use(secureMiddleware.Handler)

	r.Get("/",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idx := rand.Intn(len(subtitles))
			pages.Index("Hello, world!", subtitles[idx]).
				Render(r.Context(), w)
		}),
	)

	if os.Getenv("BEDROCK_ENV") == "localdev" {
		fmt.Println("Enabling /public/ handler...")
		fmt.Println(os.ReadDir("."))
		fmt.Println(os.ReadDir("public"))
		fs := http.FileServer(http.Dir("public"))
		r.Handle("/public/*", http.StripPrefix("/public/", fs))
	}

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}
