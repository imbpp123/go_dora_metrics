package main

import (
	"log"
	"os"
	"strconv"

	"github.com/xanzy/go-gitlab"

	"dora-metrics/internal/metrics"
	"dora-metrics/internal/metrics/graph"
	"dora-metrics/internal/metrics/repository"
)

func main() {
	git, err := gitlab.NewClient(
		os.Getenv("GITLAB_TOKEN"),
		gitlab.WithBaseURL(os.Getenv("GITLAB_URL")),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	df := metrics.NewDf(
		repository.NewDf(git),
		graph.NewEChartGenerator(),
	)

	projectID, _ := strconv.Atoi(os.Getenv("API_PROJECT_ID"))
	df.Calculate(projectID)
}
