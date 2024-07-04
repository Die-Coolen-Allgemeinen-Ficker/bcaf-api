package util

type Smp struct {
	Name              string  `json:"name"`
	Ip                *string `json:"ip"`
	Description       string  `json:"description"`
	Version           string  `json:"version"`
	Modpack           *string `json:"modpack"`
	StartingTimestamp int64   `json:"startingTimestamp"`
}

type SmpWorld struct {
	Name     string  `json:"name"`
	Download string  `json:"download"`
	Version  string  `json:"version"`
	Modpack  *string `json:"modpack"`
}