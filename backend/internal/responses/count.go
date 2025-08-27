package responses

type CountResponse struct {
	Count int `json:"count"`
}

func NewCountResponse(count int) *CountResponse {
	return &CountResponse{
		Count: count,
	}
}
