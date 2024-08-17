# start watch mode
find . -name '*.go' | entr -c gotestsum --format=testdox
