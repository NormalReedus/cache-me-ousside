package utils

type void struct{}

var nothing void

type Set[T comparable] map[T]void

func (set Set[T]) Add(elem T) {
	set[elem] = nothing
}

// Actually not needed for this specific project
func (set Set[T]) Remove(elem T) {
	delete(set, elem)
}

func (set Set[T]) Elements() []T {
	elements := make([]T, 0, len(set))
	for k := range set {
		elements = append(elements, k)
	}
	return elements
}
