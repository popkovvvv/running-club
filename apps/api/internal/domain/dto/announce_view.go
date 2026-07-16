package dto

import "github.com/google/uuid"

type AnnounceView struct {
	ID          uuid.UUID         `json:"id"`
	Place       string            `json:"place"`
	Day         string            `json:"day"`
	Time        string            `json:"time"`
	Group       string            `json:"group"`
	Note        string            `json:"note"`
	Going       int               `json:"going"`
	SignedUp    bool              `json:"signedUp"`
	ScheduleCta string            `json:"scheduleCta"`
	StartsOn    string            `json:"startsOn,omitempty"`
	Attendees   []GoingPersonView `json:"attendees"`
}

type GoingPersonView struct {
	ID   uuid.UUID `json:"id"`
	Init string    `json:"init"`
	Name string    `json:"name"`
}

func NewGoingPersonView(id uuid.UUID, init, name string) GoingPersonView {
	return GoingPersonView{ID: id, Init: init, Name: name}
}

func NewAnnounceView(
	id uuid.UUID,
	place, day, timeLabel, group, note string,
	going int,
	signedUp bool,
	scheduleCta string,
) AnnounceView {
	return AnnounceView{
		ID:          id,
		Place:       place,
		Day:         day,
		Time:        timeLabel,
		Group:       group,
		Note:        note,
		Going:       going,
		SignedUp:    signedUp,
		ScheduleCta: scheduleCta,
		Attendees:   []GoingPersonView{},
	}
}

func (v AnnounceView) WithStartsOn(iso string) AnnounceView {
	v.StartsOn = iso
	return v
}

func (v AnnounceView) WithAttendees(attendees []GoingPersonView) AnnounceView {
	v.Attendees = attendees
	if len(attendees) > v.Going {
		v.Going = len(attendees)
	}
	return v
}

type CreateAnnounceRequest struct {
	Place    string `json:"place"`
	Day      string `json:"day"`
	Time     string `json:"time"`
	Group    string `json:"group"`
	Note     string `json:"note"`
	StartsOn string `json:"startsOn,omitempty"`
}

type CalendarResponse struct {
	Year  int                `json:"year"`
	Month int                `json:"month"`
	Label string             `json:"label"`
	Cells []CalendarCellView `json:"cells"`
}

func NewCalendarResponse(year, month int, label string, cells []CalendarCellView) *CalendarResponse {
	return &CalendarResponse{Year: year, Month: month, Label: label, Cells: cells}
}

type CalendarCellView struct {
	Key     string `json:"key"`
	N       int    `json:"n"`
	Blank   bool   `json:"blank"`
	Has     bool   `json:"has"`
	IsToday bool   `json:"isToday"`
	Iso     string `json:"iso,omitempty"`
	Bg      string `json:"bg"`
	Fg      string `json:"fg"`
	Dot     string `json:"dot"`
}

func NewBlankCalendarCell(key, bg, fg, dot string) CalendarCellView {
	return CalendarCellView{
		Key:   key,
		Blank: true,
		Bg:    bg,
		Fg:    fg,
		Dot:   dot,
	}
}

func NewCalendarCell(key string, n int, has, isToday bool, iso, bg, fg, dot string) CalendarCellView {
	return CalendarCellView{
		Key:     key,
		N:       n,
		Blank:   false,
		Has:     has,
		IsToday: isToday,
		Iso:     iso,
		Bg:      bg,
		Fg:      fg,
		Dot:     dot,
	}
}
