package metrics

import (
	"fmt"
	"sort"

	"dora-metrics/internal/metrics/domain"
	"dora-metrics/internal/metrics/graph"
	"dora-metrics/internal/metrics/repository"
)

type Metric interface {
	Calculate()
}

type DfRepository interface {
	FindDataPoint(projectId int, options repository.FindDataPointOptions) []domain.DfDataPoint
}

type Df struct {
	dfRepository   DfRepository
	chartGenerator *graph.EChartGenerator
}

func NewDf(dfRepository DfRepository, chartGenerator *graph.EChartGenerator) *Df {
	return &Df{
		dfRepository:   dfRepository,
		chartGenerator: chartGenerator,
	}
}

func (m *Df) Calculate(projectID int) {
	// collect data from repository (repository is responsible for gitlab fetching)
	fmt.Println("Collect data from gitlab...")
	data := m.dfRepository.FindDataPoint(projectID, repository.FindDataPointOptions{
		Name:   withString("deploy-k8s:production"),
		Status: withString("success"),
		Page:   200,
	})

	releaseData := domain.NewReleasePeriodByMonth(data)
	sort.Slice(releaseData, func(i, j int) bool {
		return data[i].Date().After(data[j].Date())
	})

	fmt.Println("releases data:")
	for _, item := range releaseData {
		fmt.Printf("%s\n", item.String())
	}

	m.chartGenerator.Render(releaseData, graph.ChartRenderOptions{
		Title: "DORA - Deployment Frequency",
		File:  "text.html",
	})

	fmt.Println("Done!")
}

func withString(v string) *string {
	p := new(string)
	*p = v
	return p
}
