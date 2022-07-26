package registrable

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Router struct {
	*http.ServeMux

	// any any other dependencies
}

func (self Router) Register(r HandlerRegistration) {
	if r.HandlerFunc == nil {
		panic("no handler set")
	}

	if r.HandlerFunc != nil {
		self.ServeMux.HandleFunc(r.Path, r.HandlerFunc)
	}
}

func NewRouter() *Router {
	r := &Router{http.NewServeMux()}
	RegisterMethods[HandlerRegistration](r)

	return r
}

type HandlerRegistration struct {
	// Path the endpoint is registered at
	Path string

	// Your http handler func
	HandlerFunc http.HandlerFunc

	// and any other settings you configured inside your Register
}

func (self Router) ServeFoo() Registration {
	return HandlerRegistration{
		Path: "/foo",
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "fighters")
		},
	}
}

func (self Router) ServeBar() Registration {
	return HandlerRegistration{
		Path: "/bar",
		HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "stool")
		},
	}
}

func TestHTTPRegistration(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		path     string
		expected string
	}{
		{path: "/foo", expected: "fighters"},
		{path: "/bar", expected: "stool"},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", test.path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		res := w.Result()
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		if string(data) != test.expected {
			t.Errorf("expected %v got %v", test.expected, string(data))
		}
	}
}
