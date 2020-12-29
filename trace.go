package mongotrace

import (
	"context"
	"github.com/google/uuid"
	"github.com/nsf/jsondiff"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	OperationADD    = "ADD"
	OperationUPDATE = "UPDATE"
	OperationDELETE = "DELETE"
)

// Trace is a representation of a action perform in database.
type Trace struct {
	// Default infos
	ID        string    `json:"_id" bson:"_id"`
	CreatedAt time.Time `json:"CreatedAt" bson:"CreatedAt"`

	// Collection is the collection where action is perform
	Collection string `json:"Collection" bson:"Collection"`

	/*Operation is the operation who want to perform on db
	OperationADD for InsertOne
	OperationUPDATE for UpdateOne
	OperationDELETE for DeleteOne*/
	Operation string `json:"Operation" bson:"Operation"`

	/*DocumentBefore is the json representation of document before the operation is performs on database
	Need when Operation is UPDATE or DELETE*/
	DocumentBefore string `json:"DocumentBefore" bson:"DocumentBefore"`

	/*DocumentAfter is the json representation of document after the operation is performs on database
	Need when Operation is ADD or UPDATE*/
	DocumentAfter string `json:"DocumentAfter" bson:"DocumentAfter"`

	/*Difference is the html representation of difference between DocumentBefore and DocumentAfter.
	For a good readable render, you must put this html inside a <pre> tags*/
	Difference string `json:"Difference" bson:"Difference"`
}

//save Save the trace in the given traceCollection
func (trace *Trace) save(database *mongo.Database, traceCollection string) error {
	trace.ID = uuid.Must(uuid.NewRandom()).String()
	trace.CreatedAt = time.Now().UTC()
	_, err := database.Collection(traceCollection).InsertOne(context.TODO(), trace)
	return err
}

// TraceOperationAdd set the DocumentAfter and save the trace in the given traceCollection
func TraceOperationAdd(database *mongo.Database, traceCollection string, collection string, documentAfter interface{}) (*Trace, error) {
	trace := Trace{
		Collection:    collection,
		Operation:     OperationADD,
		DocumentAfter: StructToJson(documentAfter),
	}
	err := trace.save(database, traceCollection)
	if err != nil {
		return nil, err
	}

	return &trace, nil
}

// TraceOperationDelete set the DocumentBefore and save the trace in the given traceCollection
func TraceOperationDelete(database *mongo.Database, traceCollection string, collection string, documentBefore interface{}) (*Trace, error) {
	trace := Trace{
		Collection:     collection,
		Operation:      OperationDELETE,
		DocumentBefore: StructToJson(documentBefore),
	}
	err := trace.save(database, traceCollection)
	if err != nil {
		return nil, err
	}

	return &trace, nil
}

// TraceOperationUpdateWithFilter set the DocumentBefore and DocumentAfter(searching by the filter params) get the difference between the two documents and save the trace in the given traceCollection
func TraceOperationUpdateWithFilter(database *mongo.Database, traceCollection string, collection string, documentBefore interface{}, filter interface{}) (*Trace, error) {
	trace := Trace{
		Collection:     collection,
		Operation:      OperationUPDATE,
		DocumentBefore: StructToJson(documentBefore),
	}
	err := trace.setDocumentAfterBySearch(database, filter)

	if err != nil {
		return nil, err
	}

	trace.getDifference()

	err = trace.save(database, traceCollection)
	if err != nil {
		return nil, err
	}

	return &trace, nil
}

// TraceOperationUpdate set the DocumentBefore and DocumentAfter get the difference between the two documents and save the trace in the given traceCollection
func TraceOperationUpdate(database *mongo.Database, traceCollection string, collection string, documentBefore interface{}, documentAfter interface{}) (*Trace, error) {
	trace := Trace{
		Collection:     collection,
		Operation:      OperationUPDATE,
		DocumentBefore: StructToJson(documentBefore),
		DocumentAfter:  StructToJson(documentAfter),
	}

	trace.getDifference()

	err := trace.save(database, traceCollection)
	if err != nil {
		return nil, err
	}

	return &trace, nil
}

// setDocumentAfterBySearch search a document in the Collection parsing it to json and set it as DocumentAfter
func (trace *Trace) setDocumentAfterBySearch(database *mongo.Database, filter interface{}) error {
	var document bson.M
	err := database.Collection(trace.Collection).FindOne(context.TODO(), filter).Decode(&document)
	if err != nil {
		return err
	}

	trace.DocumentAfter = StructToJson(document)
	return nil
}

// getDifference sort the html representation of difference between DocumentBefore an DocumentAfter
func (trace *Trace) getDifference() {
	opts := jsondiff.DefaultHTMLOptions()
	opts.PrintTypes = false
	_, desc := jsondiff.Compare([]byte(trace.DocumentBefore), []byte(trace.DocumentAfter), &opts)
	trace.Difference = desc
}
