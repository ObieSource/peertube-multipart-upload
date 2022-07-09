package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
)

var (
	ReadEnvVarsFailed = errors.New("Read Environment Variables Failed")

	trueStrs = []string{
		"true",
		"t",
		"yes",
		"y",
		"1",
	}
	falseStrs = []string{
		"false",
		"f",
		"no",
		"n",
		"0",
	}
	VideoChunkSize VideoFileByteCounter = 1024 * 1024 * 1
)

func ReadEnvironmentVars() (input MultipartUploadHandlerHandlerInput, erro error, failtext []string) {
	var ok bool
	var fail bool
	var err error
	tagsraw := new(string)
	descfile := new(string)
	suppfile := new(string)

	var StringReqEnvVars map[string](*string) = map[string](*string){
		"PTHOST":  &input.Hostname,
		"PTUSER":  &input.Username,
		"PTFILE":  &input.FileName,
		"PTTITLE": &input.DisplayName,
		"PTTYPE":  &input.ContentType,
		"PTLANG":  &input.Language,
	}
	var StringEnvVars map[string](*string) = map[string](*string){
		"PTTAGS": tagsraw,
		"PTDESC": descfile,
		"PTSUPP": suppfile,
		"PTPASS": &input.Password,
	}
	for key, val := range StringReqEnvVars {
		if *val, ok = os.LookupEnv(key); !ok {
			failtext = append(failtext, fmt.Sprintf("ERROR: Environment variable %s is required.\n", key))
			fail = true
		}
	}
	for key, val := range StringEnvVars {
		*val = os.Getenv(key)
	}

	// clean hostname
	input.Hostname = CleanHostname(input.Hostname)

	// get the text for the description and support
	if *descfile != "" { // skip if not specified
		input.DescriptionText, err = GetDescriptionTextFromFilename(*descfile)
		if err != nil {
			failtext = append(failtext, fmt.Sprintf("Error while reading description file: %+v\n", err))
			fail = true
		}
	}

	if *suppfile != "" {
		input.SupportText, err = GetDescriptionTextFromFilename(*suppfile)
		if err != nil {
			failtext = append(failtext, fmt.Sprintf("Error while reading support text file: %+v\n", err))
			fail = true
		}
	}

	// set tags
	input.Tags, err, _ = GetTagsFromEnv(*tagsraw)
	if err != nil {
		failtext = append(failtext, fmt.Sprintf("Tags error: %+v\n", err))
		fail = true
	}

	// set video file
	input.File, err = GetVideoFileReader(input.FileName, VideoChunkSize)
	if err != nil {
		failtext = append(failtext, fmt.Sprintf("Error when reading video file: \"%+v\"", err))
		fail = true
	}

	// input.FileName should only be actual file name
	spl := strings.Split(input.FileName, "/")
	input.FileName = spl[len(spl)-1]

	IntReqEnvVars := map[string]*int{
		"PTCHAN": &input.ChannelID,
		"PTPRIV": &input.Privacy,
		"PTLIC":  &input.Licence,
	}
	IntEnvVars := map[string]*int{
		"PTCAT": &input.Category,
	}
	for key, val := range IntReqEnvVars {
		var err error
		env := os.Getenv(key)
		*val, err = strconv.Atoi(env)
		if err != nil {
			failtext = append(failtext, fmt.Sprintf("ERROR: Environment integer variable %s is required and/or was not an integer.\n", key))
			fail = true
		}
	}
	for key, val := range IntEnvVars {
		var err error
		env := os.Getenv(key)
		*val, err = strconv.Atoi(env)
		if err != nil && env != "" {
			failtext = append(failtext, fmt.Sprintf("ERROR: Environment variable %s must be an int (recieved %s)\n", key, env))
			fail = true
		} else if env == "" {
			*val = 0
		}
	}

	// default category is 10 (entertainment)
	if input.Category == 0 {
		input.Category = 10
	}
	if input.Licence == 0 {
		input.Licence = 1
	}

	BoolEnvVars := map[string]*bool{
		"PTCOMMENTS":  &input.CommentsEnabled,
		"PTDOWNLOADS": &input.DownloadEnabled,
		"PTNSFW":      &input.NSFW,
	}
	// hard code boolean values
	input.CommentsEnabled = true
	input.DownloadEnabled = true
	input.NSFW = false

	for key, val := range BoolEnvVars {
		// all unrequired
		env := os.Getenv(key)
		if in(env, trueStrs) {
			*val = true
		} else if in(env, falseStrs) {
			*val = false
		} else if env != "" {
			failtext = append(failtext, fmt.Sprintf("Unknown value for environment variable %s: \"%s\"\n", key, env))
		}
		/*
			Otherwise, skip and stay with
			default values
		*/
	}

	erro = nil
	if fail {
		erro = ReadEnvVarsFailed
		return
	}

	err, verified, verifyFailText := VerifyEnvironmentValues(input)
	if err != nil {
		erro = err
		return
	}

	if len(verifyFailText) != 0 || !verified {
		for _, line := range verifyFailText {
			failtext = append(failtext, line)
		}
		erro = ReadEnvVarsFailed
	}

	/*
		Note: password is set earlier in this function.
	*/
	/*
		Set password after this function is
		run so that errors relating to
		unset environment variables are
		shown to the user before entering
		the password.
	*/
	return
}

func PasswordSecret() (passwd string, err error) {
	/*
		Borrowed from https://stackoverflow.com/a/32768479
	*/
	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return
	}
	passwd = strings.TrimSpace(string(bytePassword))

	return
}

func in(substr string, all []string) bool {
	for _, s := range all {
		if s == substr {
			return true
		}
	}
	return false
}
