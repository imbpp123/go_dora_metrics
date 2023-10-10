package domain

import (
	"fmt"
	"time"
)

type ReleasePeriod struct {
	date     time.Time
	releases []DfDataPoint
}

func NewReleasePeriodByMonthFromDataPoint(data DfDataPoint) *ReleasePeriod {
	year := data.Date().Year()
	month := data.Date().Month()

	return &ReleasePeriod{
		date:     time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
		releases: []DfDataPoint{data},
	}
}

func NewReleasePeriodByMonth(data []DfDataPoint) (result []*ReleasePeriod) {
	versionIndex := make(map[string]int)

	for _, item := range data {
		if _, ok := versionIndex[item.Minor()]; ok {
			// this version is added already
			continue
		}
		versionIndex[item.Minor()] = 1

		release := NewReleasePeriodByMonthFromDataPoint(item)
		if curr := findReleaseByDate(result, release.date); curr == nil {
			result = append(result, release)
		} else {
			curr.AddRelease(item)
		}
	}

	return
}

func (r *ReleasePeriod) Date() time.Time {
	return r.date
}

func (r *ReleasePeriod) Qty() int {
	return len(r.releases)
}

func (r *ReleasePeriod) String() string {
	text := fmt.Sprintf("date=%s, qty=%d, releases:\n", r.date.Format(time.DateOnly), r.Qty())
	for _, point := range r.releases {
		text = text + "\t" + point.String() + "\n"
	}
	return text
}

func (r *ReleasePeriod) AddRelease(item DfDataPoint) {
	for _, point := range r.releases {
		if point.IsSameVersion(item) {
			return
		}
	}

	r.releases = append(r.releases, item)
}

func findReleaseByDate(arr []*ReleasePeriod, date time.Time) *ReleasePeriod {
	for _, d := range arr {
		if d.date.Equal(date) {
			return d
		}
	}

	return nil
}
