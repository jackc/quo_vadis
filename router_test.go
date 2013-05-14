package router

import (
	"net/http"
	"testing"
)

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
