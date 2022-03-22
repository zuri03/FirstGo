FROM golang:1.17 as builder
WORKDIR /go/src/github.com/zuri03/FirstGo/
COPY ./ ./
RUN CGO_ENABLED=0 go build -o  build/FirstGo-Server cmd/FirstGo-Server/main.go


FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /go/src/github.com/zuri03/FirstGo/
COPY --from=builder /go/src/github.com/zuri03/FirstGo/build/FirstGo-Server ./cmd/FirstGo-Server/main
COPY --from=builder /go/src/github.com/zuri03/FirstGo/.env ./.env
WORKDIR /go/src/github.com/zuri03/FirstGo/cmd/FirstGo-Server/
EXPOSE 8080
ENTRYPOINT [ "./main" ]
