package keys

type Key struct {
	Kid     string
	Name    string
	Version string `json:",omitempty"`
}

type KeyList struct {
	Keys []Key
}
