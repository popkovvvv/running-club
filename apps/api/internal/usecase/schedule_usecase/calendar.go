package schedule_usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/nikpopkov/running-club/api/internal/domain/dto"
)

func (u *UseCase) Calendar(ctx context.Context, accent, onAccent, accentSoft, text string) (*dto.CalendarResponse, error) {
	sessionDays := map[int]bool{2: true, 7: true, 9: true, 14: true, 16: true, 21: true, 23: true, 28: true, 30: true}
	today := 9
	first := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	startDow := (int(first.Weekday()) + 6) % 7
	cells := make([]dto.CalendarCellView, 0, 42)
	for i := 0; i < startDow; i++ {
		cells = append(cells, dto.NewBlankCalendarCell(fmt.Sprintf("b%d", i), "transparent", text, "transparent"))
	}
	for d := 1; d <= 31; d++ {
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
