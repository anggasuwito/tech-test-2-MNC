package entity

type (
	AuthRegisterRequest struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		PIN         string `json:"pin"`
	}

	AuthRegisterResponse struct {
		*UserAccount
	}

	AuthLoginRequest struct {
		PhoneNumber string `json:"phone_number"`
		PIN         string `json:"pin"`
	}

	AuthLoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
