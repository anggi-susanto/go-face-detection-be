package mongo

import (
	"context"

	"github.com/anggi-susanto/go-face-detection-be/config"
	"github.com/anggi-susanto/go-face-detection-be/domain"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PhotoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewPhotoRepository creates a new instance of the PhotoRepository struct
// with the provided MongoDB client and configuration.
//
// Parameters:
// - client: A pointer to a mongo.Client object representing the MongoDB client.
// - config: A pointer to a config.MongoConfig object representing the MongoDB configuration.
//
// Returns:
// - A pointer to a PhotoRepository object representing the newly created repository.
func NewPhotoRepository(client *mongo.Client, config *config.MongoConfig) *PhotoRepository {
	return &PhotoRepository{
		client:     client,
		collection: client.Database(config.Database).Collection(config.Collection),
	}
}

// Create inserts a new photo document into the MongoDB collection.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - photo: A pointer to a domain.Photo object representing the photo to be inserted.
//
// Returns:
// - error: An error object if there was an error inserting the photo, otherwise nil.
func (p *PhotoRepository) Create(ctx context.Context, photo *domain.Photo) error {
	_, err := p.collection.InsertOne(ctx, photo)
	if err != nil {
		// Log the error and return it
		logrus.Error(err)
		return err
	}
	return nil
}

// FindByID finds a photo document in the MongoDB collection by its ID.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - id: The ID of the photo to be found.
//
// Returns:
// - photo: A pointer to a domain.Photo object representing the found photo, or nil if not found.
// - error: An error object if there was an error finding the photo, otherwise nil.
func (p *PhotoRepository) FindByID(ctx context.Context, id string) (*domain.Photo, error) {
	var photo domain.Photo
	err := p.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&photo)
	if err != nil {
		// Log the error and return it
		logrus.Error(err)
		return nil, err
	}
	return &photo, nil
}

// FindAll finds all photo documents in the MongoDB collection.
//
// Parameters:
// - ctx: The context.Context object for the function.
//
// Returns:
// - photos: A slice of domain.Photo objects representing all the found photos.
// - error: An error object if there was an error finding the photos, otherwise nil.
func (p *PhotoRepository) FindAll(ctx context.Context) ([]domain.Photo, error) {
	var photos []domain.Photo
	cursor, err := p.collection.Find(ctx, bson.M{})
	if err != nil {
		// Log the error and return it
		logrus.Error(err)
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var photo domain.Photo
		err := cursor.Decode(&photo)
		if err != nil {
			// Log the error and return it
			logrus.Error(err)
			return nil, err
		}
		photos = append(photos, photo)
	}
	return photos, nil
}

// Delete deletes a photo document from the MongoDB collection by its ID.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - id: The ID of the photo to be deleted.
//
// Returns:
// - error: An error object if there was an error deleting the photo, otherwise nil.
func (p *PhotoRepository) Delete(ctx context.Context, id string) error {
	_, err := p.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		// Log the error and return it
		logrus.Error(err)
		return err
	}
	return nil
}

// Update updates a photo document in the MongoDB collection.
//
// Parameters:
// - ctx: The context.Context object for the function.
// - photo: A pointer to a domain.Photo object representing the photo to be updated.
//
// Returns:
// - error: An error object if there was an error updating the photo, otherwise nil.
func (p *PhotoRepository) Update(ctx context.Context, photo *domain.Photo) error {
	_, err := p.collection.UpdateOne(ctx, bson.M{"_id": photo.ID}, bson.M{"$set": photo})
	if err != nil {
		// Log the error and return it
		logrus.Error(err)
		return err
	}
	return nil
}
