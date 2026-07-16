package polyline_test

import (
	"testing"

	"github.com/nikpopkov/running-club/api/internal/pkg/polyline"
	"github.com/stretchr/testify/require"
)

func TestToSVG(t *testing.T) {
	t.Parallel()
	svg := polyline.ToSVG("_p~iF~ps|U_ulLnnqC_mqNvxq`@", 300, 140)
	require.NotEmpty(t, svg.Path)
	require.Greater(t, svg.SX, 0.0)
	require.Greater(t, svg.EX, 0.0)
}
