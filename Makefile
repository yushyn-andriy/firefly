WEBASSEMBLY_GOOS=js
WEBASSEMBLY_GOARCH=wasm
WEBASSEMBLY_TARGET_PATH=./webassembly/yushyn-andriy.github.io/compiler.wasm
WEBASSEMBLY_SOURCE_PATH=./webassembly/yushyn-andriy.github.io/cmd/wasm

SERVER_PATH=./webassembly/yushyn-andriy.github.io

wasm:
	GOOS=$(WEBASSEMBLY_GOOS) GOARCH=$(WEBASSEMBLY_GOARCH) go build -o $(WEBASSEMBLY_TARGET_PATH)  $(WEBASSEMBLY_SOURCE_PATH)


server: wasm
	cd  $(SERVER_PATH) && go run ./cmd/server/main.go


update_page: wasm
	cd  $(SERVER_PATH) && git add . && git commit -m "update" && git push origin gh-pages


firefly:
	go run main.go

