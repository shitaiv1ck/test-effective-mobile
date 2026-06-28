package httpresponse

type ErrorDTO struct {
	Message string `json:"message" example:"short error text"`
	Err     string `json:"error"   example:"full error text"`
}
