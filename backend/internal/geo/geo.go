// Package geo memeriksa apakah sebuah titik (alamat pembeli / titik penjemputan toko)
// berada di dalam wilayah layanan (service_areas.boundary, GeoJSON Polygon).
package geo

import (
	"encoding/json"
	"fmt"
)

type Polygon struct {
	Type        string         `json:"type"`
	Coordinates [][][2]float64 `json:"coordinates"` // [ring][point][lng, lat]
}

func ParsePolygon(raw []byte) (Polygon, error) {
	var p Polygon
	if err := json.Unmarshal(raw, &p); err != nil {
		return p, err
	}
	if p.Type != "Polygon" || len(p.Coordinates) == 0 || len(p.Coordinates[0]) < 4 {
		return p, fmt.Errorf("boundary harus GeoJSON Polygon dengan minimal 4 titik")
	}
	return p, nil
}

// Contains memakai ray casting pada ring luar; ring berikutnya dianggap lubang.
func (p Polygon) Contains(lat, lng float64) bool {
	if len(p.Coordinates) == 0 || !inRing(p.Coordinates[0], lat, lng) {
		return false
	}
	for _, hole := range p.Coordinates[1:] {
		if inRing(hole, lat, lng) {
			return false
		}
	}
	return true
}

func inRing(ring [][2]float64, lat, lng float64) bool {
	in := false
	for i, j := 0, len(ring)-1; i < len(ring); j, i = i, i+1 {
		xi, yi := ring[i][0], ring[i][1]
		xj, yj := ring[j][0], ring[j][1]
		if (yi > lat) != (yj > lat) && lng < (xj-xi)*(lat-yi)/(yj-yi)+xi {
			in = !in
		}
	}
	return in
}
