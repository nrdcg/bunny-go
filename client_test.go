package bunny

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckRespWithEmptyUnsuccessfulResp(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://test.de", nil)
	require.NoError(t, err)

	resp := http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	clt := NewClient("")

	err = clt.checkResp(req, &resp)
	require.Error(t, err)
	require.IsType(t, &HTTPError{}, err)

	var httpErr *HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Empty(t, httpErr.Errors)
}

func TestCheckRespWithJSONBody(t *testing.T) {
	apiErr := APIError{
		ErrorKey: "err",
		Field:    "id",
		Message:  "something br0ke",
	}

	buf, err := json.Marshal(&apiErr)
	require.NoError(t, err)

	const reqURL = "http://test.de"
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	require.NoError(t, err)

	hdr := http.Header{}
	hdr.Add("Content-Type", "application/json; charset=utf-8")

	resp := http.Response{
		Header:     hdr,
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader(buf)),
	}

	clt := NewClient("")

	err = clt.checkResp(req, &resp)
	require.Error(t, err)
	require.IsType(t, &APIError{}, err, "error: "+err.Error())

	var retAPIErr *APIError
	assert.ErrorAs(t, err, &retAPIErr)
	assert.Equal(t, apiErr.ErrorKey, retAPIErr.ErrorKey, "unexpected errorKey value")
	assert.Equal(t, apiErr.Field, retAPIErr.Field, "unexpected field value")
	assert.Equal(t, apiErr.Message, retAPIErr.Message, "unexpected message value")

	assert.Equal(t, reqURL, retAPIErr.RequestURL, "unexpected RequestURL")
	assert.Equal(t, resp.StatusCode, retAPIErr.StatusCode, "unexpected status code")
	assert.Equal(t, buf, retAPIErr.RespBody)
}

func TestCheckRespWithJSONBodyAndMissingContentType(t *testing.T) {
	buf, err := json.Marshal(&APIError{Message: "something br0ke"})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)

	resp := http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewReader(buf)),
	}

	clt := NewClient("")

	err = clt.checkResp(req, &resp)
	require.Error(t, err)
	require.IsType(t, &HTTPError{}, err, "error: "+err.Error())

	var retErr *HTTPError
	assert.ErrorAs(t, err, &retErr)
	assert.Equal(t, buf, retErr.RespBody)

	assert.EqualError(t, retErr.Errors[0], "processing response failed: Content-Type header is missing or empty")
}

func TestUnmarshalHTTPJSONBody(t *testing.T) {
	hostnameval := "hello"
	msgIn := Hostname{
		Value: &hostnameval,
	}
	buf, err := json.Marshal(&msgIn)
	require.NoError(t, err)

	hdr := http.Header{}
	hdr.Add("Content-Type", "application/json; charset=utf-8")
	resp := http.Response{
		Body:   io.NopCloser(bytes.NewReader(buf)),
		Header: hdr,
	}

	clt := NewClient("")

	var msgOut Hostname

	err = clt.unmarshalHTTPJSONBody(&resp, "", &msgOut)
	require.NoError(t, err)

	require.NotNil(t, msgOut.Value)
	require.Equal(t, *msgIn.Value, *msgOut.Value)
}

func TestUnmarshalHTTPJSONBodyWithMissingContentType(t *testing.T) {
	msgIn := Hostname{}
	buf, err := json.Marshal(&msgIn)
	require.NoError(t, err)

	code := 200
	resp := http.Response{
		StatusCode: code,
		Body:       io.NopCloser(bytes.NewReader(buf)),
	}

	clt := NewClient("")

	var msgOut Hostname

	url := "http://test.de"
	err = clt.unmarshalHTTPJSONBody(&resp, url, &msgOut)
	require.Error(t, err)

	require.IsType(t, &HTTPError{}, err)

	var httpErr *HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, httpErr.RequestURL, url)
	assert.Equal(t, httpErr.StatusCode, code)
	assert.Len(t, httpErr.Errors, 1)
	assert.EqualError(t, httpErr.Errors[0], "processing response failed: Content-Type header is missing or empty")
	assert.Equal(t, buf, httpErr.RespBody)
}

func TestUnmarshalHTTPJSONBodyWithWrongContentType(t *testing.T) {
	msgIn := Hostname{}
	buf, err := json.Marshal(&msgIn)
	require.NoError(t, err)

	hdr := http.Header{}
	hdr.Add("Content-Type", "application/binary")

	code := 200
	resp := http.Response{
		StatusCode: code,
		Header:     hdr,
		Body:       io.NopCloser(bytes.NewReader(buf)),
	}

	clt := NewClient("")

	var msgOut Hostname

	url := "http://test.de"
	err = clt.unmarshalHTTPJSONBody(&resp, url, &msgOut)
	require.Error(t, err)

	require.IsType(t, &HTTPError{}, err)

	var httpErr *HTTPError
	assert.ErrorAs(t, err, &httpErr)
	assert.Equal(t, httpErr.RequestURL, url)
	assert.Equal(t, httpErr.StatusCode, code)
	assert.Equal(t, buf, httpErr.RespBody)
	assert.Len(t, httpErr.Errors, 1)
	assert.EqualError(t, httpErr.Errors[0], "processing response failed: expected Content-Type to be \"application/json\", got: \"application/binary\"")
}
