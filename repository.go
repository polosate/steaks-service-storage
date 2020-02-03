// storage-service/repository.go
package main


import (
	"context"
	pb "github.com/polosate/storage-service/proto/storage"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

type Storage struct {
	Id                   int32    `json:"id"`
	Capacity             int32    `json:"capacity"`
	Name                 string   `json:"name"`
	Available            bool     `json:"available"`
	OwnerId              string   `json:"owner_id"`
}

type Specification struct {
	Capacity int32 `json:"capacity"`
}

func MarshalSpecification(specification *pb.Specification) *Specification {
	return &Specification{
		Capacity: specification.Capacity,
	}
}


func MarshalStorage(storage *pb.Storage) *Storage {
	return &Storage{
		Id: storage.Id,
		Name: storage.Name,
		Available: storage.Available,
		Capacity: storage.Capacity,
		OwnerId: storage.OwnerId,
	}
}

func UnmarshalStorage(storage *Storage) *pb.Storage {
	return &pb.Storage{
		Id: storage.Id,
		Name: storage.Name,
		Available: storage.Available,
		Capacity: storage.Capacity,
		OwnerId: storage.OwnerId,
	}
}

type repository interface {
	FindAvailable(context context.Context, specification *Specification) (*Storage, error)
	Create(context context.Context, storage *Storage) error

}

type StorageRepository struct {
	collection *mongo.Collection
}

func (repository *StorageRepository) Create(ctx context.Context, storage *Storage) error {
	_, err := repository.collection.InsertOne(ctx, storage)
	return err
}


func (repository *StorageRepository) FindAvailable(ctx context.Context, spec *Specification) (*Storage, error) {
	filter := bson.D{{
		"capacity",
		bson.D{{
			"$lte",
			spec.Capacity,
		},
		},
	}}
	storage := &Storage{}
	if err := repository.collection.FindOne(ctx, filter).Decode(storage); err != nil {
		return nil, err
	}
	return storage, nil
}
