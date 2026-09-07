package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// semverRegex matches semantic versioning tags like v1.2.3, v1.2.3-beta, v1.2.3+build.
var semverRegex = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)([-+].*)?$`)

// SemVer represents a parsed semantic version.
type SemVer struct {
	Major        int
	Minor        int
	Patch        int
	Suffix       string // Prerelease and/or build metadata suffix (e.g., "-beta+build")
	Full         string
	IsPrerelease bool // True only if suffix starts with "-" (not for build metadata only)
}

// ParseSemVer parses a semantic version tag and returns its components.
func ParseSemVer(tag string) (*SemVer, error) {
	matches := semverRegex.FindStringSubmatch(tag)
	if matches == nil {
		return nil, fmt.Errorf("tag %q does not match semantic versioning format (expected vX.Y.Z)", tag)
	}
	major, minor, patch, err := atoi3(matches[1], matches[2], matches[3])
	if err != nil {
		return nil, fmt.Errorf("tag %q has out of range version numbers: %w", tag, err)
	}
	suffix := matches[4]
	return &SemVer{
		Major:  major,
		Minor:  minor,
		Patch:  patch,
		Suffix: suffix,
		Full:   tag,
		// Per semver spec: prerelease versions have a hyphen suffix (e.g., -beta, -rc.1).
		// Build metadata uses a + suffix (e.g., +build.123) and is NOT a prerelease.
		IsPrerelease: strings.HasPrefix(suffix, "-"),
	}, nil
}

// atoi3 converts the three version number components, failing on the first error.
func atoi3(a, b, c string) (int, int, int, error) {
	x, err := strconv.Atoi(a)
	if err != nil {
		return 0, 0, 0, err
	}
	y, err := strconv.Atoi(b)
	if err != nil {
		return 0, 0, 0, err
	}
	z, err := strconv.Atoi(c)
	if err != nil {
		return 0, 0, 0, err
	}
	return x, y, z, nil
}

// MajorTag returns the major version tag (e.g., "v1").
func (s *SemVer) MajorTag() string {
	return fmt.Sprintf("v%d", s.Major)
}

// MinorTag returns the minor version tag (e.g., "v1.2").
func (s *SemVer) MinorTag() string {
	return fmt.Sprintf("v%d.%d", s.Major, s.Minor)
}

// SemVerGreaterThan returns true if a represents a higher version than b.
func SemVerGreaterThan(a, b *SemVer) bool {
	if a.Major != b.Major {
		return a.Major > b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor > b.Minor
	}
	if a.Patch != b.Patch {
		return a.Patch > b.Patch
	}
	// Same version numbers: a release outranks a prerelease.
	return b.IsPrerelease && !a.IsPrerelease
}
