all: clean build install

clean:
	@echo "Cleaning"
	rm -f net
	sudo rm -f /usr/local/bin/net

build: clean
	@echo "Building net command line tool"
	go build -o net .

install: build
	@echo "Installing net command line tool in /usr/local/bin"
	sudo cp net /usr/local/bin/net
	@echo "Done"
	net -h