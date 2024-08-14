package model

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// OAuth User
type UserLoginRequest struct {
	Password    string `json:"password" validate:"required"`
	Email       string `json:"email" validate:"omitempty,email"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,phone_number"`
}
type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type UserLogoutResponse struct {
}

// OAuth Workspace
type WorkspaceLoginRequest struct {
	WorkspaceId string `json:"workspace_id" validate:"required"`
}
type WorkspaceLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type WorkspaceLogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type WorkspaceLogoutResponse struct {
}
