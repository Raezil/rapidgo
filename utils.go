package rapidgo

import (
	"fmt"
	"os"
	"strings"
)

func ResolvePort(addr []string) string {
	if len(addr) == 0 {
		if port := os.Getenv("PORT"); port != "" {
			return ":" + port
		}
		fmt.Println("PORT environment variable not set. Using default port 8080")
		return ":8080"
	} else if len(addr) == 1 {
		return ":" + addr[0]
	} else {
		panic("Too many parameters")
	}
}

// Check if the path is dynamic
func IsDynamic(path string) bool {
	if len(path) > 0 && (strings.Contains(path, ":") || strings.Contains(path, "*")) {
		return true
	}
	return false
}

// Generate Static Route Key value
func GenerateStaticRouteKey(method, path string) string {
	return method + "#" + path
}
