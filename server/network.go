package server

import (
	"log"
	"net/http"
)

type Adapter func(http.Handler) http.Handler

func Logging(l *log.Logger) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Println(r.Method, r.URL.Path)
			h.ServeHTTP(w, r)
		})
	}
}

func WithHeader(key, value string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		})
	}
}

// Override method by checking its header
// One thing you can do is to "tunnel" HTTP Methods inside another HTTP Header.
// Basically you have a header that says "No, seriously, I know I got here via
// a POST, but use this one instead."
// You would still POST, but then you'd have
// "X-HTTP-Method-Override:PUT" as a header.
func SupportXHTTPMethodOverride() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := r.Header.Get("X-HTTP-Method-Override")
			if len(m) > 0 {
				r.Method = m
			}
			h.ServeHTTP(w, r)
		})
	}
}

func RecoverHandler() Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("panic: %+v", err)
					http.Error(w, http.StatusText(500), 500)
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		// adapt original handler with multiple outer functions
		h = adapter(h)
	}
	return h
}
