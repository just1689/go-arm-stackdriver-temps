package gast

import (
	"fmt"
	"github.com/sirupsen/logrus"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	"strconv"
	"time"
)

func StartReader(filename string, f func(string) ([]byte, error), ms int) (reply chan int, stop chan bool) {

	reply = make(chan int, 32)
	stop = make(chan bool)
	d := time.Duration(ms) * time.Millisecond

	go func() {
	outer:
		for {
			select {
			case <-stop:
				logrus.Info("Stopping reader")
				close(reply)
				close(stop)
				break outer
			case <-time.After(d):
				dat, err := f(filename)
				if err != nil {
					logrus.Error(err)
					break
				}
				str := string(dat)
				i, err := strconv.Atoi(str)
				if err != nil {
					logrus.Error(err)
					break
				}
				reply <- i
			}
		}
	}()
	return
}

func StartAggregator(in chan int, batchSize int) (out chan []*monitoringpb.Point) {
	d := time.Duration(60) * time.Second
	out = make(chan []*monitoringpb.Point)

	var slice []int
	go func() {
		for {
			select {
			case <-time.After(d):
				if len(slice) > 0 {
					convertAndSend(slice, out)
					slice = slice[:0]
				}
			case i := <-in:
				fmt.Println("Read: ", i)
				slice = append(slice, i)
				if len(slice) >= batchSize {
					convertAndSend(slice, out)
					slice = slice[:0]
				}
			}
		}
	}()
	return
}

func StartWriter(in chan []*monitoringpb.Point, f func(points []*monitoringpb.Point)) {
	go func() {
		for ps := range in {
			f(ps)
		}
	}()
	return
}

func convertAndSend(sl []int, out chan []*monitoringpb.Point) {
	result := make([]*monitoringpb.Point, len(sl))
	for i, s := range sl {
		result[i] = makeDataPoint(int64(s))
	}
	out <- result
}
