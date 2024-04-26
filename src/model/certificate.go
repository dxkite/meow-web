package model

type Certificate struct {
	Base

	Domain      []string `json:"domain"`
	Description string   `json:"description"`
	Key         string   `json:"key"`
	Certificate string   `json:"certificate"`

	ExpireAt uint64 `json:"expire_at"`
}
