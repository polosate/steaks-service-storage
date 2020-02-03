// steaks-service-storage/main.go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/micro/go-micro"
	pb "github.com/polosate/storage-service/proto/storage"
)

type repository interface {
	FindAvailable(*pb.Specification) (*pb.Storage, error)
}

type StorageRepository struct {
	storages []*pb.Storage
}

func (repo *StorageRepository) FindAvailable(spec *pb.Specification) (*pb.Storage, error) {
	for _, storage := range repo.storages {
		if spec.Capacity <= storage.Capacity {
			return storage, nil
		}
	}
	return nil, errors.New("No storage found by that spec.")
}

// Our grpc service handler
type service struct {
	repo repository
}

func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {
	storage, err := s.repo.FindAvailable(req)
	if err != nil {
		return err
	}

	res.Storage = storage
	return nil
}

func main() {
	storages := []*pb.Storage{
		&pb.Storage{Id: 1, Name: "Storage 01", Capacity: 1},
	}
	repo := &StorageRepository{storages}

	srv := micro.NewService(
		micro.Name("steaks.storage.service"),
	)

	srv.Init()

	// Register our implementation with
	pb.RegisterStorageServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}

