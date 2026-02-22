.PHONY: test demo

test:
	clear; go test -v ./...

demo:
	clear; go test -v -run TestVisual ./...
