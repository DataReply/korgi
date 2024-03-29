/*
Copyright © 2020  Artyom Topchyan a.topchyan@reply.de

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func SanitizeAppName(input string) string {
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

func GetRealmeDir() string {
	dir, _ := os.Getwd()
	return ConcatDirs(dir, "realm")
}
func GetNamespaceDir(namespace string) string {
	return ConcatDirs(GetRealmeDir(), "namespaces", namespace)
}
