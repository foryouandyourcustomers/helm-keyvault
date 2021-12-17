package structs

type Key struct {
	Kid     KeyvaultObjectId
	Name    string
	Version string `json:",omitempty"`
}

type KeyList struct {
	Keys []Key
}
