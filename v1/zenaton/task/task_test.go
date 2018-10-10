package task_test

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Task", func() {
	Context("When creating a task with unserializable input", func() {
		It("should panic", func() {

			defer func() {
				r := recover()
				fmt.Println("r: ", r)
				Expect(r).To(Equal("task: must be able to json unmarshal into the handler type... json: cannot unmarshal object into Go struct field Unserializable.Error of type error"))
			}()
			UnserializableTask.New(TestErrorType{"test error message"})
		})
	})
})

var UnserializableTask = task.NewCustom("UnserializableTask", &Unserializable{})

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
}

func (tet TestErrorType) Error() string {
	return tet.Message
}
