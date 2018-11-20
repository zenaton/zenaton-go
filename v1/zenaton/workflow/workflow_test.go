package workflow_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var _ = Describe("Workflow", func() {

	Context("New Definition", func() {
		Context("When it has the same name as an existing def", func() {
			It("should panic", func() {
				workflow.New("testRedundantNameWithNew", func() (interface{}, error) { return nil, nil })
				defer func() {
					r := recover()
					fmt.Println("r: ", r)
					Expect(r).To(Equal("workflow definition with name 'testRedundantNameWithNew' already exists"))
				}()
				workflow.New("testRedundantNameWithNew", func() (interface{}, error) { return nil, nil })
			})
		})
	})

	Context("NewCustom Definition", func() {

		Context("with the same name as an existing def", func() {
			It("should panic", func() {
				workflow.NewCustom("testRedundantNameWithNewCustom", &TestHandler{})
				defer func() {
					r := recover()
					fmt.Println("r: ", r)
					Expect(r).To(Equal("workflow definition with name 'testRedundantNameWithNewCustom' already exists"))
				}()
				workflow.NewCustom("testRedundantNameWithNewCustom", &TestHandler{})
			})
		})

		Context("without a pointer", func() {
			It("should panic", func() {
				defer func() {
					r := recover()
					fmt.Println("r: ", r)
					Expect(r).To(Equal("must pass a pointer to NewCustom"))
				}()
				workflow.NewCustom("testRedundantNameWithNewCustom", TestHandler{})
			})
		})

		Context("with an unserializable Handler type", func() {
			It("should panic", func() {

				defer func() {
					r := recover()
					fmt.Println("r: ", r)
					Expect(r).To(Equal("workflow: Handler type 'unserializableHandler' must be able to be marshaled to json. json: unsupported type: func()"))
				}()
				workflow.NewCustom("testRedundantNameWithNewCustom", &unserializableHandler{})
			})
		})
	})

	Context("New Instance", func() {

		Context("with unserializable input", func() {
			It("should panic", func() {

				defer func() {
					r := recover()
					Expect(r).To(Equal("workflow: must be able to json unmarshal into the handler type... json: cannot unmarshal object into Go struct field Unserializable.Error of type error"))
				}()
				UnserializableWorkflow.New(TestErrorType{Message: "test error message", Prefix: "test: "})
			})
		})
	})

	//Context("Given a workflow instance", func() {
	//	It("should be able to Execute the workflow", func() {
	//		testTask := task.New("testTask",
	//			func()(interface{}, error){
	//				return "test task output", nil
	//			})
	//		testWorkflow := workflow.New("parallelWorkflow",
	//			func() (interface{}, error) {
	//				var output string
	//				testTask.New().Execute().Output(&output)
	//				return output, nil
	//			})
	//
	//		var workflowOutput string
	//		testWorkflow.New().Execute().Output(&workflowOutput)
	//
	//
	//	})
	//})
})

type unserializableHandler struct{ Func func() }

func (u *unserializableHandler) Handle() (interface{}, error) { return nil, nil }

type TestHandler struct{}

func (th TestHandler) Handle() (interface{}, error) { return nil, nil }

var UnserializableWorkflow = workflow.NewCustom("UnserializableWorkflow", &Unserializable{})

type Unserializable struct {
	Error error
}

func (u *Unserializable) Init(err error) {
	u.Error = err
}

func (u *Unserializable) Handle() (interface{}, error) {
	return nil, u.Error
}

type TestErrorType struct {
	Message string
	Prefix  string
}

func (tet TestErrorType) Error() string {
	return tet.Prefix + tet.Message
}
