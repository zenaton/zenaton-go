package log

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"
)

var PrintTask = task.NewCustom("PrintTask", &Print{})

type Print struct {
	Values string
}

func (p *Print) Handle() (interface{}, error) {
	return fmt.Println(p.Values)
}

func (p *Print) Init(str string) {
	p.Values = str
}

func Println(values ...interface{}) {
	str := spew.Sprint(values...)
	PrintTask.New(str).Execute()
}
