

build: clean
	go build -o ksheet main.go

clean:
	rm -rf ksheet