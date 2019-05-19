// Chef Organization
// Utility functions
package co

import (
        "fmt"
	"regexp"
)

// IsUSAG identifies chef server internal groups.
// These groups are skipped when processing group membership and group reporting.
func IsUSAG(group string) bool {
	match, err := regexp.MatchString("^[0-9a-f]+$", group)
	if err != nil {
		fmt.Println("Issue with regex", err)
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return len(group) == 32 && match
}

// Unique takes an array and returns the unique elements.
func Unique(in []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range in {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
