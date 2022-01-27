#!/bin/bash
cwd=$(pwd)
path=$GOPATH/src/github.com/prometheus/prometheus/tsdb


function latextree() {
	# path from find     | after dirname  | what we want
	# ./wal.go           | .              | /
	# ./foo/bar.go       | ./foo          | /foo
	{
		cd $path
		find . -name '*.go' | xargs dirname | sed -e 's#\.$#/#' -e 's#^\.##' | LC_ALL=C sort | uniq 
	} | go run main.go
}

function deps() {
	path='github.com/prometheus/prometheus/tsdb'
	cd $GOPATH/src/$path
	# for all go packages, retrieve package name and list of imports
	find . -name '*.go' | xargs grep "$path/" |
		sed 's# //.*##'                   | # remove code comments
		sed 's/\/\?[^/]\+\.go://'         | # ./file/like/so.go: -> ./file/like 
		sed -e 's#^\.#/#' -e 's#//#/#g'   | # normalize ./ to / and . to /, this also canonicalizes double slashes in import paths
		tr -d '"'                         | # imports have "" around them
		awk '{print $1,$NF}'              | # select first and last field, i.o.w. ignore the alias if there is any
		sed "s# $path# #"                 | # in 2nd field, represent <path> as '.' to correspond to first field.
		LC_ALL=C sort | uniq                # same import may appear in multiple files in same package
}

function latexarrows() {
	deps | while read from to; do
		color=blue
		[[ $from == '/' ]] && color=red
		[[ $from == '/errors' ]] && color=red
		[[ $from == '/tombstones' ]] && color=gray
		[[ $from == '/tsdbutil' ]] && color=yellow
		[[ $from == '/index' ]] && color=purple
		[[ $from == '/chunks' ]] && color=brown
		[[ $from == '/agent' ]] && color=orange
		[[ $from == '/wal' ]] && color=green
		printf '\draw [thin, %s, -{Triangle[]}] (%s.west) [bend left] to node [pos=.25, left, inner sep=1pt, xshift=-5pt] {} (%s) ;\n' $color $from $to
	done
}

cat <<EOF > plot.tex

\documentclass[tikz,border=5pt]{standalone}
\usepackage{forest}
\usetikzlibrary{arrows.meta}
\begin{document}
  \begin{forest}
    for tree={
      parent anchor=south,
      child anchor=north,
      tier/.wrap pgfmath arg={tier#1}{level()},
      font=\sffamily
    }
    `latextree`
    `latexarrows`
  \end{forest}
\end{document}
EOF

IMAGE=blang/latex:ubuntu
exec docker run --rm -i --user="$(id -u):$(id -g)" --net=none -v "$PWD":/data "$IMAGE" pdflatex plot.tex
