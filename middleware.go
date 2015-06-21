package rex

import "net/http"

type middleware struct {
	cache http.Handler
	stack []func(http.Handler) http.Handler
}

// Implements the net/http Handler interface and calls the middleware stack.
func (self *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if self.cache == nil {
		// setup the whole middleware modules in a FIFO chain.
		var next http.Handler = http.DefaultServeMux
		for index := len(self.stack) - 1; index >= 0; index-- {
			next = self.stack[index](next)
		}
		self.cache = next
	}
	self.cache.ServeHTTP(w, r)
}
