all: clean net

clean:
	@echo "Cleaning"
	rm -f net

net:
	@echo "Building net command line tool"
	go build -o net .