package domain

type RegisterRequest struct {
	Username string `validate:"max=20,min=3"`
	Password string `validate:"max=20,min=6"`
}

type RegisterResponse struct {
	Status bool `json:"status"`
}

type AuthRequest struct {
	Username string `validate:"max=20,min=3"`
	Password string `validate:"max=20,min=6"`
}

type AuthResponse struct {
	Status bool   `json:"status"`
	Token  string `json:"token,omitempty"`
}

type AuthGameRequest struct {
	Token string `json:"token"`
}

type AuthGameResponse struct {
	Status   bool   `json:"status"`
	Username string `json:"username,omitempty"`
	Message  string `json:"message,omitempty"`
}

type GetUsersRequest struct {
	Limit  int `validate:"max=100,min=0"`
	Offset int `validate:"min=0"`
}

type GetUsersResponse struct {
	Status bool         `json:"status"`
	Users  []*DummyUser `json:"users"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}
