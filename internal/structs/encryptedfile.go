package structs

type EncryptedFile struct {
	Kid          KeyvaultObjectId `json:"kid,omitempty"`
	Data         string           `json:"data,omitempty"`
	LastModified JTime            `json:"lastmodified,omitempty"`
}
