package httpresponse

type ErrorDTO struct {
	Message string `json:"message"`
	Err     string `json:"error"`
}
