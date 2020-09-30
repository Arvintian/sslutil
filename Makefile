
build: bin/sslutil

bin/sslutil: main.go gen_cert.go
	rm -rf bin && mkdir bin
	go build -o bin/sslutil main.go gen_cert.go

clean:
	rm -rf bin