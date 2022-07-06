package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func VerifyEnvironmentValues(input MultipartUploadHandlerHandlerInput) (err error, correct bool, failtext []string) {
	correct = true

	host := input.Hostname

	/*
		category, language, license, and privacy have
		more constrained inputs. Use additional API
		calls to verify that these are correct.
	*/

	// Category: /api/v1/videos/categories
	categoriesRequest, err := http.Get(fmt.Sprintf("%s/api/v1/videos/categories", host))
	if err != nil {
		return
	}

	defer categoriesRequest.Body.Close()
	categoriesBytes, err := io.ReadAll(categoriesRequest.Body)
	if err != nil {
		return
	}

	var categoriesMap map[string]string

	if err = json.Unmarshal(categoriesBytes, &categoriesMap); err != nil {
		return
	}

	// test that the category is one of the integers.
	categoryExists := false
	for key, _ := range categoriesMap {
		if key == strconv.Itoa(input.Category) {
			categoryExists = true
			break
		}
	}
	if !categoryExists {
		failtext = append(failtext, fmt.Sprintf("Specied category \"%d\" is not valid. See %s/api/v1/videos/categories", input.Category, host))
		correct = false
	}

	// Language: /api/v1/videos/languages
	languagesRequest, err := http.Get(fmt.Sprintf("%s/api/v1/videos/languages", host))
	if err != nil {
		return
	}

	defer languagesRequest.Body.Close()

	languagesBytes, err := io.ReadAll(languagesRequest.Body)
	if err != nil {
		return
	}

	var languagesMap map[string]string

	if err = json.Unmarshal(languagesBytes, &languagesMap); err != nil {
		return
	}

	languageExists := false

	for key, _ := range languagesMap {
		if key == input.Language {
			languageExists = true
			break
		}
	}

	if !languageExists {
		failtext = append(failtext, fmt.Sprintf("Specified languages \"%s\" is not valid. See %s/api/v1/videos/categories", input.Language, host))
		correct = false
	}

	// licenses: /api/v1/videos/licences
	licensesRequest, err := http.Get(fmt.Sprintf("%s/api/v1/videos/licences", host))
	if err != nil {
		return
	}

	defer licensesRequest.Body.Close()

	licensesBytes, err := io.ReadAll(licensesRequest.Body)
	if err != nil {
		return
	}

	var licensesMap map[string]string

	if err = json.Unmarshal(licensesBytes, &licensesMap); err != nil {
		return
	}

	licenseExists := false
	for key, _ := range licensesMap {
		if key == strconv.Itoa(input.Licence) {
			licenseExists = true
			break
		}
	}

	if !licenseExists {
		failtext = append(failtext, fmt.Sprintf("Specified licence \"%d\" is not valid. See %s/api/v1/videos/licences", input.Licence, host))
		correct = false
	}

	// privacy: /api/v1/videos/privacies
	privacyRequest, err := http.Get(fmt.Sprintf("%s/api/v1/videos/privacies", host))
	if err != nil {
		return
	}

	defer privacyRequest.Body.Close()

	privacyBytes, err := io.ReadAll(privacyRequest.Body)
	if err != nil {
		return
	}

	var privacyMap map[string]string

	if err = json.Unmarshal(privacyBytes, &privacyMap); err != nil {
		return
	}

	privacyExists := false
	for key, _ := range privacyMap {
		if key == strconv.Itoa(input.Privacy) {
			privacyExists = true
			break
		}
	}

	if !privacyExists {
		failtext = append(failtext, fmt.Sprintf("Specified privacy \"%d\" is not valid. See %s/api/v1/videos/privacies", input.Licence, host))
		correct = false
	}

	return
}
