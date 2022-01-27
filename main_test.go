package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCommonPrefix(t *testing.T) {
	for i, c := range []struct {
		a    string
		b    string
		want string
	}{
		{"", "", ""},
		{"x", "", ""},
		{"", "x", ""},
		{"x", "x", "x"},
		{"/a/b", "/a/b", "/a/b"},
		{"/a/b/", "/a/b", "/a/b"},
		{"/a/bb/", "/a/bb", "/a/bb"},
		{"/a/b/", "/a/b/", "/a/b/"},
		{"/aa/bb", "/aa/b", "/aa/"},
		{"/aa", "/a", "/"},
		{"/foo/bar/baz", "/foo/bar/baz/quux", "/foo/bar/baz"},
	} {
		got := commonPrefix(c.a, c.b)
		if got != c.want {
			t.Errorf("testcase %d failed: commonPrefix('%s','%s') should be '%s', got '%s'", i, c.a, c.b, c.want, got)
		}
	}
}

func TestMakeTree(t *testing.T) {
	in := `/
/foo
/foo/bar
/foo/bar/baz
/foo/bar/other
/foo/baz
/bar/x/y
`

	want := `[/, name=/
  [foo, name=/foo
    [bar, name=/foo/bar
      [baz, name=/foo/bar/baz]
      [other, name=/foo/bar/other]
    ]
    [baz, name=/foo/baz]
  ]
  [bar, name=/bar
    [x, name=/bar/x
      [y, name=/bar/x/y]
    ]
  ]
]
`
	var buf bytes.Buffer
	makeTree(bytes.NewReader([]byte(in)), &buf)
	got := buf.String()

	if diff := cmp.Diff(want, got); diff != "" {
		fmt.Println(want)
		fmt.Println(got)
		t.Errorf("makeTree() mismatch (-want +got):\n%s", diff)
	}
}

func TestMakeTreeDeepFirst(t *testing.T) {
	in := `/
/foo/bar/baz
/foo/bar/baz/quux
`

	want := `[/, name=/
  [foo, name=/foo
    [bar, name=/foo/bar
      [baz, name=/foo/bar/baz
        [quux, name=/foo/bar/baz/quux]
      ]
    ]
  ]
]
`
	var buf bytes.Buffer
	makeTree(bytes.NewReader([]byte(in)), &buf)
	got := buf.String()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("makeTree() mismatch (-want +got):\n%s", diff)
	}
}
