package config

// Check if a slice contains a value.
// Is used to remove invalid method keys in Cache and Bust maps on the config.
func contains[T comparable](slice []T, target T) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}
