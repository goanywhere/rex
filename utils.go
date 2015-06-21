package rex

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}
