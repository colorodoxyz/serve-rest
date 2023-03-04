#all: clean client server test
all: server test

cleanCli:
	rm -rf build/client
cleanServ:
	rm -rf build/server
client: cleanCli
	mkdir -p build/client
	go build -o build/client/client.out src/client/client.go
server: cleanServ
	mkdir -p build/server/scripts/certificates
	cp scripts/server/* build/server/scripts
	go build -o build/server/server.out src/server/server.go
test:
	docker-compose up
