package schedule_usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nikpopkov/running-club/api/internal/domain/dto"
	"github.com/nikpopkov/running-club/api/internal/domain/model"
)

func (u *UseCase) Calendar(ctx context.Context, userID uuid.UUID, role string) (*dto.CalendarResponse, error) {
	club, err := u.clubFor(ctx, userID, role)
	if err != nil {
		if errors.Is(err, model.ErrNotMember) || errors.Is(err, model.ErrNotFound) {
			return dto.NewCalendarResponse([]dto.CalendarCellView{}), nil
		}
		return nil, err
	}

	items, err := u.announceRepo.FindByClub(ctx, club.ID)
	if err != nil {
		return nil, fmt.Errorf("announceRepo.FindByClub: %w", err)
	}

	now := time.Now()
	sessionDays := make(map[int]bool)
	for _, a := range items {
		if a.StartsOn == nil {
			continue
		}
		if a.StartsOn.Year() == now.Year() && a.StartsOn.Month() == now.Month() {
			sessionDays[a.StartsOn.Day()] = true
		}
	}

	accent, onAccent, accentSoft, text := calendarColors(club.AccentHex)
	today := now.Day()
	first := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()
	startDow := (int(first.Weekday()) + 6) % 7

	cells := make([]dto.CalendarCellView, 0, 42)
	for i := 0; i < startDow; i++ {
		cells = append(cells, dto.NewBlankCalendarCell(fmt.Sprintf("b%d", i), "transparent", text, "transparent"))
	}
	for d := 1; d <= daysInMonth; d++ {
		has := sessionDays[d]
		bg, fg, dot := "transparent", text, "transparent"
		if d == today {
			bg, fg = accent, onAccent
		} else if has {
			bg = accentSoft
			dot = accent
		}
		cells = append(cells, dto.NewCalendarCell(fmt.Sprintf("d%d", d), d, has, d == today, bg, fg, dot))
	}
	return dto.NewCalendarResponse(cells), nil
}

func calendarColors(accentHex string) (accent, onAccent, accentSoft, text string) {
	text = "#f4f6f7"
	accent = strings.TrimSpace(accentHex)
	if accent == "" {
		accent = "#ff5c22"
	}
	r, g, b, ok := parseHexRGB(accent)
	if !ok {
		return accent, "#ffffff", "rgba(255,92,34,.15)", text
	}
	lum := (0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 255
	onAccent = "#ffffff"
	if lum > 0.6 {
		onAccent = "#0a0b0c"
	}
	accentSoft = fmt.Sprintf("rgba(%d,%d,%d,.15)", r, g, b)
	return accent, onAccent, accentSoft, text
}

func parseHexRGB(hex string) (r, g, b int, ok bool) {
	h := strings.TrimPrefix(hex, "#")
	if len(h) != 6 {
		return 0, 0, 0, false
	}
	rv, err1 := strconv.ParseInt(h[0:2], 16, 0)
	gv, err2 := strconv.ParseInt(h[2:4], 16, 0)
	bv, err3 := strconv.ParseInt(h[4:6], 16, 0)
	if err1 != nil || err2 != nil || err3 != nil {
		return 0, 0, 0, false
	}
	return int(rv), int(gv), int(bv), true
}
