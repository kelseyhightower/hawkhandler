package hawkhandler

import (
	"net/http"

	"github.com/tent/hawk-go"
)

// hawkHandler holds the hawk configuration.
type hawkHandler struct {
	credentialsLookupFunc hawk.CredentialsLookupFunc
	handler               http.Handler
	whiteList             []string
}

// ServeHTTP
func (h *hawkHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, path := range h.whiteList {
		if r.URL.Path == path {
			h.handler.ServeHTTP(w, r)
			return
		}
	}
	auth, err := hawk.NewAuthFromRequest(r, h.credentialsLookupFunc, nil)
	if err != nil {
		w.Header().Set("WWW-Authenticate", "Hawk")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := auth.Valid(); err != nil {
		w.Header().Set("WWW-Authenticate", "Hawk")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	h.handler.ServeHTTP(w, r)
}

// HawkHandler wraps an http.Handler providing Hawk authentication to every request.
// Paths listed in the whitelist will bypass all authentication checks.
func HawkHandler(h http.Handler, whiteList []string, f hawk.CredentialsLookupFunc) http.Handler {
	return &hawkHandler{f, h, whiteList}
}
