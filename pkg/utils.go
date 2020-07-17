package pkg

import (
	"fmt"
	"regexp"
)

func SanitzeAppName(input string) string {
	r, _ := regexp.Compile("\\.ya?ml")
	return r.ReplaceAllString(input, "")
}

func ConcatDir(a string, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}
