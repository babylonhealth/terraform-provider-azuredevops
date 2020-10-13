package model

type CheckPayloadType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CheckResource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
