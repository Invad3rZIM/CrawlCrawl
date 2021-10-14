package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//Simple Parsing Testing
func Test_parseBodyForURLs(t *testing.T) {
	bodyInput := "<a href='https:google.com'> and <a  href= 'http://sometestinput.com'> and <a href= \"https://somebrokenlink .com\" >"
	assert.Equal(t, []string{"https:google.com", "http://sometestinput.com"}, parseBodyForURLs(bodyInput))
}

///Simple Get Request Testing
func Test_getRequest_ExpectedPass(t *testing.T) {
	urlInput := "http://www.rescale.com"
	_, err := getRequest(urlInput)

	assert.Equal(t, nil, err)
}

func Test_getRequest_ExpectedFailure(t *testing.T) {
	urlInput := "http://www.rescale.com   "
	_, err := getRequest(urlInput)

	assert.NotEqual(t, nil, err)
}
