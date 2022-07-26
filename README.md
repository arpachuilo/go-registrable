# go-registrable

[![Go Reference](https://pkg.go.dev/badge/github.com/arpachuilo/go-registrable.svg)](https://pkg.go.dev/github.com/arpachuilo/go-registrable)

Just another way to help compartmentalize your go code.

## Why

Started playing with this idea of organizing my code because I got tired of the distance between my http handler functions and their addition to the multiplexer (along with various other settings such as authorization). Just wanted everything to be bundled to make it easy to spin up new routes. Can be made to work with whatever framework you may already be using could even be used for things other than setting up http handlers.

If what you are registering requires a certain ordering, checkout `OrderedRegistrable`.

## How

The code for this is fairly short. All it does is use reflection to look for any methods will return a `Registration` on you struct. It runs those methods, then pipes it to your `Register` function.

## Example

```go

// our 'Registrable'
type Router struct {
	*http.ServeMux

	// any other dependencies such as a database connection
}

// does all the boilerplate work
func (self Router) Register(r HandlerRegistration) {
	if r.HandlerFunc == nil {
		panic("no handler set")
	}

	if r.HandlerFunc != nil {
		self.ServeMux.HandleFunc(r.Path, r.HandlerFunc)
	}
}

func NewRouter() *Router {
  // instantiate
	r := &Router{http.NewServeMux()}

  // register methods
	RegisterMethods[HandlerRegistration](r)

	return r
}


// our 'Registration'
type HandlerRegistration struct {
	// path the endpoint is registered at
	Path string

	// your http handler func
	HandlerFunc http.HandlerFunc

	// any other settings you need setup inside your Register
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
```
