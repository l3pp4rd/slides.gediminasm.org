.PHONY: run

run: # starts go talks
	go get -u golang.org/x/tools/present
	go install golang.org/x/tools/cmd/present
	git submodule update --init
	present -orighost=localhost -http=localhost:3999
