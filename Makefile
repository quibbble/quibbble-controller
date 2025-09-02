create_cluster:
	k3d cluster create quibbble \
		--k3s-arg "--disable=traefik@server:*" \
		--port "80:80@loadbalancer" \
		--agents 1

delete_cluster:
	k3d cluster delete quibbble

docker_build:
	docker build -t quibbble/server:latest -f build/server.Dockerfile .
	docker build -t quibbble/controller:latest -f build/controller.Dockerfile .
	docker build -t quibbble/watcher:latest -f build/watcher.Dockerfile .

docker_import:
	k3d image import quibbble/server:latest --cluster quibbble
	k3d image import quibbble/controller:latest --cluster quibbble
	k3d image import quibbble/watcher:latest --cluster quibbble

docker_run:
	docker run -p 8080:8080 quibbble/server:latest
	docker run -p 8080:8080 quibbble/controller:latest
	docker run quibbble/watcher:latest

test:
	go test ./...

clean: 
	go clean -testcache

proto_go:
	protoc --proto_path=. --go_out=sdks/go --go-grpc_out=sdks/go proto/sdk.proto \
		&& mv sdks/go/sdk/* sdks/go \
		&& rm -r sdks/go/sdk

proto_python:
	python3 -m grpc_tools.protoc -I=proto --python_out=sdks/python --grpc_python_out=sdks/python proto/sdk.proto

proto_js:
	protoc --proto_path=. -I=proto --js_out=sdks/js proto/sdk.proto \
		&& protoc --proto_path=. -I=proto --grpc-web_out=mode=grpcwebtext:sdks/js proto/sdk.proto \
		&& mv sdks/js/proto/* sdks/js && rm -r sdks/js/proto

proto_tictactoe_go:
	protoc --proto_path=. --go_out=examples/tictactoe examples/tictactoe/proto/tictactoe.proto \
		&& mv examples/tictactoe/tictactoe/* examples/tictactoe/go \
		&& rm -r examples/tictactoe/tictactoe
