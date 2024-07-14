package util

type Message struct {
	Content          string `json:"content"`
	Author           string `json:"authorId"`
	ChannelId        string `json:"channelId"`
	CreatedTimestamp int64  `json:"createdTimestamp"`
}

type MessageCount struct {
	Id     string                 `json:"_id"`
	Counts map[string]struct {
		Count      int64 `json:"count"`
		Characters int   `json:"characters"`
	}                             `json:"counts"`
}