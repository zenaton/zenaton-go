package recursive

import (
	"strconv"

	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

func NewRecursiveWorkflow(id, max int) *workflow.Workflow {

	return &workflow.Workflow{
		Name: "RecursiveWorkflow",

		HandleFunc: func() interface{} {

			for counter := 0; counter < 10; counter++ {
				NewDisplayTask(strconv.Itoa(counter)).Execute()
			}

			NewRelaunchTask(id, max).Execute()
			return nil
		},

		ID: func() string {
			return strconv.Itoa(id)
		},
	}
}
