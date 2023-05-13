package image

import "github.com/avislash/nftstamper/lib/filter"

type OpacityFilter struct {
	Filters map[string]filter.Filter[string, float64]
}

func NewOpacityFilter(opacityFilters map[string]filter.Filter[string, float64]) *OpacityFilter {
	return &OpacityFilter{opacityFilters}
}
