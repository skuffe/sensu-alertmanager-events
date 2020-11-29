package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	plugin.AgentAPIURL = "http://127.0.0.1:3031/events"
	event := v2.FixtureEvent("entity1", "check1")
	status, err := checkArgs(event)
	assert.NoError(err)
	assert.Equal(sensu.CheckStateOK, status)
	status, err = checkArgs(event)
	assert.NoError(err)
	assert.Equal(sensu.CheckStateOK, status)
	plugin.AlertmanagerLabelEntity = "cluster"
	plugin.SensuProxyEntity = "k8s-cluster"
	status, err = checkArgs(event)
	assert.Error(err)
	assert.Equal(sensu.CheckStateWarning, status)
}

func TestSubmitEventAgentAPI(t *testing.T) {
	testcases := []struct {
		httpStatus  int
		expectError bool
	}{
		{http.StatusOK, false},
		{http.StatusBadRequest, true},
	}
	for _, tc := range testcases {
		assert := assert.New(t)
		event := v2.FixtureEvent("entity1", "check1")
		var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(err)
			eV := &v2.Event{}
			err = json.Unmarshal(body, eV)
			w.WriteHeader(tc.httpStatus)
		}))
		_, err := url.ParseRequestURI(test.URL)
		plugin.AgentAPIURL = test.URL
		err = submitEventAgentAPI(event)
		if tc.expectError {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestStringInSlice(t *testing.T) {
	testSlice := []string{"foo", "bar", "test"}
	testString := "test"
	testResult := stringInSlice(testString, testSlice)
	assert.True(t, testResult)
}
