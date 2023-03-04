all: clean perms client server test

clean:
	rm -rf build

cleanCli:
	rm -rf build/client

cleanServ:
	rm -rf build/server

perms:
	chmod -R 755 scripts

client:
	mkdir -p build/client/scripts
	cp -r testingJsons build/client
	cp scripts/client/* build/client/scripts
	go build -o build/client/client.out src/client/client.go

server:
	mkdir -p build/server/scripts/certificates
	cp scripts/server/* build/server/scripts
	go build -o build/server/server.out src/server/server.go
test:
	docker-compose up
