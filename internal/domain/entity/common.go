package entity

type JWTClaimAccountInfo struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address   string `json:"address"`
}

type JWTClaim struct {
	ID          string               `json:"id"`
	ExpiredAt   int64                `json:"expired_at"`
	AccountInfo *JWTClaimAccountInfo `json:"account_info"`
}
type Filter struct {
	Field string `json:"field"`
	Value string `json:"value"`
}
type ListPaginationRequest struct {
	Limit  int64     `json:"limit"`
	Page   int64     `json:"page"`
	Search []*Filter `json:"search"`
	Sort   []*Filter `json:"sort"`
}

type ListPaginationResponse struct {
	CurrentPage int64 `json:"current_page"`
	TotalPage   int64 `json:"total_page"`
	TotalData   int64 `json:"total_data"`
	PerPage     int64 `json:"per_page"`
}
