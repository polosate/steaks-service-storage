// storage-service/main.go
package main

import (
	"context"
	"os"

	"fmt"
	"github.com/micro/go-micro"
	pb "github.com/polosate/storage-service/proto/storage"
	"log"
)

const (
	defaultHost = "datastore:27017"
)

func createDummyData(repo repository) {
	storages := []*Storage{
		{Id: 0, Name: "Storage_00", Available: true, Capacity: 0, OwnerId: "qwe"},
		{Id: 1, Name: "Storage_01", Available: true, Capacity: 1, OwnerId: "qwe"},
	}
	for _, v := range storages {
		repo.Create(context.Background(), v)
	}
}

func main() {

	srv := micro.NewService(
		micro.Name("steaks.storage.service"),
	)

	srv.Init()

	uri := os.Getenv("DB_HOST")
	if uri == "" {
		uri = defaultHost
	}

	client, err := CreateClient(context.Background(), uri, 0)
	if err != nil {
		log.Panic(err)
	}
	defer client.Disconnect(context.Background())

	storagesCollection := client.Database("steaks").Collection("storages")

	repository := &StorageRepository{storagesCollection}

	createDummyData(repository)

	h := &handler{repository}

	// Register handlers
	pb.RegisterStorageServiceHandler(srv.Server(), h)

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

