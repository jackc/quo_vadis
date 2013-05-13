package router

import (
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

func TestRouting(t *testing.T) {
	//router := NewRouter()

}
