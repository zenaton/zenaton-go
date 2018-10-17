package task_test

import (
	"fmt"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
	"time"

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

	FContext("when dispatching a task locally", func() {
		It("it should actually run the task asynchronously", func() {

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

			a.New().Dispatch()
			b.New().Dispatch()

			time.Sleep(20 * time.Millisecond)

			Expect(output).To(Equal("ba"))
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
