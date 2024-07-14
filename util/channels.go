package util

type Channel struct {
	Id               string `json:"_id"`
	Name             string `json:"name"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
	Archived         bool   `json:"archived"`
}