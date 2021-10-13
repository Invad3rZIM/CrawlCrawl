package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Simple Parsing Testing
func Test_parseBodyForURLs(t *testing.T) {
	bodyInput := "'https:google.com' and 'http://sometestinput.com' and \"https://somebrokenlink .com\" "
	assert.Equal(t, parseBodyForURLs(bodyInput), []string{"https:google.com", "http://sometestinput.com"})
}

///Simple Get Request Testing
func Test_getRequest_ExpectedPass(t *testing.T) {
	urlInput := "http://www.rescale.com"
	_, err := getRequest(urlInput)

	assert.Equal(t, err, nil)
}

func Test_getRequest_ExpectedFailure(t *testing.T) {
	urlInput := "http://www.rescale.com   "
	_, err := getRequest(urlInput)

	assert.NotEqual(t, err, nil)
}
