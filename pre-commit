#!/bin/sh
# NOTE: This script does not play well when adding partial changes, such as with
# git add -p or git commit -p.

test_fmt() {
    hash gofmt 2>&- || { echo >&2 "gofmt not in PATH."; exit 1; }
    IFS='
'
    exitcode=0
    for file in `git diff --cached --name-only --diff-filter=ACM | grep '\.go$'`
    do
        output=`gofmt -w "$file"`
        if test -n "$output"
        then
            # any output is a syntax error
            echo >&2 "$output"
            exitcode=1
        fi
        git add "$file"
    done
    exit $exitcode
}

case "$1" in
    --about )
        echo "Check Go code formatting"
        ;;
    * )
        test_fmt
        ;;
esac
