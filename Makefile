build:
	release/release.sh

clean:
	go clean -modcache
	rm release/bin/ninetails-*