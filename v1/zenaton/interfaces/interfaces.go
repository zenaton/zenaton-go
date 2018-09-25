package interfaces

//todo: move this somewhere else
type Handler interface {
	Handle() (interface{}, error)
}

type Starter interface {
	Start()
}

type Successer interface {
	Success()
}

type IDer interface {
	ID() string
}

type OnEventer interface {
	OnEvent(string, interface{})
}

// todo: maybe I don't want this to be exported, so only tasks and workflows can implement this interface
