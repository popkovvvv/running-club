package dto

import "github.com/google/uuid"

type AnnounceView struct {
	ID          uuid.UUID `json:"id"`
	Place       string    `json:"place"`
	Day         string    `json:"day"`
	Time        string    `json:"time"`
	Group       string    `json:"group"`
	Note        string    `json:"note"`
	Going       int       `json:"going"`
	SignedUp    bool      `json:"signedUp"`
	ScheduleCta string    `json:"scheduleCta"`
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
	}
}

type CreateAnnounceRequest struct {
	Place string `json:"place"`
	Day   string `json:"day"`
	Time  string `json:"time"`
	Group string `json:"group"`
	Note  string `json:"note"`
}

type CalendarResponse struct {
	Cells []CalendarCellView `json:"cells"`
}

func NewCalendarResponse(cells []CalendarCellView) *CalendarResponse {
	return &CalendarResponse{Cells: cells}
}

type CalendarCellView struct {
	Key     string `json:"key"`
	N       int    `json:"n"`
	Blank   bool   `json:"blank"`
	Has     bool   `json:"has"`
	IsToday bool   `json:"isToday"`
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

func NewCalendarCell(key string, n int, has, isToday bool, bg, fg, dot string) CalendarCellView {
	return CalendarCellView{
		Key:     key,
		N:       n,
		Blank:   false,
		Has:     has,
		IsToday: isToday,
		Bg:      bg,
		Fg:      fg,
		Dot:     dot,
	}
}
