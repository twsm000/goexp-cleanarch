package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Method int

const (
	Connect Method = iota
	Delete
	Get
	Head
	Options
	Patch
	Post
	Put
	Trace
)

func NewHandler(path string, method Method, handlerFunc http.HandlerFunc) Handler {
	return Handler{
		Path:        path,
		Method:      method,
		HandlerFunc: handlerFunc,
	}
}

type Handler struct {
	Path   string
	Method Method
	http.HandlerFunc
}

type WebServer struct {
	Router        chi.Router
	Handlers      []Handler
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      []Handler{},
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(handler Handler) {
	s.Handlers = append(s.Handlers, handler)
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	var handle func(string, http.HandlerFunc)
	for _, handler := range s.Handlers {
		switch handler.Method {
		case Connect:
			handle = s.Router.Connect
		case Delete:
			handle = s.Router.Delete
		case Get:
			handle = s.Router.Get
		case Head:
			handle = s.Router.Head
		case Options:
			handle = s.Router.Options
		case Patch:
			handle = s.Router.Patch
		case Post:
			handle = s.Router.Post
		case Put:
			handle = s.Router.Put
		case Trace:
			handle = s.Router.Trace
		default:
			log.Println("ERROR - invalid route:", handler.Path, handler.Method)
			continue
		}
		handle(handler.Path, handler.HandlerFunc)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
