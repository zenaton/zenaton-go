package client

import (
	"runtime"

	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

func SetEnv(envFile string) {

	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		panic(thisFilePath)
	}

	thisFilePath = strings.Replace(thisFilePath, "/client.go", "/"+envFile, -1)
	variables, err := godotenv.Read(thisFilePath)
	if err != nil {
		panic("Error loading .env file " + err.Error())
	}

	godotenv.Load(thisFilePath)
	if err != nil {
		panic(err)
	}

	//make sure that all required environment variables are present
	appID, ok := variables["ZENATON_APP_ID"]
	if !ok {
		panic("Please add ZENATON_APP_ID env variable (https://zenaton.com/app/api)")
	}

	apiToken, ok := variables["ZENATON_API_TOKEN"]
	if !ok {
		panic("Please add ZENATON_API_TOKEN env variable (https://zenaton.com/app/api)")
	}

	appEnv, ok := variables["ZENATON_APP_ENV"]
	if !ok {
		panic("Please add ZENATON_APP_ENV env variable(https://zenaton.com/app/api)")
	}

	api := os.Getenv("ZENATON_API_URL")
	if api == "" {
		panic("Please add ZENATON_APP_ENV env variable(https://zenaton.com/app/api)")
	}

	zenaton.InitClient(appID, apiToken, appEnv)
}
