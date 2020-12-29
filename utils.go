package mongotrace

import "encoding/json"

// StructToJson Convert interface to json
func StructToJson(data interface{}) string {
	var jsonData []byte
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonData)
}
