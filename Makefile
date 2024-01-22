
build:
	mkdir -p build && CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o build/argocd-tag-plugin

clean:
	rm -rf build
