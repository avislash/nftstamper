package filter

type Filter[K comparable, V any] struct {
	Name    string
	Default V
	Weights map[K]V
}
