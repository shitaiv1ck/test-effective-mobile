package subshttp

import (
	"github.com/google/uuid"
	"github.com/shitaiv1ck/test-effective-mobile/internal/core/domains"
)

type SubDTORequest struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required,len=7"`
	EndDate     *string   `json:"end_date" validate:"omitempty,len=7"`
}

type SubDTOResponse struct {
	ID          uuid.UUID
	ServiceName string
	Price       int
	UserID      uuid.UUID
	StartDate   string
	EndDate     *string `json:"end_date" validate:"required,len=7"`
}

type PatchSubDTORequest struct {
	Price   domains.Nullable[int]    `json:"price" validate:"omitempty"`
	EndDate domains.Nullable[string] `json:"end_date" validate:"omitempty"`
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
