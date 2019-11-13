package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic/v7"
)

type Service struct {
	Config struct {
		Prefixes   []string
		DaysToKeep int
	}

	client *elastic.Client
}

// Defaults that can be overridden.
const DefaultElasticsearchURL = "http://elasticsearch:9200"
const DefaultDelay = 43200 * time.Second
const DefaultDaysToKeep = 7
const DefaultMaxRetries = 5
const DefaultRunOnce = false

func (svc *Service) Init() error {
	var err error
	var elasticsearchUrl string
	var delay time.Duration
	var daysToKeep int
	var maxRetries int
	var runOnce bool

	// Check if there is a ELASTICSEARCH_URL passed in.
	if os.Getenv("ELASTICSEARCH_URL") == "" {
		elasticsearchUrl = DefaultElasticsearchURL
	} else {
		elasticsearchUrl = os.Getenv("ELASTICSEARCH_URL")
	}

	// Check if there is a DELAY passed in.
	if os.Getenv("DELAY") == "" {
		delay = DefaultDelay
	} else {
		i, err := strconv.Atoi(os.Getenv("DELAY"))
		if err != nil {
			return err
		}

		delay = time.Duration(i) * time.Second
	}

	// Check if there is a DAYS_TO_KEEP passed in.
	if os.Getenv("DAYS_TO_KEEP") == "" {
		daysToKeep = DefaultDaysToKeep
	} else {
		i, err := strconv.Atoi(os.Getenv("DAYS_TO_KEEP"))
		if err != nil {
			return err
		}

		daysToKeep = i
	}

	// Check if there is a MAX_RETRIES passed in.
	if os.Getenv("MAX_RETRIES") == "" {
		maxRetries = DefaultMaxRetries
	} else {
		maxRetries, err = strconv.Atoi(os.Getenv("MAX_RETRIES"))
		if err != nil {
			return err
		}
	}

	// Check if there is a RUN_ONCE passed in.
	if os.Getenv("RUN_ONCE") == "" {
		runOnce = DefaultRunOnce
	} else {
		runOnce, err = strconv.ParseBool(os.Getenv("RUN_ONCE"))
		if err != nil {
			return err
		}
	}

	// Check if there is a PREFIXES passed in.
	if os.Getenv("PREFIXES") == "" {
		return errors.New("Prefixes is a required field.")
	}

	// Setup the elastic client.
	client, err := elastic.NewClient(
		elastic.SetURL(elasticsearchUrl),
		elastic.SetSniff(false),
		elastic.SetMaxRetries(maxRetries),
	)
	if err != nil {
		return err
	}

	ctx := context.Background()

	svc.client = client
	svc.Config.Prefixes = strings.Split(os.Getenv("PREFIXES"), ",")
	svc.Config.DaysToKeep = daysToKeep

	// Continuously loop on a delay.
	for {
		if err := svc.Run(); err != nil {
			return err
		}

		// Run once mode for unit tests.
		if runOnce {
			return nil
		}

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil
		}
	}
}

func (svc *Service) Run() error {
	pruneTime := time.Now().Add(-24 * time.Duration(svc.Config.DaysToKeep) * time.Hour)

	// Loop through all the prefixes that need to be pruned.
	for _, prefix := range svc.Config.Prefixes {
		// Get all of the indices with the given prefix.
		indices, err := svc.client.IndexGet(fmt.Sprintf("%s-*", string(prefix))).Do(context.Background())
		if err != nil {
			return err
		}

		// Loop through indices and check the time to see if it is before the pruneTime.
		for _, i := range indices {
			index := i.Settings["index"].(map[string]interface{})["provided_name"].(string)

			indexTime, err := time.Parse("2006-01-02", index[len(prefix)+1:])
			if err != nil {
				fmt.Println(fmt.Sprintf("Error parsing index: %s", index))
				continue
			}

			// If the indexTime is before the pruneTime then it is safe to delete the index.
			if indexTime.Before(pruneTime) {
				_, err := svc.client.DeleteIndex(index).Do(context.Background())
				if err != nil {
					fmt.Println(fmt.Sprintf("Failed deleting index: %s", index))
				}

				fmt.Println(fmt.Sprintf("Deleted index: %s", index))
			} else {
				fmt.Println(fmt.Sprintf("Skipped index: %s", index))
			}
		}
	}

	return nil
}
