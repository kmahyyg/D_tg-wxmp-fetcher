ifeq ($(OS),Windows_NT)
    EXT=.exe
else
    EXT=
endif

proto:
	protoc -I proto proto/wxfetcher.proto --go_out=plugins=grpc:proto

wxfetcher:
	go build -o bin/wxfetcher$(EXT) ./cmd/wxfetcher
