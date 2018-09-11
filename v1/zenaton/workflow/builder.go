package workflow

import (
	"github.com/zenaton/zenaton-go/v1/zenaton/client"
)

var builderInstance *QueryBuilder

type QueryBuilder struct {
	workflowClass string
	id            string
	client        *client.Client
}

func NewBuilder(workflow *WorkflowType) *QueryBuilder {
	if builderInstance == nil {
		builderInstance = &QueryBuilder{
			client:        client.NewClient(false),
			workflowClass: workflow.name,
		}
	}
	return builderInstance
}

func (b *QueryBuilder) WhereID(id string) *QueryBuilder {
	b.id = id
	return b
}

func (b *QueryBuilder) Find() (*Workflow, error) {
	output, err := b.client.FindWorkflow(b.workflowClass, b.id)
	if err != nil {
		return nil, err
	}

	properties := output["data"]["properties"].(string)
	name := output["data"]["name"].(string)

	return NewWorkflowManager().GetWorkflow(name, properties), nil
}

func (b *QueryBuilder) Send(eventName string, eventData interface{}) {
	b.client.SendEvent(b.workflowClass, b.id, eventName, eventData)
}

/**
 * Kill a workflow instance
 */

func (b *QueryBuilder) Kill() (*QueryBuilder, error) {
	err := b.client.KillWorkflow(b.workflowClass, b.id)
	return b, err
}

/**
* Pause a workflow instance
 */

func (b *QueryBuilder) Pause() (*QueryBuilder, error) {
	err := b.client.PauseWorkflow(b.workflowClass, b.id)
	return b, err
}

/**
* Resume a workflow instance
 */

func (b *QueryBuilder) Resume() (*QueryBuilder, error) {
	err := b.client.ResumeWorkflow(b.workflowClass, b.id)
	return b, err
}
