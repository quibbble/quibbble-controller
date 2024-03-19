build_qs:
	docker build -t quibbble/server:latest -f build/server.Dockerfile .

push_qs:
	docker push quibbble/server:latest

run_qs:
	docker run -p 8080:8080 quibbble/server:latest

build_qc:
	docker build -t quibbble/controller:latest -f build/controller.Dockerfile .

push_qc:
	docker push quibbble/controller:latest

run_qc:
	docker run -p 8080:8080 quibbble/controller:latest
