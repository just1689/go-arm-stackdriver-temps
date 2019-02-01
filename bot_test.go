package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestStartReader(t *testing.T) {

	filename := "temp"
	ms := 10

	reply, stop := StartReader(filename, fakeReader, ms)

	go func() {
		time.Sleep(110 * time.Millisecond)
		stop <- true
	}()

	go func() {
		for item := range reply {
			logrus.Info(item)
		}
	}()

	time.Sleep(120 * time.Millisecond)

}

func TestStartAggregator(t *testing.T) {
	ok := false
	filename := "temp"
	ms := 10

	reply, stop := StartReader(filename, fakeReader, ms)
	out := StartAggregator(reply, 10)

	go func() {
		time.Sleep(150 * time.Millisecond)
		stop <- true
	}()

	go func() {
		logrus.Info("Going to listen for items coming back")
		item := <-out
		logrus.Info("Got an item back!")
		b, err := json.Marshal(item)
		if err != nil {
			t.Error(err)
		}
		ok = true
		logrus.Info(string(b))
	}()

	time.Sleep(160 * time.Millisecond)
	if !ok {
		t.Error("Failed to pass the read test")
	}

}

func fakeReader(filename string) (b []byte, err error) {
	b = []byte("49000")
	return
}
