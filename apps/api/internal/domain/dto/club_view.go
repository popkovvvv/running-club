package dto

import "github.com/google/uuid"

type ClubView struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	InviteCode string    `json:"inviteCode,omitempty"`
	AccentHex  string    `json:"accentHex"`
	Students   int       `json:"students"`
}

func NewClubView(id uuid.UUID, name, accentHex string, students int) *ClubView {
	return &ClubView{
		ID:        id,
		Name:      name,
		AccentHex: accentHex,
		Students:  students,
	}
}

func (v *ClubView) WithInviteCode(code string) *ClubView {
	v.InviteCode = code
	return v
}

type JoinClubRequest struct {
	Code string `json:"code"`
}

type PaletteRequest struct {
	AccentHex string `json:"accentHex"`
}

type StudentView struct {
	ID   uuid.UUID `json:"id"`
	Init string    `json:"init"`
	Name string    `json:"name"`
	Sub  string    `json:"sub"`
	Km   string    `json:"km"`
	Comp int       `json:"comp"`
}

func NewStudentView(id uuid.UUID, init, name, sub, km string, comp int) StudentView {
	return StudentView{
		ID:   id,
		Init: init,
		Name: name,
		Sub:  sub,
		Km:   km,
		Comp: comp,
	}
}
