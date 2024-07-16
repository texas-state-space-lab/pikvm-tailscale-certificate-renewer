package helpers

import (
	"fmt"
	"regexp"
)

// SetLine sets the line in the contents if the regex matches
func SetLine(contents []string, regex *regexp.Regexp, replacementLine string) []string {
	for i, line := range contents {
		if regex.MatchString(line) {
			contents[i] = replacementLine
			return contents
		}
	}

	return append(contents, fmt.Sprintf("%s\n", replacementLine))
}
