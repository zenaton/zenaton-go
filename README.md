# Zenaton library for Go

This Zenaton library for Go lets you code and launch workflows using Zenaton platform. You can sign up for an account at [https://zenaton/com](http://zenaton.com)

**DISCLAIMER** The go library is currently in public beta. Please open an
issue if you find any bugs.

## Requirements

This library has been tested with Go 1.10 and later.

## Installation

Execute:

    $ go get github.com/zenaton/zenaton-go/...

## Usage

For more detailed examples, please check [Zenaton go examples](https://github.com/zenaton/examples-go).

### Client Initialization

You will need to export three environment variables: `ZENATON_APP_ID`, `ZENATON_API_TOKEN`, `ZENATON_APP_ENV`. You"ll find them [here](https://zenaton/app/api).

Then you can initialize your Zenaton client:
```go
import(
    "os"
	"github.com/zenaton/zenaton-go/v1/zenaton"
)

appID = os.Getenv["ZENATON_APP_ID"]
apiToken = os.Getenv["ZENATON_API_TOKEN"]
appEnv = os.Getenv["ZENATON_APP_ENV"]

zenaton.InitClient(appID, apiToken, appEnv)
```

### Writing Workflows and Tasks

Writing a workflow is as simple as:

```go
import "github.com/zenaton/zenaton-go/v1/zenaton/workflow"

var MyWorkflow = workflow.New("MyWorkflow",
	func() (interface{}, error) {
        // Your workflow Implementation
	})
```
Note that your workflow implementation should be idempotent. See [documentation](https://zenaton.com/app/documentation#workflow-basics-implementation).

Writing a task is as simple as:
```go
import "github.com/zenaton/zenaton-go/v1/zenaton/task"

var MyTask = task.New("MyTask",
	func() (interface{}, error) {
        // Your task Implementation
	})
```

### Launching a workflow

Once your Zenaton client is initialized, you can start a workflow with

```go
MyWorkflow.New().Dispatch()
```

### Worker Installation

Your workflow's tasks will be executed on your worker servers. Please install a Zenaton worker on it:

    $ curl https://install.zenaton.com | sh

that you can start and configure with

    $ zenaton listen --boot=path/to/boot.go

where `boot.go` is a file that will be included before each task execution - this file should import all workflows. See [example boot.go](https://github.com/zenaton/examples-go/blob/master/boot/boot.go).

## Documentation

Please see https://zenaton.com/documentation for complete documentation.

## Development

To run tests, you can use go -v ./...
## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/zenaton/zenaton-go. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

The library is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).
