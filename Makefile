all: clean net

clean:
	@echo "Cleaning"
	rm -f net
	rm -f /usr/local/bin/net

net:
	@echo "Building net command line tool"
	go build -o net .

install: clean net
	@echo "Installing net command line tool in /usr/local/bin"
	cp net /usr/local/bin/net