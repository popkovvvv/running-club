//go:build unit

package schedule_usecase_test

import (
	"context"
	"testing"

	"github.com/nikpopkov/running-club/api/internal/usecase/schedule_usecase"
	"github.com/stretchr/testify/require"
)

func TestCalendar(t *testing.T) {
	t.Parallel()
	m := newMocks(t)
	uc := schedule_usecase.NewUseCase(m.announceRepo, m.clubRepo, m.membershipRepo)
	res, err := uc.Calendar(context.Background(), "#ff5c22", "#fff", "soft", "#f4f6f7")
	require.NoError(t, err)
	require.NotEmpty(t, res.Cells)
	var blanks int
	for _, c := range res.Cells {
		if c.Blank {
			blanks++
			require.Equal(t, 0, c.N)
			require.Equal(t, "transparent", c.Dot)
			require.Equal(t, "transparent", c.Bg)
		}
	}
	require.Greater(t, blanks, 0)
}
