package mongotrace

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	mgo_bson "gopkg.in/mgo.v2/bson"
)

// StructToJson Convert interface to json
func StructToJson(data interface{}) string {
	switch data.(type) {
	case bson.A, bson.D, bson.M:
		jsonDocument, _ := mgo_bson.MarshalJSON(data)
		return string(jsonDocument)
	default:
		var jsonData []byte
		jsonData, err := json.Marshal(data)
		if err != nil {
			return ""
		}
		return string(jsonData)
	}
}
