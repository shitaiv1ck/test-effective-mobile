package statshttp

import "github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"

type StatsDTOResponse struct {
	TotalPrice int `json:"total_price" example:"4200"`
}

func ToDTO(stats domains.Statistics) StatsDTOResponse {
	return StatsDTOResponse{
		TotalPrice: stats.TotalPrice,
	}
}
