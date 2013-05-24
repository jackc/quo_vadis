package quo_vadis

import (
	"io"
	"net/http"
	"strings"
)

type Endpoint struct {
	handler          http.Handler
	placeholderNames []string
}

type Router struct {
	methodEndpoints map[string]*Endpoint
	staticHandlers  map[string]*Router
	placeholder     *Router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	segments := segmentizePath(req.URL.Path)
	if endpoint, matchedPlaceholders, ok := r.FindEndpoint(req.Method, segments, []string{}); ok {
		addRouteParametersToRequest(endpoint.placeholderNames, matchedPlaceholders, req)
		endpoint.handler.ServeHTTP(w, req)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Not Found")
	}
}

func (r *Router) addRouteFromSegments(method string, segments []string, endpoint *Endpoint) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		var subrouter *Router
		if strings.HasPrefix(head, ":") {
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
		subrouter.addRouteFromSegments(method, tail, endpoint)
	} else {
		r.methodEndpoints[method] = endpoint
	}
}

func (r *Router) AddRoute(method string, path string, handler http.Handler) {
	segments := segmentizePath(path)
	placeholderNames := extractPlaceholderNames(segments)
	endpoint := &Endpoint{handler: handler, placeholderNames: placeholderNames}
	r.addRouteFromSegments(method, segments, endpoint)
}

func (r *Router) FindEndpoint(method string, segments, matchedPlaceholders []string) (*Endpoint, []string, bool) {
	if len(segments) > 0 {
		head, tail := segments[0], segments[1:]
		if subrouter, present := r.staticHandlers[head]; present {
			return subrouter.FindEndpoint(method, tail, matchedPlaceholders)
		} else if r.placeholder != nil {
			matchedPlaceholders = append(matchedPlaceholders, head)
			return r.placeholder.FindEndpoint(method, tail, matchedPlaceholders)
		} else {
			return nil, nil, false
		}
	}
	endpoint, present := r.methodEndpoints[method]
	return endpoint, matchedPlaceholders, present
}

func segmentizePath(path string) (segments []string) {
	for _, s := range strings.Split(path, "/") {
		if len(s) != 0 {
			segments = append(segments, s)
		}
	}
	return
}

func extractPlaceholderNames(segments []string) (placeholderNames []string) {
	for _, s := range segments {
		if strings.HasPrefix(s, ":") {
			placeholderNames = append(placeholderNames, s[1:])
		}
	}
	return
}

func addRouteParametersToRequest(names, values []string, req *http.Request) {
	query := req.URL.Query()
	for i := 0; i < len(names); i++ {
		query.Set(names[i], values[i])
	}
	req.URL.RawQuery = query.Encode()
}

func NewRouter() (r *Router) {
	r = new(Router)
	r.methodEndpoints = make(map[string]*Endpoint)
	r.staticHandlers = make(map[string]*Router)
	return
}
