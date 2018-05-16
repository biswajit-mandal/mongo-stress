package main

import (
	"context"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"log"
)

type MongoDBHandler struct {
	mgoDB   *mgo.Database
}

func createMongoClientDB() *mgo.Database {
	var (
		mongo = "mongodb://" + MongoIP + ":" + MongoPort
	)
	log.Println("Connecting Mongo", mongo)
	client, err := mgo.Connect(context.Background(), mongo, nil)
	if err != nil {
		log.Println("Getting mgoSession as:", err)
		createMongoClientDB()
	}
	return client.Database(FlowsDBStr)
}

func createRowInDB(mg *MongoDBHandler, collection string, data SchemaData) error {
	db := mg.mgoDB
	c := db.Collection(collection)
	if _, err := c.InsertOne(context.Background(), data); err != nil {
		stats.InsertErrorCount++
		return err
	}
	stats.InsertSuccessCount++
	return nil
}

func getCollectionCout(mg *MongoDBHandler, collection string) (int64, error) {
	var (
		result int64
		err    error
	)
	db := mg.mgoDB
	c := db.Collection(collection)
	if result, err = c.Count(context.Background(), nil); err != nil {
		return 0, err
	}
	return result, nil
}

func DBSetup(mg *MongoDBHandler) {
	mg.mgoDB = createMongoClientDB()
	return
}
