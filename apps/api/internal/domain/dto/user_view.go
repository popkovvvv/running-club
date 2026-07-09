package dto

import "github.com/google/uuid"

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserView `json:"user"`
}

func NewAuthResponse(token string, user UserView) *AuthResponse {
	return &AuthResponse{Token: token, User: user}
}

type UserView struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	InClub    bool       `json:"inClub"`
	NeedsClub bool       `json:"needsClub"`
	ClubID    *uuid.UUID `json:"clubId,omitempty"`
}

func NewUserView(id uuid.UUID, name, email, role string) *UserView {
	return &UserView{
		ID:    id,
		Name:  name,
		Email: email,
		Role:  role,
	}
}

func (v *UserView) WithClub(clubID uuid.UUID) *UserView {
	v.InClub = true
	v.ClubID = &clubID
	return v
}

func (v *UserView) MarkInClub() *UserView {
	v.InClub = true
	return v
}

func (v *UserView) MarkNeedsClub() *UserView {
	v.NeedsClub = true
	v.InClub = false
	return v
}
