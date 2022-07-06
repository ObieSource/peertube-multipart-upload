package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	log.Println(categoriesMap)

	return
}
