FirstGo:
	go build cmd/main.go -o bin/FirstGo
	
Server Test:
	go test pkg/server/server.go pkg/server/server_test.go
