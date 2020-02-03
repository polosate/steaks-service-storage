package main

import (
	"context"
	pb "github.com/polosate/storage-service/proto/storage"
)

type handler struct {
	repository
}

func (s *handler) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

	storage, err := s.repository.FindAvailable(ctx, MarshalSpecification(req))

	if err != nil {
		return err
	}

	res.Storage  = UnmarshalStorage(storage)
	return nil
}

func (s *handler) Create(ctx context.Context, req *pb.Storage, res *pb.Response) error {
	if err := s.repository.Create(ctx, MarshalStorage(req)); err != nil {
		return err
	}

	res.Storage = req

	return nil
}
