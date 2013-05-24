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
	handler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) {})
	benchmarkRouter.AddRoute("GET", "/", handler)
	benchmarkRouter.AddRoute("GET", "/foo", handler)
	benchmarkRouter.AddRoute("GET", "/foo/bar", handler)
	benchmarkRouter.AddRoute("GET", "/foo/baz", handler)
	benchmarkRouter.AddRoute("GET", "/foo/bar/baz/quz", handler)
	benchmarkRouter.AddRoute("GET", "/people", handler)
	benchmarkRouter.AddRoute("GET", "/people/search", handler)
	benchmarkRouter.AddRoute("GET", "/people/:id", handler)
	benchmarkRouter.AddRoute("GET", "/users", handler)
	benchmarkRouter.AddRoute("GET", "/users/:id", handler)
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

func TestExtractParameterNames(t *testing.T) {
	test := func(segments []string, expected []string) {
		actual := extractParameterNames(segments)
		if len(actual) != len(expected) {
			t.Errorf("Expected \"%v\" to have %v parameters, but it actually had %v", segments, expected, actual)
			return
		}

		for i := 0; i < len(actual); i++ {
			if actual[i] != expected[i] {
				t.Errorf("Expected \"%v\" to have %v parameters, but it actually had %v", segments, expected, actual)
				return
			}
		}
	}

	test([]string{}, []string{})
	test([]string{"foo"}, []string{})
	test([]string{"foo", ":id"}, []string{"id"})
	test([]string{"foo", ":id", "edit"}, []string{"id"})
}

func TestRouter(t *testing.T) {
	router := NewRouter()

	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "root")
	})
	router.AddRoute("GET", "/", rootHandler)

	widgetIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "widgetIndex")
	})
	router.AddRoute("GET", "/widget", widgetIndexHandler)

	widgetShowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "widgetShow: %v", r.URL.Query().Get("id"))
	})
	router.AddRoute("GET", "/widget/:id", widgetShowHandler)

	widgetEditHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "widgetEdit: %v", r.URL.Query().Get("id"))
	})
	router.AddRoute("GET", "/widget/:id/edit", widgetEditHandler)

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
	get("/widget/1", 200, "widgetShow: 1")
	get("/widget/1/edit", 200, "widgetEdit: 1")

	get("/missing", 404, "404 Not Found")
	get("/widget/1/missing", 404, "404 Not Found")
}

func getBench(b *testing.B, handler http.Handler, path string, expectedCode int) {
	response := httptest.NewRecorder()
	request, err := http.NewRequest("GET", "http://example.com"+path, nil)
	if err != nil {
		b.Fatalf("Unable to create test GET request for %v", path)
	}

	handler.ServeHTTP(response, request)
	if response.Code != expectedCode {
		b.Fatalf("GET %v: expected HTTP code %v, received %v", path, expectedCode, response.Code)
	}
}

func BenchmarkRoutedRequest(b *testing.B) {
	router := getBenchmarkRouter()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		getBench(b, router, "/widgets/important", 200)
	}
}

func BenchmarkFindEndpointRoot(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/"), []string{})
	}
}

func BenchmarkFindEndpointSegment1(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/foo"), []string{})
	}
}

func BenchmarkFindEndpointSegment2(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/people/search"), []string{})
	}
}

func BenchmarkFindEndpointSegment2Placeholder(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/people/1"), []string{})
	}
}

func BenchmarkFindEndpointSegment4(b *testing.B) {
	router := getBenchmarkRouter()

	for i := 0; i < b.N; i++ {
		router.findEndpoint("GET", segmentizePath("/foo/bar/baz/quz"), []string{})
	}
}
