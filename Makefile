

build-image:
	docker build -f ./Dockerfile.server -t credscan/server .
	docker build -f ./Dockerfile.worker -t credscan/worker .

run-server:
	go run cmd/server/server.go
run-worker:
	go run cmd/worker/worker.go