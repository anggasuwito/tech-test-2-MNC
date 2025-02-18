package entity

type (
	UserAccount struct {
		ID          string `json:"id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}

	UpdateProfileRequest struct {
		AccountID string `json:"-"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address   string `json:"address"`
	}

	UpdateProfileResponse struct {
		*UserAccount
	}
)
