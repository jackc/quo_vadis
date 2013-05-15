package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var benchmarkRouter *Router

func getBenchmarkRouter() *Router {
	if benchmarkRouter != nil {
		return benchmarkRouter
	}

	benchmarkRouter := NewRouter()
	handler := func(http.ResponseWriter, *http.Request) {}
	benchmarkRouter.AddRoute("/", handler)
	benchmarkRouter.AddRoute("/foo", handler)
	benchmarkRouter.AddRoute("/foo/bar", handler)
	benchmarkRouter.AddRoute("/foo/baz", handler)
	benchmarkRouter.AddRoute("/people", handler)
	benchmarkRouter.AddRoute("/people/search", handler)
	benchmarkRouter.AddRoute("/people/:id", handler)
	benchmarkRouter.AddRoute("/users", handler)
	benchmarkRouter.AddRoute("/users/:id", handler)
	benchmarkRouter.AddRoute("/widgets", handler)
	benchmarkRouter.AddRoute("/widgets/important", handler)

	return benchmarkRouter
}

func TestSegmentizePath(t *testing.T) {
	test := func(path string, expected []string) {
		actual := segmentizePath(path)
		if len(actual) != len(expected) {
			t.Errorf("Expected \"%v\" to be segmented into %v, but it actually was %v", path, expected, actual)
			return
		}

		for i := 0; i < len(actual); i++ {
			if actual[i] != expected[i] {
				t.Errorf("Expected \"%v\" to be segmented into %v, but it actually was %v", path, expected, actual)
				return
			}
		}
	}

	test("/", []string{})
	test("/foo", []string{"foo"})
	test("/foo/", []string{"foo"})
	test("/foo/bar", []string{"foo", "bar"})
	test("/foo/bar/", []string{"foo", "bar"})
	test("/foo/bar/baz", []string{"foo", "bar", "baz"})
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "root")
	}
	router.AddRoute("/", rootHandler)

	widgetIndexHandler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "widgetIndex")
	}
	router.AddRoute("/widget", widgetIndexHandler)

	get := func(path string, expectedCode int, expectedBody string) {
		response := httptest.NewRecorder()
		request, err := http.NewRequest("GET", "http://example.com"+path, nil)
		if err != nil {
			t.Errorf("Unable to create test GET request for %v", path)
		}

		router.ServeHTTP(response, request)
		if response.Code != expectedCode {
			t.Errorf("GET %v: expected HTTP code %v, received %v", path, expectedCode, response.Code)
		}
		if response.Body.String() != expectedBody {
			t.Errorf("GET %v: expected HTTP response body \"%v\", received \"%v\"", path, expectedBody, response.Body.String())
		}
	}

	get("/", 200, "root")
	get("/widget", 200, "widgetIndex")

	get("/missing", 404, "404 Not Found")
	get("/widget/missing", 404, "404 Not Found")
}

func BenchmarkFindHandlerRoot(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler(segmentizePath("/"))
	}
}

func BenchmarkFindHandlerSingleLevel(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler(segmentizePath("/foo"))
	}
}

func BenchmarkFindHandlerSecondLevel(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler(segmentizePath("/people/search"))
	}
}
