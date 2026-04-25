package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"numbernama-go/utils"
)

type (
	IRepository interface {
		FindOne(ctx context.Context, collection string, filter any, result any, opts ...options.Lister[options.FindOneOptions]) error
		FindMany(ctx context.Context, collection string, filter any, result any, opts ...options.Lister[options.FindOptions]) error
		InsertOne(ctx context.Context, collection string, document any) (*mongo.InsertOneResult, error)
		InsertMany(ctx context.Context, collection string, documents []any) (*mongo.InsertManyResult, error)
		RemoveFields(ctx context.Context, collection string, filter any, removeFields map[string]any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
		ReplaceOne(ctx context.Context, collection string, filter any, update any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error)
		UpdateOne(ctx context.Context, collection string, filter any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
		UpdateMany(ctx context.Context, collection string, filter any, update any) (*mongo.UpdateResult, error)
		DeleteOne(ctx context.Context, collection string, filter any) (*mongo.DeleteResult, error)
		DeleteMany(ctx context.Context, collection string, filter any) (*mongo.DeleteResult, error)
		CountDocuments(ctx context.Context, collection string, filter any) (int64, error)
		Aggregate(ctx context.Context, collection string, pipeline any, result any) error
	}

	Repository struct {
		mongoClient *mongo.Client
		database    string
	}
)

// PersistenceEnabled is the feature-flag switch point.
// Keep false until Mongo dependencies/config are ready.
const PersistenceEnabled = false

// NewRepository creates a new RepositoryStore with the given Mongo client
func NewRepository(mongoClient *mongo.Client) IRepository {
	return &Repository{
		mongoClient: mongoClient,
		database:    "numbernama",
	}
}

// region methods

// FindOne – Find a single document by a filter.
func (r *Repository) FindOne(ctx context.Context, collection string, filter any, result any, opts ...options.Lister[options.FindOneOptions]) error {
	return r.getCollection(collection).FindOne(ctx, filter, opts...).Decode(result)
}

// FindMany – Find multiple documents by a filter.
func (r *Repository) FindMany(ctx context.Context, collection string, filter any, result any, opts ...options.Lister[options.FindOptions]) error {
	cursor, err := r.getCollection(collection).Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			utils.Log(ctx).Warnf("failed to close cursor: %v", err)
		}
	}()
	return cursor.All(ctx, result)
}

// InsertOne – Insert a single document.
func (r *Repository) InsertOne(ctx context.Context, collection string, document any) (*mongo.InsertOneResult, error) {
	return r.getCollection(collection).InsertOne(ctx, document)
}

// InsertMany – Insert multiple documents.
func (r *Repository) InsertMany(ctx context.Context, collection string, documents []any) (*mongo.InsertManyResult, error) {
	return r.getCollection(collection).InsertMany(ctx, documents)
}

func (r *Repository) RemoveFields(ctx context.Context, collection string, filter any, fieldsToRemove map[string]any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	if len(fieldsToRemove) == 0 {
		return nil, nil
	}

	update := bson.D{{Key: "$unset", Value: fieldsToRemove}}
	return r.getCollection(collection).UpdateOne(ctx, filter, update, opts...)
}

// ReplaceOne – Update a single document.
func (r *Repository) ReplaceOne(ctx context.Context, collection string, filter any, entity any, opts ...options.Lister[options.ReplaceOptions]) (*mongo.UpdateResult, error) {
	return r.getCollection(collection).ReplaceOne(ctx, filter, entity, opts...)
}

// UpdateOne – Update a single document.
func (r *Repository) UpdateOne(ctx context.Context, collection string, filter any, entity any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
	update := bson.D{{Key: "$set", Value: entity}}
	return r.getCollection(collection).UpdateOne(ctx, filter, update, opts...)
}

// UpdateOne2 – Update a single document.
func (r *Repository) UpdateOne2(ctx context.Context, collection string, filter any, entity any) (*mongo.UpdateResult, error) {
	updateBSON, err := entityToBSON(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity, err: %s", err.Error())
	}
	update := bson.D{{Key: "$set", Value: updateBSON}}
	return r.getCollection(collection).UpdateOne(ctx, filter, update)
}

// UpdateMany – Update multiple documents.
func (r *Repository) UpdateMany(ctx context.Context, collection string, filter any, entity any) (*mongo.UpdateResult, error) {
	return r.getCollection(collection).UpdateMany(ctx, filter, entity)
}

// DeleteOne – Delete a single document.
func (r *Repository) DeleteOne(ctx context.Context, collection string, filter any) (*mongo.DeleteResult, error) {
	return r.getCollection(collection).DeleteOne(ctx, filter)
}

// DeleteMany – Delete multiple documents.
func (r *Repository) DeleteMany(ctx context.Context, collection string, filter any) (*mongo.DeleteResult, error) {
	return r.getCollection(collection).DeleteMany(ctx, filter)
}

// CountDocuments – Count the number of documents matching a filter.
func (r *Repository) CountDocuments(ctx context.Context, collection string, filter any) (int64, error) {
	return r.getCollection(collection).CountDocuments(ctx, filter)
}

// Aggregate – Perform aggregation.
func (r *Repository) Aggregate(ctx context.Context, collection string, pipeline any, result any) error {
	cursor, err := r.getCollection(collection).Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			utils.Log(ctx).Warnf("failed to close cursor: %v", err)
		}
	}()
	return cursor.All(ctx, result)
}

// region Support function

func (r *Repository) getCollection(name string) *mongo.Collection {
	return r.mongoClient.Database(r.database).Collection(name)
}

func entityToBSON(entity any) (bson.D, error) {
	updateJSON, err := json.Marshal(entity)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal entity: %v", err)
	}

	var updateBSON bson.D
	err = bson.UnmarshalExtJSON(updateJSON, true, &updateBSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal entity into BSON: %v", err)
	}

	return updateBSON, nil
}
