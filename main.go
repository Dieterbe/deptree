package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	makeTree(os.Stdin, os.Stdout)
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}

// commonPrefix returns the common path prefix
// see test cases
func commonPrefix(a, b string) string {

	if a == b {
		return a
	}

	if len(a) == 0 || len(b) == 0 {
		return ""
	}

	// if one of them is like /abc/de and the other one is /abc/de/ or /abc/de/foo, then we want to return /abc/de
	if strings.HasPrefix(a, b+"/") {
		return b
	}
	if strings.HasPrefix(b, a+"/") {
		return a
	}

	// maximum possible length
	maxl := len(a)
	if len(b) < maxl {
		maxl = len(b)
	}

	var lastSlash int
	var i int
	for i = 0; i < maxl; i++ {
		if a[i] != b[i] {
			return a[:lastSlash]
		}
		if a[i] == '/' {
			lastSlash = i
		}
	}

	return a[:lastSlash+1]
}

func makeTree(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)
	var indent int
	var i int
	var prevAbs string // absolute path

	// note that none of the printed lines ends on a newline
	// for some of the lines we need the subsequent line to determine whether the previous line was finished
	// so for consistency, every line always adds a newline to the previous line, when appropriate.
	for scanner.Scan() {
		abs := scanner.Text()
		if len(abs) != 0 {
			// before we print our node, do any tree level adjustments that should come after previous entry, if any.
			prefix := commonPrefix(abs, prevAbs)

			// close down however many levels of the tree we're leaving, if any
			closing := strings.TrimPrefix(prevAbs, prefix)
			if len(closing) > 0 {
				// close down the previous node.
				_, err := fmt.Fprint(out, "]")
				perror(err)
				indent -= 2
			}
			for cnt := 0; cnt < strings.Count(closing, "/")-1; cnt++ {
				_, err := fmt.Fprintf(out, "\n%s]", strings.Repeat(" ", indent))
				perror(err)
				indent -= 2
			}

			// open however many levels of the tree we are entering
			opening := strings.TrimPrefix(abs, prefix)
			if len(opening) > 1 && opening[0] == '/' {
				opening = opening[1:]
			}
			nodes := strings.Split(opening, "/")
			thisAbs := prefix
			if abs == "/" {
				abs = ""
				thisAbs = ""
				nodes = []string{""}
			}

			for _, node := range nodes {
				thisAbs += "/" + node

				// for the very first line we start at indent 0. for all others, we bump the indent and add a newline to whatever the last line was.
				if i > 0 {
					_, err := fmt.Fprintln(out, "")
					perror(err)
					indent += 2
				}
				if node == "" {
					// we assume this only hapens for the root node. could this happen if path is like foo//bar ?
					// we want to print it as '/'
					node = "/"
				}
				_, err := fmt.Fprintf(out, "%s[%s, name=%s", strings.Repeat(" ", indent), node, thisAbs)
				perror(err)
			}

			prevAbs = abs
			i++
		}
	}
	// close all levels that are still open
	_, err := fmt.Fprintln(out, "]")
	perror(err)
	for indent >= 2 {
		indent -= 2
		_, err := fmt.Fprintf(out, "%s]\n", strings.Repeat(" ", indent))
		perror(err)
	}
	return nil
}
