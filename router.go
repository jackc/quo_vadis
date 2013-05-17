package quo_vadis

import (
	"io"
	"net/http"
	"strings"
)

type Router struct {
	methodHandlers map[string]func(http.ResponseWriter, *http.Request)
	staticHandlers map[string]*Router
	placeholder    *Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := segmentizePath(req.URL.Path)
	if handler, ok := r.FindHandler(req.Method, segments); ok {
		handler(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) addRouteFromSegments(method string, segments []string, handler func(http.ResponseWriter, *http.Request)) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		var subrouter *Router
		if head == "?" {
			if r.placeholder == nil {
				r.placeholder = NewRouter()
			}
			subrouter = r.placeholder
		} else {
			if _, present := r.staticHandlers[head]; !present {
				r.staticHandlers[head] = NewRouter()
			}
			subrouter = r.staticHandlers[head]
		}
		subrouter.addRouteFromSegments(method, tail, handler)
	} else {
		r.methodHandlers[method] = handler
	}
}

func (r *Router) AddRoute(method string, path string, handler func(http.ResponseWriter, *http.Request)) {
	segments := segmentizePath(path)
	r.addRouteFromSegments(method, segments, handler)
}

func (r *Router) FindHandler(method string, segments []string) (handler func(http.ResponseWriter, *http.Request), present bool) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if subrouter, present := r.staticHandlers[head]; present {
			return subrouter.FindHandler(method, tail)
		} else if r.placeholder != nil {
			return r.placeholder.FindHandler(method, tail)
		} else {
			return nil, false
		}
	}
	handler, present = r.methodHandlers[method]
	return
}

func segmentizePath(path string) (segments []string) {
	for _, s := range strings.Split(path, "/") {
		if len(s) != 0 {
			segments = append(segments, s)
		}
	}
	return
}

func NewRouter() (r *Router) {
	r = new(Router)
	r.methodHandlers = make(map[string]func(http.ResponseWriter, *http.Request))
	r.staticHandlers = make(map[string]*Router)
	return
}
