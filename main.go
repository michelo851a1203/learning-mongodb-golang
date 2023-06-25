package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	mongoClient, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI("mongodb://michael:secret@localhost:27017/"),
	)

	defer func() {
		cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Fatalf("mongodb disconnect error : %v", err)
		}
	}()

	if err != nil {
		log.Fatalf("connection error :%v", err)
		return
	}

	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("ping mongodb error :%v", err)
		return
	}
	fmt.Println("ping success")

	// database and collection
	database := mongoClient.Database("demo")
	sampleCollection := database.Collection("sampleCollection")
	sampleCollection.Drop(ctx)

	// insert one

	insertedDocument := bson.M{
		"name":       "michael",
		"content":    "test content",
		"bank_money": 1000,
		"create_at":  time.Now(),
	}
	insertedResult, err := sampleCollection.InsertOne(ctx, insertedDocument)

	if err != nil {
		log.Fatalf("inserted error : %v", err)
		return
	}
	fmt.Println("======= inserted id ================")
	log.Printf("inserted ID is : %v", insertedResult.InsertedID)

	// query all data
	fmt.Println("== query all data ==")
	cursor, err := sampleCollection.Find(ctx, options.Find())
	if err != nil {
		log.Fatalf("find collection err : %v", err)
		return
	}
	var queryResult []bson.M
	if err := cursor.All(ctx, &queryResult); err != nil {
		log.Fatalf("query mongodb result")
		return
	}

	for _, doc := range queryResult {
		fmt.Println(doc)
	}

	// insert many data
	fmt.Println("=========== inserted many data ===============")
	insertedManyDocument := []interface{}{
		bson.M{
			"name":       "Andy",
			"content":    "new test content",
			"bank_money": 1500,
			"create_at":  time.Now().Add(36 * time.Hour),
		},
		bson.M{
			"name":       "Jack",
			"content":    "jack content",
			"bank_money": 800,
			"create_at":  time.Now().Add(12 * time.Hour),
		},
	}

	insertedManyResult, err := sampleCollection.InsertMany(ctx, insertedManyDocument)
	if err != nil {
		log.Fatalf("inserted many error : %v", err)
		return
	}

	for _, doc := range insertedManyResult.InsertedIDs {
		fmt.Println(doc)
	}

	fmt.Println("=========== query specific data =====================")
	// query specific data
	filter := bson.D{
		bson.E{
			Key: "bank_money",
			Value: bson.D{
				bson.E{
					Key:   "$gt",
					Value: 900,
				},
			},
		},
	}

	filterCursor, err := sampleCollection.Find(
		ctx,
		filter,
	)
	if err != nil {
		log.Fatalf("filter query data error : %v", err)
		return
	}
	var filterResult []bson.M
	err = filterCursor.All(ctx, &filterResult)
	if err != nil {
		log.Fatalf("filter result %v", err)
		return
	}

	for _, filterDoc := range filterResult {
		fmt.Println(filterDoc)
	}

	updateManyFilter := bson.D{
		bson.E{
			Key:   "name",
			Value: "michael",
		},
	}

	updateSet := bson.D{
		bson.E{
			Key: "$set",
			Value: bson.D{
				bson.E{
					Key:   "bank_money",
					Value: 2000,
				},
			},
		},
	}
	// update
	updateManyResult, err := sampleCollection.UpdateMany(
		ctx,
		updateManyFilter,
		updateSet,
	)
	if err != nil {
		log.Fatalf("update error : %v", err)
		return
	}

	fmt.Println("========= updated modified count ===========")
	fmt.Println(updateManyResult.ModifiedCount)

	// check if updated with find solution
	checkedCursor, err := sampleCollection.Find(
		ctx,
		bson.D{
			bson.E{
				Key:   "name",
				Value: "michael",
			},
		},
	)
	if err != nil {
		log.Fatalf("check result error : %v", err)
		return
	}
	var checkedResult []bson.M
	err = checkedCursor.All(ctx, &checkedResult)
	if err != nil {
		log.Fatalf("get check information error : %v", err)
		return
	}
	fmt.Println("=========== checked updated result ==============")
	for _, checkedDoc := range checkedResult {
		fmt.Println(checkedDoc)
	}
	fmt.Println("===============================")
	// delete Many

	deleteManyResult, err := sampleCollection.DeleteMany(
		ctx,
		bson.D{
			bson.E{
				Key: "bank_money",
				Value: bson.D{
					bson.E{
						Key:   "$lt",
						Value: 1000,
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatalf("delete many data error : %v", err)
		return
	}
	fmt.Println("===== delete many data modified count =====")
	fmt.Println(deleteManyResult.DeletedCount)
}
