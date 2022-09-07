package restconf

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	agent := newRestconfAgent()
	defer agent.Close()

	expectedClient := &Client{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		restconfURL: agent.URL,
		username:    "username",
		password:    "password",
	}
	actualClient := New(30*time.Second, agent.URL, expectedClient.username, expectedClient.password)

	require.Equal(t, expectedClient, actualClient)
}

func TestNewLoopbackInterface(t *testing.T) {
	testCase := struct {
		testName   string
		httpClient *http.Client
		username   string
		password   string
		config     IetfInterfaceRequest
	}{
		testName: "create a loopback interface",
		username: "username",
		password: "password",
		config: IetfInterfaceRequest{IetfInterface: IetfInterface{
			Name:    "loopback200",
			Type:    IetfInterfaceType,
			Enabled: true,
			IPv4: Ipv4{
				Address: []Address{{
					Ip:      "20.0.0.1",
					Netmask: "255.255.255.255",
				}}}}},
	}

	agent := newRestconfAgent()
	defer agent.Close()

	client := New(30*time.Second, agent.URL, testCase.username, testCase.password)
	response, err := client.NewLoopbackInterface(context.Background(), testCase.config)

	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, http.StatusOK, response.StatusCode)
}

func TestFormInterfaceConfig(t *testing.T) {
	config := &IetfInterfaceRequest{IetfInterface: IetfInterface{
		Name:    "Loopback200",
		Type:    IetfInterfaceType,
		Enabled: true,
		IPv4: Ipv4{
			Address: []Address{{
				Ip:      "20.0.0.1",
				Netmask: "255.255.255.255",
			}}}}}

	actualConfig, err := json.Marshal(config)
	require.NoError(t, err)

	expectedConfig, err := ioutil.ReadFile("./test-resources/create-loopback-interface-yang-data.json")
	require.NoError(t, err)

	require.JSONEq(t, string(expectedConfig), string(actualConfig))
}

func newRestconfAgent() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}
