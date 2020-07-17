package pkg

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func ExplodeArg(args []string) []interface{} {
	new := make([]interface{}, len(args))

	for i, v := range args {
		new[i] = v
	}
	return append(new)
}

func SanitzeAppName(input string) string {
	r, _ := regexp.Compile("\\.ya?ml")
	return r.ReplaceAllString(input, "")
}

func ConcatDirs(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
		sb.WriteString("/")
	}
	return filepath.Clean(sb.String())

}
func GetNamespaceDir(namespace string) string {
	dir, _ := os.Getwd()
	return ConcatDirs(dir, "realm", "namespaces", namespace)
}
