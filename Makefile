build:
	echo "Build proto file"
	protoc -I. --go_out=plugins=micro:. \
		  proto/storage/storage.proto

	echo "Build docker image"
	docker build -t storage-service .

run:
	echo "Run docker image"
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 storage-service
