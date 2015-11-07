.PHONY: run

run: # starts go talks
	go get -u golang.org/x/tools/present
	git submodule update --init
	present -orighost=localhost -http=localhost:3999
