package parallel_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenaton/zenaton-go/v1/zenaton/parallel"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
	"time"
)

var _ = Describe("Parallel", func() {
	Context("given a workflow and a task", func(){
		It("should be able to run a workflow and task in parallel", func(){
			var output string
			a := task.New("a", func() (interface{}, error) {
				time.Sleep(10 * time.Millisecond)
				output = output + "a"
				return nil, nil
			})

			b := task.New("b", func() (interface{}, error) {
				output = output + "b"
				return nil, nil
			})

			testWorkflow := workflow.New("testWorkflow",
				func() (interface{}, error) {
					b.New().Execute()
					return nil, nil
				})

			parallelWorkflow := workflow.New("parallelWorkflow",
				func()(interface{}, error) {
					parallel.Jobs{
						testWorkflow.New(),
						a.New(),
					}.Execute()
					return nil, nil
				})

			parallelWorkflow.New().Execute()
			Expect(output).To(Equal("ba"))
		})
	})
})
