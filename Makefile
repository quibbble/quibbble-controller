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

proto_quibbble:
	protoc --proto_path=proto --go_out=pkg --go-grpc_out=pkg proto/game.proto 
	protoc --proto_path=proto --go_out=pkg --go-grpc_out=pkg \
		--go_opt=Mgame.proto=github.com/quibbble/quibbble-controller/pkg/game \
		proto/sdk.proto
	protoc --proto_path=proto --go_out=pkg --go-grpc_out=pkg \
		--go_opt=Mgame.proto=github.com/quibbble/quibbble-controller/pkg/game \
		proto/controller.proto

proto_tictactoe:
	protoc --proto_path=examples/tictactoe/proto --go_out=examples/tictactoe examples/tictactoe/proto/tictactoe.proto \
		&& mv examples/tictactoe/tictactoe/* examples/tictactoe/impl/go \
		&& rm -r examples/tictactoe/tictactoe
