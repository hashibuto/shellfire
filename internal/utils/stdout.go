package utils

import (
	"io/ioutil"
	"os"
)

// CapturesStdout captures stdout from the wrapped function
func CaptureStdout(wrapped func() error) ([]byte, error) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	defer func() {
		os.Stdout = oldStdout
	}()
	os.Stdout = w

	err := wrapped()
	if err != nil {
		return nil, err
	}

	w.Close()
	return ioutil.ReadAll(r)
}
