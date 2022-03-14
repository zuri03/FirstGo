FirstGo-Server:
	go build cmd/FirstGo-Server/main.go -o bin/FirstGo-Server

FirstGo-Socket:
	go build cmd/FirstGo-Socket/main.go -o bin/FirstGo-Socket
	
all: FirstGo-Server FirstGo-Socket

Server-Test:
	go test pkg/server/server.go pkg/server/server_test.go pkg/server/Analyzer.go pkg/server/types.go

clean:
	rm -f bin/FirstGo-Server bin/FirtsGo-Socket
