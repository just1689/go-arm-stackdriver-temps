package main

import (
	"cloud.google.com/go/monitoring/apiv3"
	"context"
	"flag"
	"fmt"
	"github.com/team142/go-arm-stackdriver-temps/gast"
	"log"
	"time"
)

var (
	projectID = flag.String("projectid", "ex-remote-pi", "A GCP project ID")
	deviceID  = flag.String("deviceid", "desktop", "The label for stackdriver")
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

	readings, stop := gast.StartReader(*filename, gast.ReadFile, 60000)
	batches := gast.StartAggregator(readings, 1)
	gast.StartWriter(batches, gast.BuildWriter(*projectID, *deviceID, client, ctx))

	<-time.After(10 * time.Minute)
	stop <- true

	// Closes the client and flushes the data to Stackdriver.
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}

	fmt.Printf("Done writing time series data.\n")
}
