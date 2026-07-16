package polyline

import (
	"fmt"
	"math"
)

func Decode(encoded string) ([][2]float64, error) {
	if encoded == "" {
		return nil, nil
	}
	var coords [][2]float64
	index, lat, lng := 0, 0, 0
	for index < len(encoded) {
		var result, shift int
		for {
			if index >= len(encoded) {
				return nil, fmt.Errorf("truncated polyline")
			}
			b := int(encoded[index]) - 63
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		dlat := result >> 1
		if result&1 != 0 {
			dlat = ^dlat
		}
		lat += dlat

		result, shift = 0, 0
		for {
			if index >= len(encoded) {
				return nil, fmt.Errorf("truncated polyline")
			}
			b := int(encoded[index]) - 63
			index++
			result |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		dlng := result >> 1
		if result&1 != 0 {
			dlng = ^dlng
		}
		lng += dlng

		coords = append(coords, [2]float64{float64(lat) / 1e5, float64(lng) / 1e5})
	}
	return coords, nil
}

type SVGPath struct {
	Path string
	SX   float64
	SY   float64
	EX   float64
	EY   float64
}

func ToSVG(encoded string, width, height float64) SVGPath {
	coords, err := Decode(encoded)
	if err != nil || len(coords) < 2 {
		return SVGPath{}
	}
	minLat, maxLat := coords[0][0], coords[0][0]
	minLng, maxLng := coords[0][1], coords[0][1]
	for _, c := range coords[1:] {
		minLat = math.Min(minLat, c[0])
		maxLat = math.Max(maxLat, c[0])
		minLng = math.Min(minLng, c[1])
		maxLng = math.Max(maxLng, c[1])
	}
	latSpan := maxLat - minLat
	lngSpan := maxLng - minLng
	if latSpan == 0 {
		latSpan = 1e-6
	}
	if lngSpan == 0 {
		lngSpan = 1e-6
	}
	padding := 12.0
	innerW := width - padding*2
	innerH := height - padding*2

	project := func(c [2]float64) (float64, float64) {
		x := padding + (c[1]-minLng)/lngSpan*innerW
		y := padding + (maxLat-c[0])/latSpan*innerH
		return x, y
	}

	sx, sy := project(coords[0])
	path := fmt.Sprintf("M %.1f %.1f", sx, sy)
	for _, c := range coords[1:] {
		x, y := project(c)
		path += fmt.Sprintf(" L %.1f %.1f", x, y)
	}
	ex, ey := project(coords[len(coords)-1])
	return SVGPath{Path: path, SX: sx, SY: sy, EX: ex, EY: ey}
}
