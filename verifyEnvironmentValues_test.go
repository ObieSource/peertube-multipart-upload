package main

import (
	"testing"
)

var verifyEnvironmentValuesTestsSuccess []MultipartUploadHandlerHandlerInput = []MultipartUploadHandlerHandlerInput{
	// the four values that are checked, and the hostname
	{
		Hostname: "https://peertube.cpy.re",
		Category: 1,
		Privacy:  1,
		Licence:  4,
		Language: "en",
	},
}

type vevtf struct {
	Input  MultipartUploadHandlerHandlerInput
	Errors int
}

var verifyEnvironmentValuesTestsFails []vevtf = []vevtf{
	// the four values that are checked, and the hostname
	{
		MultipartUploadHandlerHandlerInput{
			Hostname: "https://peertube.cpy.re",
			Category: 0,  // wrong
			Privacy:  12, // wrong
			Licence:  4,
			Language: "en",
		}, 2},
	{
		MultipartUploadHandlerHandlerInput{
			Hostname: "https://peertube.cpy.re",
			Category: 200, // wrong
			Privacy:  1,
			Licence:  15,               // wrong
			Language: "not a language", // wrong
		}, 3},
	{
		MultipartUploadHandlerHandlerInput{
			Hostname: "https://peertube.cpy.re",
			Category: 1,
			Privacy:  1,
			Licence:  50, // wrong
			Language: "en",
		}, 1},
}

func TestVerifyEnvironmentValuesSuccess(t *testing.T) {
	// verify that it succeeds.
	for i, c := range verifyEnvironmentValuesTestsSuccess {
		err, correct, failtext := VerifyEnvironmentValues(c)

		if err != nil {
			t.Skipf("Error returned during VerifyEnvironmentValues: %+v", err)
		} else if !correct || len(failtext) != 0 {
			t.Errorf("Success test #%d - For succeeding environment variables, VerifyEnvironmentValues returned a fail, correct=%v, failtext=%+v", i, correct, failtext)
		}
	}
}

func TestVerifyEnvironmentValuesFails(t *testing.T) {
	for i, c := range verifyEnvironmentValuesTestsFails {
		err, correct, failtext := VerifyEnvironmentValues(c.Input)
		if err != nil {
			t.Skipf("Error returned during VerifyEnvironmentValues: %+v", err)
		} else if correct || len(failtext) != c.Errors {
			t.Errorf("Fail test #%d - For failing environment variables, VerifyEnvironmentValues returned a success or wrong number of errors, correct=%v, failtext length=%d, failtext expected length=%d, failtext=%+v", i, correct, len(failtext), c.Errors, failtext)
		}
	}
}
