package dsl

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	Variable struct {
		Name         string // 变量名称
		Type         int64  // 变量类型
		Value        string // 变量输入的值
		DefaultValue string // 默认的变量值
		Helper       string // 变量输入时的提示信息
		Hidden       bool   // 是否隐藏
	}

	// Workflow is the type used to express the workflow definition. Variables are a map of valuables. Variables can be
	// used as input to Activity.
	Workflow struct {
		Id        string              // 流程id
		Name      string              // 流程名称
		Variables map[string]Variable // 流程变量
		Tasks     []*Task             // 任务节点列表
	}

	// Task 任务单元
	Task struct {
		Activity  string   // 具体执行的activity
		Name      string   // 任务名称
		Arguments []string // 任务参数
		Result    string   // 结果
		DependsOn []string // 依赖的前置任务
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]string) error
	}
)

// SimpleDSLWorkflow workflow definition
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	bindings := make(map[string]string)
	//workflowcheck:ignore Only iterates for building another map
	for k, v := range dslWorkflow.Variables {
		if v.Value != "" {
			bindings[k] = v.Value
		} else {
			bindings[k] = v.DefaultValue
		}
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	var results = make(map[string]string)

	for _, task := range dslWorkflow.Tasks {
		if len(task.DependsOn) > 0 {
			for _, dep := range task.DependsOn {
				workflow.Await(ctx, func() bool {
					_, exists := results[dep]
					return exists
				})
			}
		}

		var result string
		var err error

		inputParam := makeInput(task.Arguments, bindings)
		err = workflow.ExecuteActivity(ctx, task.Name, inputParam).Get(ctx, &result)
		if err != nil {
			logger.Error("failed to execute task", "task", task.Name, "error", err)
			return nil, err
		}
		if task.Result != "" {
			bindings[task.Result] = result
		}

		results[task.Result] = result

	}

	logger.Info("DSL Workflow completed.")
	return nil, nil
}

func (t Task) execute(ctx workflow.Context, bindings map[string]string) error {
	inputParam := makeInput(t.Arguments, bindings)
	var result string
	err := workflow.ExecuteActivity(ctx, t.Name, inputParam).Get(ctx, &result)
	if err != nil {
		return err
	}
	if t.Result != "" {
		bindings[t.Result] = result
	}
	return nil
}

//func (p Parallel) execute(ctx workflow.Context, bindings map[string]string) error {
//	//
//	// You can use the context passed in to activity as a way to cancel the activity like standard GO way.
//	// Cancelling a parent context will cancel all the derived contexts as well.
//	//
//
//	// In the parallel block, we want to execute all of them in parallel and wait for all of them.
//	// if one activity fails then we want to cancel all the rest of them as well.
//	childCtx, cancelHandler := workflow.WithCancel(ctx)
//	selector := workflow.NewSelector(ctx)
//	var activityErr error
//	for _, s := range p.Branches {
//		f := executeAsync(s, childCtx, bindings)
//		selector.AddFuture(f, func(f workflow.Future) {
//			err := f.Get(ctx, nil)
//			if err != nil {
//				// cancel all pending activities
//				cancelHandler()
//				activityErr = err
//			}
//		})
//	}
//
//	for i := 0; i < len(p.Branches); i++ {
//		selector.Select(ctx) // this will wait for one branch
//		if activityErr != nil {
//			return activityErr
//		}
//	}
//
//	return nil
//}

func executeAsync(exe executable, ctx workflow.Context, bindings map[string]string) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, bindings)
		settable.Set(nil, err)
	})
	return future
}

func makeInput(argNames []string, argsMap map[string]string) []string {
	var args []string
	for _, arg := range argNames {
		args = append(args, argsMap[arg])
	}
	return args
}
