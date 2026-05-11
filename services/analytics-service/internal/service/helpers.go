package service

import (
	"encoding/json"
	"sort"

	"github.com/finalops/analytics-service/internal/model"
)

func jsonMarshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func sortProgressionByDate(points []model.ProgressionPoint) {
	sort.Slice(points, func(i, j int) bool {
		return points[i].Date < points[j].Date
	})
}
