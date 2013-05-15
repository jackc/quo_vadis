package router

import (
	"io"
	"net/http"
	"strings"
)

type Router struct {
	handler        func(http.ResponseWriter, *http.Request)
	staticHandlers map[string]*Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := segmentizePath(req.URL.Path)
	if handler, ok := r.FindHandler(segments); ok {
		handler(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) addRouteFromSegments(segments []string, handler func(http.ResponseWriter, *http.Request)) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if _, present := r.staticHandlers[head]; !present {
			r.staticHandlers[head] = NewRouter()
		}
		r.staticHandlers[head].addRouteFromSegments(tail, handler)

	} else {
		r.handler = handler
	}
}

func (r *Router) AddRoute(path string, handler func(http.ResponseWriter, *http.Request)) {
	segments := segmentizePath(path)
	r.addRouteFromSegments(segments, handler)
}

func (r *Router) FindHandler(segments []string) (handler func(http.ResponseWriter, *http.Request), present bool) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if subrouter, present := r.staticHandlers[head]; present {
			return subrouter.FindHandler(tail)
		}
	} else {
		if r.handler != nil {
			return r.handler, true
		}
	}
	return nil, false
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
	r.staticHandlers = make(map[string]*Router)
	return
}
