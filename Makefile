build:
	release/build.sh

clean:
	go clean -modcache
	rm release/bin/ninetails-*
