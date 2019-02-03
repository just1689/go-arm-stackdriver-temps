package main

import (
	"cloud.google.com/go/monitoring/apiv3"
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/team142/go-arm-stackdriver-temps/gast"
	"log"
	"os"
)

var (
	projectID = flag.String("projectid", "", "A GCP project ID")
	deviceID  = flag.String("deviceid", "", "The label for stackdriver")
	filename  = flag.String("filename", "", "File to read")
)

func main() {

	if projectID == nil || *projectID == "" {
		logrus.Fatal("No projectid provided")
	}
	if deviceID == nil || *deviceID == "" {
		logrus.Fatal("No deviceid provided")
	}
	if filename == nil || *filename == "" {
		logrus.Fatal("No filename provided")
	}

	sigs := make(chan os.Signal, 1)

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

	<-sigs
	stop <- true

	// Closes the client and flushes the data to Stackdriver.
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}

	fmt.Printf("Done writing time series data.\n")
}
