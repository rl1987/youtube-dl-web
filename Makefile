build: deps
	go build .

deps:
	go get -u ./...

clean:
	go clean 

