package cadence

import (
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.uber.org/cadence/workflow"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	return interpreter.InterpreterImpl(interpreter.NewUnifiedContext(ctx), newCadenceWorkflowProvider(), input)
}

func WaitforStateCompletionWorkflow(ctx workflow.Context) (*service.WaitForStateCompletionWorkflowOutput, error) {
	return interpreter.WaitForStateCompletionWorkflowImpl(interpreter.NewUnifiedContext(ctx), newCadenceWorkflowProvider())
}
