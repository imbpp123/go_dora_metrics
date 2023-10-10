package repository

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/xanzy/go-gitlab"

	"dora-metrics/internal/metrics/domain"
)

type FindDataPointOptions struct {
	Page   int
	Name   *string
	Status *string
}

type Df struct {
	client *gitlab.Client
}

func NewDf(client *gitlab.Client) *Df {
	return &Df{
		client: client,
	}
}

func (r *Df) FindDataPoint(projectId int, options FindDataPointOptions) (data []domain.DfDataPoint) {
	var wg sync.WaitGroup

	page := 0
	limit := 100
	dataPointChannel := make(chan domain.DfDataPoint, 1000)
	limitChannel := make(chan struct{}, limit)

	go func() {
		for dataItem := range dataPointChannel {
			fmt.Printf("\t\t\tfound point\n")
			data = append(data, dataItem)
		}
	}()

	for {
		wg.Add(1)
		go func(currentPage int, ch chan<- domain.DfDataPoint, wg *sync.WaitGroup) {
			defer wg.Done()

			pipelines := r.findPipelines(projectId, currentPage)
			fmt.Printf("page: %d\n", currentPage)
			for _, item := range pipelines {
				pItem := item

				wg.Add(1)
				limitChannel <- struct{}{}
				go func(item *gitlab.PipelineInfo, ch chan<- domain.DfDataPoint, wg *sync.WaitGroup) {
					defer func() {
						<-limitChannel
						wg.Done()
					}()

					bridges := r.findJobsBridges(projectId, item.ID)
					for _, bridgeItem := range bridges {
						if options.Name != nil && *options.Name != bridgeItem.Name {
							continue
						}
						if options.Status != nil && *options.Status != bridgeItem.Status {
							continue
						}
						if bridgeItem.FinishedAt == nil {
							continue
						}

						if newPoint, err := domain.NewDfDataPoint(*bridgeItem.FinishedAt, bridgeItem.Ref, bridgeItem.Commit.Title); err == nil {
							fmt.Printf("\t\tadd bridge: %s\n", bridgeItem.WebURL)
							ch <- *newPoint
						}
					}
				}(pItem, ch, wg)

				wg.Add(1)
				limitChannel <- struct{}{}
				go func(item *gitlab.PipelineInfo, ch chan<- domain.DfDataPoint, wg *sync.WaitGroup) {
					defer func() {
						<-limitChannel
						wg.Done()
					}()

					jobs := r.findJobs(projectId, item.ID)
					for _, job := range jobs {

						if job.Name != "deploy:production" {
							continue
						}
						if options.Status != nil && *options.Status != job.Status {
							continue
						}
						if job.FinishedAt == nil {
							continue
						}

						if newPoint, err := domain.NewDfDataPoint(*job.FinishedAt, job.Ref, job.Commit.Title); err == nil {
							fmt.Printf("\t\tadd bridge: %s\n", job.WebURL)
							ch <- *newPoint
						}
					}
				}(pItem, ch, wg)
			}
		}(page, dataPointChannel, &wg)

		page++
		if page > options.Page {
			break
		}
	}

	wg.Wait()
	close(dataPointChannel)

	sort.Slice(data, func(i, j int) bool {
		return data[i].Date().After(data[j].Date())
	})

	return
}

func (r *Df) findPipelines(projectID int, page int) []*gitlab.PipelineInfo {
	pipelines, _, err := r.client.Pipelines.ListProjectPipelines(
		projectID,
		&gitlab.ListProjectPipelinesOptions{
			ListOptions: gitlab.ListOptions{
				Page: page,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return pipelines
}

func (r *Df) findJobsBridges(projectID int, pipelineID int) []*gitlab.Bridge {
	bridges, _, err := r.client.Jobs.ListPipelineBridges(
		projectID,
		pipelineID,
		&gitlab.ListJobsOptions{},
	)
	if err != nil {
		fmt.Printf("error happened while fetching bridges: %s\n", err.Error())
	}

	return bridges
}

func (r *Df) findJobs(projectID int, pipelineID int) []*gitlab.Job {
	jobs, _, err := r.client.Jobs.ListPipelineJobs(
		projectID,
		pipelineID,
		&gitlab.ListJobsOptions{},
	)
	if err != nil {
		fmt.Printf("error happened while fetching jobs: %s\n", err.Error())
	}

	return jobs
}
