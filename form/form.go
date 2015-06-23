package form

import (
	"net/http"

	. "github.com/gorilla/schema"
)

var schema = NewDecoder()

type Validator interface {
	Validate() error
}

// Parse parsed the raw query from the URL and updates request.Form,
// decode the from to the given struct with Validator implemented.
func parse(r *http.Request, form Validator) (err error) {
	if err = r.ParseForm(); err == nil {
		if err = schema.Decode(form, r.Form); err == nil {
			err = form.Validate()
		}
	}
	return
}
