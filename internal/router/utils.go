package router

import (
	"strings"
)

// Takes a route pattern (string to turn into regex) and a map of all route params in a route handler (ctx.AllParams()) and returns the routePattern with route params replaced with their arguments.
// Example: will replace /users/:id with /users/123 when given map[id:123]
func replaceRouteParams(paramMap map[string]string, routePatternTemplates []string) []string {
	// Copy original slice so we can return a new one
	// newRoutePatterns := routePatternTemplates
	newRoutePatterns := make([]string, len(routePatternTemplates))

	// We must copy slice to avoid returning a reference to the original, underlying array on consecutive requests
	copy(newRoutePatterns, routePatternTemplates)

	for param, value := range paramMap {
		for i, pattern := range newRoutePatterns {
			newRoutePatterns[i] = strings.Replace(pattern, ":"+param, value, -1)
		}
	}

	return newRoutePatterns
}
