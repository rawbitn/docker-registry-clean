package utils

import (
	"sort"
	"strconv"
	"strings"
)

type versionParts struct {
	major int
	minor int
	patch int
	build int
}

func CustomVersionSort(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		v1 := extractVersionParts(versions[i])
		v2 := extractVersionParts(versions[j])
		return isLessVersion(v1, v2)
	})
}

func extractVersionParts(version string) versionParts {
	parts := strings.Split(version, "_")
	var int_version versionParts
	if len(parts) == 3 {
		version_part := strings.Split(parts[1], ".")
		int_version.major, _ = strconv.Atoi(version_part[0])
		int_version.minor, _ = strconv.Atoi(version_part[1])
		int_version.patch, _ = strconv.Atoi(version_part[2])
		int_version.build, _ = strconv.Atoi(parts[2])
	}
	return int_version
}

func isLessVersion(v1, v2 versionParts) bool {
	switch {
	case v1.major < v2.major:
		return true
	case v1.major > v2.major:
		return false
	case v1.minor < v2.minor:
		return true
	case v1.minor > v2.minor:
		return false
	case v1.patch < v2.patch:
		return true
	case v1.patch > v2.patch:
		return false
	case v1.build < v2.build:
		return true
	}
	return false
}
