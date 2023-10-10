package domain

import (
	"fmt"
	"strings"
	"time"
)

type (
	DfDataPoint struct {
		date        time.Time
		major       string
		minor       string
		build       string
		description string
	}
)

func NewDfDataPoint(date time.Time, release string, description string) (*DfDataPoint, error) {
	version, err := parseVersion(release)
	if err != nil {
		return nil, err
	}

	return &DfDataPoint{
		date:        date,
		major:       version[0],
		minor:       version[1],
		build:       version[2],
		description: description,
	}, nil
}

func (d *DfDataPoint) Date() time.Time {
	return d.date
}

func (d *DfDataPoint) Minor() string {
	return d.minor
}

func (d *DfDataPoint) String() string {
	return fmt.Sprintf("release=%s.%s, description=%s", d.major, d.minor, d.description)
}

func (d *DfDataPoint) IsSameVersion(t DfDataPoint) bool {
	return t.major == d.major && t.minor == d.minor
}

func parseVersion(release string) ([]string, error) {
	parts := strings.Split(release, ".")
	if len(parts) >= 3 {
		return parts, nil
	}

	return nil, fmt.Errorf("version is not good: %s", release)
}
