.PHONY: build

compilerVersion = protoc-21.12-win64
sourceDir := ./src
artifactsDir := ../pkg/grpc
compilerPath = ./$(compilerVersion)/bin/protoc.exe
includePath := ./$(compilerVersion)/include

build:
	if [ -d $(artifactsDir) ]; then rm -rf $(artifactsDir); fi

	mkdir -p $(artifactsDir)

	chmod +x $(compilerPath)

	@echo "Compiling proto files..."
	@$(compilerPath) \
		--proto_path=$(sourceDir) \
		--proto_path=$(includePath) \
		--go_out=$(artifactsDir) \
		--go_opt=paths=import,module=card-game-server-prototype/pkg/grpc \
		--go-grpc_out=$(artifactsDir) \
		--go-grpc_opt=paths=import,module=card-game-server-prototype/pkg/grpc \
		$(shell find ./src -name '*.proto')

	chmod -x $(compilerPath)

use_osx:
	$(eval compilerVersion=protoc-21.12-osx-universal_binary)
	$(eval compilerPath=./$(compilerVersion)/bin/protoc)

osx: use_osx build

install:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28; \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

fix_signal11_fail: install
