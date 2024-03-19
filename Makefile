build_server:
	docker build -t quibbble/server:latest -f build/server.Dockerfile .

push_server:
	docker push quibbble/server:latest

run_server:
	docker run -p 8080:8080 quibbble/server:latest

build_controller:
	docker build -t quibbble/controller:latest -f build/controller.Dockerfile .

push_controller:
	docker push quibbble/controller:latest

run_controller:
	docker run -p 8080:8080 quibbble/controller:latest
