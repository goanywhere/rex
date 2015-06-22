package rex

import (
	"net/http"

	"github.com/gorilla/schema"
)

var Schema = schema.NewDecoder()

type Validator interface {
	Validate() error
}

// ParseForm parsed the raw query from the URL and updates request.Form,
// decode the from to the given struct with Validator implemented.
func ParseForm(r *http.Request, form Validator) (err error) {
	if err = r.ParseForm(); err == nil {
		if err = Schema.Decode(form, r.Form); err == nil {
			err = form.Validate()
		}
	}
	return
}
