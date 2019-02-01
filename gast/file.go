package gast

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func ReadFile(filename string) (b []byte, err error) {
	b, err = ioutil.ReadFile(filename)
	if err != nil {
		logrus.Error(err)
	}
	return
}
