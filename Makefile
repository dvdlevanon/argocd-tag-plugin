
build:
	mkdir -p build && go build -buildvcs=false -o build/argocd-tag-plugin

clean:
	rm -rf build
