package router

import "strings"

// Takes a route pattern (string to turn into regex) and a map of all route params in a route handler (ctx.AllParams()) and returns the routePattern with route params replaced with their arguments.
// Example: will replace /users/:id with /users/123 when given map[id:123]
func replaceRouteParams(routePattern string, paramMap map[string]string) string {
	for param, value := range paramMap {
		routePattern = strings.Replace(routePattern, ":"+param, value, -1)
	}

	return routePattern
}
