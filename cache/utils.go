package cache

//* SET
// Used for finding unique cache entry keys matched by bust-patterns
// (multiple patterns can match same key, we don't need more than one to bust it from cache)

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

// Actually not needed for this specific project
func (set Set[T]) Has(elem T) bool {
	_, ok := set[elem]
	return ok
}

func (set Set[T]) Elements() []T {
	elements := make([]T, 0, len(set))
	for k := range set {
		elements = append(elements, k)
	}
	return elements
}
