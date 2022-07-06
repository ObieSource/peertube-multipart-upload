package main

import (
	"testing"
)

func TestVerifyEnvironmentValues(t *testing.T) {
	VerifyEnvironmentValues(MultipartUploadHandlerHandlerInput{Hostname: "https://peertube.cpy.re"})
	t.Fail()
}
