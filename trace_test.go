package mongotrace

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"reflect"
	"testing"
)

type Person struct {
	ID   string `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func TestTraceOperationUpdate(t *testing.T) {
	t.Run("TestTraceOperationUpdate", func(t *testing.T) {
		ctx := context.TODO()
		clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			t.Error(err)
			return
		}
		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			t.Error(err)
			return
		}

		db := client.Database("Test")

		initialRecord := Person{
			ID:   uuid.Must(uuid.NewRandom()).String(),
			Name: "Jane",
		}
		updatedRecord := Person{
			ID:   initialRecord.ID,
			Name: "Jane Doe",
		}
		_, err = db.Collection("Person").InsertOne(ctx, initialRecord)
		if err != nil {
			t.Error(err)
			return
		}
		filter := bson.M{"_id": initialRecord.ID}
		_, err = db.Collection("Person").UpdateOne(ctx, filter, bson.M{"$set": updatedRecord})
		if err != nil {
			t.Error(err)
			return
		}

		got, err := TraceOperationUpdate(db, "Trace", "Person", initialRecord, filter)
		if err != nil {
			t.Error(err)
			return
		}

		differenceWanted := fmt.Sprintf("{\n    \"_id\": \"%s\" (string),\n    \"name\": <span style=\"background-color: #fcff7f\">\"Jane\" (string) => \"Jane Doe\" (string)</span>\n} (object)", initialRecord.ID)
		if !reflect.DeepEqual(got.Difference, differenceWanted) {
			t.Errorf("Test: TraceOperationUpdate() \ngot =\n %v \n\n want \n %v", got.Difference, differenceWanted)
		}
	})
}
