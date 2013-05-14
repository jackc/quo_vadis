package router

import (
	"net/http"
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

func TestRouterFindHandler(t *testing.T) {
	router := NewRouter()

	if _, present := router.FindHandler(segmentizePath("/missing")); present {
		t.Error("Missing route was erroneously found")
	}

	handler := func(http.ResponseWriter, *http.Request) {}
	router.AddRoute("/foo", handler)

	if _, present := router.FindHandler(segmentizePath("/foo")); !present {
		t.Error("Did not find route when route was expected")
	}

	router.AddRoute("/foo/bar/baz", handler)
	if _, present := router.FindHandler(segmentizePath("/foo/bar/baz")); !present {
		t.Error("Did not find route when route was expected")
	}

	if _, present := router.FindHandler(segmentizePath("/foo/missing")); present {
		t.Error("Missing route was erroneously found")
	}
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
