ifeq ($(OS),Windows_NT)
    EXT=.exe
else
    EXT=
endif

proto:
	protoc -I rpc rpc/wxfetcher.proto --go_out=plugins=grpc:rpc

wxfetcher:
	go build -o bin/wxfetcher$(EXT) ./cmd/wxfetcher
