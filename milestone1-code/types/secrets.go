package types

type CreateSecretRequest struct {
	PlainText string `json:"plain_text"`
}

type CreateSecretResponse struct {
	ID string `json:"id"`
}

type GetSecretResponse struct {
	Data string `json:"data"`
}

type Error struct {
	Message string `json:"error"`
}
