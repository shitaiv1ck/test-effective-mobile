package subshttp

import (
	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
)

type SubDTORequest struct {
	ServiceName string    `json:"service_name" validate:"required"     example:"Yandex Plus"`
	Price       int       `json:"price" validate:"required"            example:"1200"`
	UserID      uuid.UUID `json:"user_id" validate:"required"          example:"64201fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date" validate:"required,len=7" example:"04-2026"`
	EndDate     *string   `json:"end_date" validate:"omitempty,len=7"  example:"05-2026"`
}

type SubDTOResponse struct {
	ID          uuid.UUID `json:"id"           example:"e38b3587-eac7-4b21-ae11-e932dfa2c907"`
	ServiceName string    `json:"service_name" example:"Yandex Plus"`
	Price       int       `json:"price"        example:"1200"`
	UserID      uuid.UUID `json:"user_id"      example:"64201fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string    `json:"start_date"   example:"04-2026"`
	EndDate     *string   `json:"end_date"     example:"05-2026"`
}

type PatchSubDTORequest struct {
	Price   domains.Nullable[int]    `json:"price" validate:"omitempty"    swaggertype:"integer" example:"1200"`
	EndDate domains.Nullable[string] `json:"end_date" validate:"omitempty" swaggertype:"string" example:"null"`
}

func ToDTO(sub domains.Sub) SubDTOResponse {
	startDate := sub.StartDate.Format("01-2006")
	var endDate *string
	if sub.EndDate != nil {
		date := sub.EndDate.Format("01-2006")
		endDate = &date
	}

	return SubDTOResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

func ToDTOs(subs []domains.Sub) []SubDTOResponse {
	subsDTO := make([]SubDTOResponse, len(subs))
	for i, sub := range subs {
		subsDTO[i] = ToDTO(sub)
	}

	return subsDTO
}
