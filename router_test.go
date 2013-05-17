package quo_vadis

import (
	"fmt"
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
	benchmarkRouter.AddRoute("GET", "/", handler)
	benchmarkRouter.AddRoute("GET", "/foo", handler)
	benchmarkRouter.AddRoute("GET", "/foo/bar", handler)
	benchmarkRouter.AddRoute("GET", "/foo/baz", handler)
	benchmarkRouter.AddRoute("GET", "/foo/bar/baz/quz", handler)
	benchmarkRouter.AddRoute("GET", "/people", handler)
	benchmarkRouter.AddRoute("GET", "/people/search", handler)
	benchmarkRouter.AddRoute("GET", "/people/?", handler)
	benchmarkRouter.AddRoute("GET", "/users", handler)
	benchmarkRouter.AddRoute("GET", "/users/?", handler)
	benchmarkRouter.AddRoute("GET", "/widgets", handler)
	benchmarkRouter.AddRoute("GET", "/widgets/important", handler)

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
	router.AddRoute("GET", "/", rootHandler)

	widgetIndexHandler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "widgetIndex")
	}
	router.AddRoute("GET", "/widget", widgetIndexHandler)

	widgetShowHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "widgetShow")
	}
	router.AddRoute("GET", "/widget/?", widgetShowHandler)

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
	get("/widget/1", 200, "widgetShow")

	get("/missing", 404, "404 Not Found")
	get("/widget/1/missing", 404, "404 Not Found")
}

func BenchmarkFindHandlerRoot(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler("GET", segmentizePath("/"))
	}
}

func BenchmarkFindHandlerSegment1(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler("GET", segmentizePath("/foo"))
	}
}

func BenchmarkFindHandlerSegment2(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler("GET", segmentizePath("/people/search"))
	}
}

func BenchmarkFindHandlerSegment2Placeholder(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler("GET", segmentizePath("/people/1"))
	}
}

func BenchmarkFindHandlerSegment4(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.FindHandler("GET", segmentizePath("/foo/bar/baz/quz"))
	}
}
