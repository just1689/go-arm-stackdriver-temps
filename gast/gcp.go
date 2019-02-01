package gast

import (
	"cloud.google.com/go/monitoring/apiv3"
	"context"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"log"
	"time"
)

func makeDataPoint(v int64) (dataPoint *monitoringpb.Point) {
	dataPoint = &monitoringpb.Point{
		Interval: &monitoringpb.TimeInterval{
			EndTime: &googlepb.Timestamp{
				Seconds: time.Now().Unix(),
			},
		},
		Value: &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: v,
			},
		},
	}
	return
}

func BuildWriter(projectID, deviceID string, client *monitoring.MetricClient, ctx context.Context) func([]*monitoringpb.Point) {
	return func(lp []*monitoringpb.Point) {
		write(projectID, deviceID, client, ctx, lp)
	}
}

func write(projectID, deviceID string, client *monitoring.MetricClient, ctx context.Context, points []*monitoringpb.Point) {
	if err := client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(projectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/stores/temps",
					Labels: map[string]string{
						"device_id": deviceID,
					},
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": projectID,
					},
				},
				Points: points,
			},
		},
	}); err != nil {
		log.Fatalf("Failed to write time series data: %v", err)
	}

}
