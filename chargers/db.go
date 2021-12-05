package chargers

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type database struct {
	db     *mongo.Database
	logger log.Logger
}

func NewDatabase(db *mongo.Database, logger log.Logger) ChargerDB {
	return &database{
		db:     db,
		logger: log.With(logger, "database", "mongoDB"),
	}
}

func (dat *database) CreateCharger(ctx context.Context, charger Charger) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := dat.db.Collection("chargers").InsertOne(ctx, charger)
	if err != nil {
		dat.logger.Log("Error inserting charger into DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetCharger(ctx context.Context, id string) (Charger, error) {
	tempCharger := Charger{}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = dat.db.Collection("chargers").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempCharger)
	if err != nil {
		dat.logger.Log("Error getting charger from DB: ", err.Error())
		return tempCharger, err
	}
	return tempCharger, nil
}
func (dat *database) DeleteCharger(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error deleting charger from DB: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := dat.db.Collection("chargers").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		dat.logger.Log("Error deleting charger from DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetChargers(ctx context.Context) ([]Charger, error) {
	tempCharger := Charger{}
	tempChargers := []Charger{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := dat.db.Collection("chargers").Find(ctx, bson.D{})
	if err != nil {
		dat.logger.Log("Error getting chargers from DB: ", err.Error())
		return tempChargers, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&tempCharger)
		if err != nil {
			dat.logger.Log("Error getting chargers from DB: ", err.Error())
			return tempChargers, err
		}
		tempChargers = append(tempChargers, tempCharger)
	}
	return tempChargers, nil
}
