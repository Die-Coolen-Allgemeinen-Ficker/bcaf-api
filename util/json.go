package util

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

func ToJSON[T any](data T) (map[string]interface{}, error) {
	d, err := json.MarshalIndent(&data, "", "    ")
	var result map[string]interface{}
	json.Unmarshal(d, &result)
	return result, err
}

func FromJSON[T any](data map[string]interface{}) (*T, error) {
	d, err := json.Marshal(&data)
	var result T
	json.Unmarshal(d, &result)
	return &result, err
}

func FromJSONRaw[T any](data bson.M) (*T, error) {
	d, err := json.Marshal(&data)
	var result T
	json.Unmarshal(d, &result)
	return &result, err
}