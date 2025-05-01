package handlers

type ErrorResponse struct {
	Error string `json:"error"`
}

type CleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}
