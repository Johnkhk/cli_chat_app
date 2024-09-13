# Test by package
this will run all tests in the package
`go test ./test/rpc/... -v`

# Run sequentially: Use this for now
`go test -p 1 ./test/... -v`

# Run parallel:
`gotestsum --format=short-verbose ./test/...`