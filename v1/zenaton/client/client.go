package client

import (
	"fmt"
	"os"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"github.com/zenaton/zenaton-go/v1/zenaton/services/http"
)

const (
	ZENATON_API_URL     = "https://zenaton.com/api/v1"
	ZENATON_WORKER_URL  = "http://localhost"
	DEFAULT_WORKER_PORT = 4001
	WORKER_API_VERSION  = "v_newton"

	MAX_ID_SIZE = 256

	APP_ENV   = "app_env"
	APP_ID    = "app_id"
	API_TOKEN = "api_token"

	ATTR_ID        = "custom_id"
	ATTR_NAME      = "name"
	ATTR_CANONICAL = "canonical_name"
	ATTR_DATA      = "data"
	ATTR_PROG      = "programming_language"
	ATTR_MODE      = "mode"

	PROG = "Javascript"

	EVENT_INPUT = "event_input"
	EVENT_NAME  = "event_name"

	WORKFLOW_KILL  = "kill"
	WORKFLOW_PAUSE = "pause"
	WORKFLOW_RUN   = "run"
)

var (
	instance *Client
	appID    string
	apiToken string
	appEnv   string
)

type Client struct {
}

// todo: figure out what's going on with the singleton thing in javascript
func InitClient(appIDx, apiTokenx, appEnvx string) {
	appID = appIDx
	apiToken = apiTokenx
	appEnv = appEnvx
}

func New(worker bool) *Client {
	if instance != nil {
		if !worker && (appID == "" || apiToken == "" || appEnv == "") {
			//todo: produce error?
			fmt.Println("Please initialize your Zenaton instance with your credentials")
			// throw new ExternalZenatonError('Please initialize your Zenaton instance with your credentials')
		}
		return instance
	}
	return &Client{}
}

//todo: figure out how to handle errors
func (c *Client) StartWorkflow(flowName, flowCanonical string) []interface{} {
	//todo: fix this so that it actually uses the ID() function?
	//customID := flow.ID()
	spew.Dump("bob", instance)
	customID := ""
	if len(customID) >= MAX_ID_SIZE {
		//todo: handle this error better
		fmt.Println(`Provided id must not exceed ` + strconv.Itoa(MAX_ID_SIZE) + ` bytes`)
	}

	body := make(map[string]interface{})
	body[ATTR_PROG] = PROG
	//body[ATTR_CANONICAL] = flowCanonical
	body[ATTR_NAME] = flowName
	//todo: use serializer here as in js
	body[ATTR_DATA] = "{}"
	body[ATTR_ID] = customID

	http.Post(c.getInstanceWorkerUrl(""), body)

	return []interface{}{}
}

//todo: fill this out from js example
func (c *Client) SendEvent(workflowName, customID, eventName string, eventData interface{}) {
	fmt.Println(workflowName, customID, eventName, eventData)
}

func (c *Client) getInstanceWorkerUrl(params string) string {
	return c.getWorkerUrl("instances", params)
}

func (c *Client) getWorkerUrl(resources string, params string) string {
	workerURL := os.Getenv("ZENATON_WORKER_URL")
	if workerURL == "" {
		workerURL = ZENATON_WORKER_URL
	}

	workerPort := os.Getenv("ZENATON_WORKER_PORT")
	if workerPort == "" {
		workerPort = strconv.Itoa(DEFAULT_WORKER_PORT)
	}

	url := workerURL + ":" + workerPort + "/api/" + WORKER_API_VERSION +
		"/" + resources + "?"

	return c.addAppEnv(url, params)
}

func (c *Client) addAppEnv(url, params string) string {

	var appEnvx string
	if appEnv != "" {
		appEnvx = APP_ENV + "=" + appEnv + "&"
	}

	var appIDx string
	if appID != "" {
		appIDx = APP_ID + "=" + appID + "&"
	}

	if params != "" {
		params = params + "&"
	}

	fmt.Println("do we get here?", appEnv, appID, params)
	return url + appEnvx + appIDx + params
	return ""
}
