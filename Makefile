
build:
	mkdir -p build && go build -o build/argocd-tag-plugin

clean:
	rm -rf build
