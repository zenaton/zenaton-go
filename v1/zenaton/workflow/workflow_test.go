package workflow_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenaton/zenaton-go/v1/zenaton/workflow"
)

var UnserializableWorkflow = workflow.NewCustom("UnserializableWorkflow", &Unserializable{})

type Unserializable struct {
	Error error
}

func (u *Unserializable) Init(err error) {
	u.Error = err
}

func (u *Unserializable) Handle() (interface{}, error){
	return nil, u.Error
}

type TestErrorType struct {
	Message string
	Prefix string
}

func (tet TestErrorType) Error() string{
	return tet.Prefix + tet.Message
}

var _ = Describe("Workflow", func() {
	Context("When creating a workflow with unserializable input", func (){
		It("should panic", func(){

			defer func(){
				r := recover()
				Expect(r).To(Equal("workflow: must be able to json unmarshal into the handler type... json: cannot unmarshal object into Go struct field Unserializable.Error of type error"))
			}()
			UnserializableWorkflow.New(TestErrorType{Message: "test error message", Prefix: "test: ",})
		})
	})
})
