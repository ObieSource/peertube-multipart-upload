package main

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadEnvironmentVars(t *testing.T) {
	var input MultipartUploadHandlerHandlerInput
	var err error
	var failtext []string

	descFileText := "description file text"
	descfile, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatalf("Error during os.CreateTemp: %+v\n", err)
	}
	defer os.Remove(descfile.Name())
	if _, err = descfile.WriteString(descFileText); err != nil {
		t.Fatalf("Error during writing to description file: %+v\n", err)
	}
	descfile.Close()

	suppFileText := "support file text"
	suppfile, err := os.CreateTemp("", "*")
	if err != nil {
		t.Fatalf("Error during os.CreateTemp: %+v\n", err)
	}
	defer os.Remove(suppfile.Name())
	if _, err = suppfile.WriteString(suppFileText); err != nil {
		t.Fatalf("Error during writing to description file: %+v\n", err)
	}
	suppfile.Close()

	os.Clearenv()
	defer os.Clearenv()

	// test with no environment variables

	input, err, failtext = ReadEnvironmentVars()
	if !errors.Is(err, ReadEnvVarsFailed) {
		t.Fatalf("ReadEnvironmentVars on empty environment set should have returned ReadEnvVarsFailed. Returned %+v\n", err)
	}
	failtextLength := 10
	if len(failtext) != failtextLength {
		t.Fatalf("failtext length was not correct size for empty environment. Wanted %d got %d.%+v\n", failtextLength, len(failtext), failtext)
	}
	t.Log(input)

	if err = os.Setenv("PTHOST", "hostname"); err != nil {
		t.Fatalf("os.Setenv returned error \"%+v\"\n", err)
	}

	input, err, failtext = ReadEnvironmentVars()
	if !errors.Is(err, ReadEnvVarsFailed) {
		t.Fatalf("ReadEnvironmentVars on empty environment set should have returned ReadEnvVarsFailed. Returned %+v\n", err)
	}
	if len(failtext) != failtextLength-1 {
		t.Fatalf("failtext length was not correct size for only one environment variable defined. Wanted %d got %d.\n", failtextLength, len(failtext))
	}

	stuff := map[string]string{
		"PTHOST":      "https://peertube.cpy.re",
		"PTUSER":      "user",
		"PTPASS":      "passwordhere",
		"PTFILE":      descfile.Name(), // hacky solution to get working file pointer for this field
		"PTTITLE":     "title",
		"PTCAT":       "1",
		"PTCHAN":      "2",
		"PTTAGS":      "abcde,fghij,klmno",
		"PTDESCFILE":  descfile.Name(),
		"PTSUPP":      suppfile.Name(),
		"PTLANG":      "en",
		"PTCOMMENTS":  "false",
		"PTDOWNLOADS": "false",
		"PTNSFW":      "true",
		"PTPRIV":      "1",
		"PTLIC":       "1",
		"PTTYPE":      "application/text",
	}

	for key, val := range stuff {
		if err = os.Setenv(key, val); err != nil {
			t.Fatalf("os.Setenv error: \"%+v\"\n", err)
		}
	}

	input, err, failtext = ReadEnvironmentVars()
	if err != nil || len(failtext) != 0 {
		t.Fatalf("On ReadEnvironmentVars with all env variables, should have returend no error but returned %+v\n (Failtest %+v)", err, (failtext))
	}

	inputResult := MultipartUploadHandlerHandlerInput{
		Hostname:        "https://peertube.cpy.re",
		Username:        "user",
		Password:        "passwordhere",
		ContentType:     "application/text",
		ChannelID:       2,
		File:            nil,
		FileName:        input.FileName,
		DisplayName:     "title",
		Privacy:         1,
		Category:        1,
		CommentsEnabled: false,
		DescriptionText: input.DescriptionText,
		DownloadEnabled: false,
		Language:        "en",
		Licence:         1,
		NSFW:            true,
		SupportText:     input.SupportText,
		Tags: []string{
			"abcde",
			"fghij",
			"klmno",
		},
	}

	input.File = nil // cmp.Equal doesn't compare file values

	if !cmp.Equal(input, inputResult) {
		t.Log(cmp.Diff(input, inputResult))
		t.Error("ReadEnvironmentVars did not return correct variables")
	}

	/*
		Different file pointers lead to
		deep equal to be false
	*/
	input.File = nil
	for _, v := range trueStrs {
		stuff["PTNSFW"] = v
		input2, err, failtext := ReadEnvironmentVars()
		input2.File = nil
		if err != nil || len(failtext) != 0 || !reflect.DeepEqual(input2, input) {
			t.Fatalf("For using value %s for true, did not work or result in the same input. Error \"%+v\" len(failtext) == %d\n", v, err, len(failtext))
		}
	}

	for _, v := range falseStrs {
		stuff["PTCOMMENTS"] = v
		input2, err, failtext := ReadEnvironmentVars()
		input2.File = nil
		if err != nil || len(failtext) != 0 || !reflect.DeepEqual(input2, input) {
			t.Fatalf("For using value %s for false, did not work or result in the same input. Error \"%+v\" len(failtext) == %d\n", v, err, len(failtext))
		}
	}

}
