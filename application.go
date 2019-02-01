package main

import (
	"cloud.google.com/go/monitoring/apiv3"
	"context"
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	projectID = flag.String("projectid", "", "A GCP project ID")
	deviceID  = flag.String("deviceid", "", "The label for stackdriver")
	filename  = flag.String("filename", "temp", "File to read")
)

func main() {

	flag.Parse()

	ctx := context.Background()

	// Creates a client.
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	readings, stop := StartReader(*filename, readFile, 1000)
	batches := StartAggregator(readings, 10)
	StartWriter(batches, buildWriter(client, ctx))

	time.After(2 * time.Minute)
	stop <- true

	// Closes the client and flushes the data to Stackdriver.
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}

	fmt.Printf("Done writing time series data.\n")
}
