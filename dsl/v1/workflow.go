package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/expr-lang/expr"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"sync"
	"temporal-workflow/dsl/v1/activities"
	"temporal-workflow/dsl/v1/dslDefinition"
	"time"
)

// WorkflowState 工作流运行状态
type WorkflowState struct {
	CurrentNodeIds     []string            // 当前运行节点id
	CompletedNodeIds   map[string]struct{} // 运行完成的节点id
	NotExecutedNodeIds map[string]struct{} // 还未执行的节点id
	SkippedNodeIds     map[string]struct{} // 跳过执行的节点id
}

// SimpleDSLWorkflow 解析DSL中的node并进行运行
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow dslDefinition.Workflow, start string) ([]byte, error) {
	// 初始化状态
	state := WorkflowState{
		CurrentNodeIds:     []string{},
		CompletedNodeIds:   make(map[string]struct{}),
		NotExecutedNodeIds: make(map[string]struct{}),
		SkippedNodeIds:     make(map[string]struct{}),
	}

	// 把所有节点标记为未运行
	for _, node := range dslWorkflow.Nodes {
		state.NotExecutedNodeIds[node.Id] = struct{}{}
	}

	// 把开始节点标记为当前运行的节点
	//state.CurrentNodeIds = append(state.CurrentNodeIds, start)

	// 注册一个查询方法
	err := workflow.SetQueryHandler(ctx, "current_state", func() (WorkflowState, error) {
		return state, nil
	})
	if err != nil {
		fmt.Println("SetQueryHandler", err)
		return nil, err
	}

	shareVariables := sync.Map{}
	shareVariables.Store("variables", dslWorkflow.Variables)

	// 对拓扑进行排序
	sortedNodes, err := dslDefinition.TopSort(dslWorkflow.Nodes, dslWorkflow.Edges)
	if err != nil {
		return nil, err
	}
	// 生成执行组
	executionGroups := dslDefinition.DetectParallelExecution(sortedNodes, dslWorkflow.Edges)

	// 获取顶级节点（入度为0的节点）
	topLevelNodes := getTopLevelNodes(dslWorkflow.Nodes, dslWorkflow.Edges)
	if len(topLevelNodes) == 0 {
		return nil, fmt.Errorf("no top-level nodes found")
	}

	// 确认start节点是顶级节点
	if _, ok := topLevelNodes[start]; !ok {
		return nil, fmt.Errorf("start node %s is not a top-level node", start)
	}
	// 获取日志记录器
	logger := workflow.GetLogger(ctx)

	//// 跟踪已经访问过的节点
	//visited := make(map[string]bool)
	//// 标记

	// 获取从start节点开始的子图
	markedNodes := dslDefinition.MarkNodesFromStart(dslWorkflow.Edges, start)

	// 按照顺序执行每个执行组
	for _, group := range executionGroups {
		// 设置一个可取消的子ctx，用于并行执行的activity
		childCtx, cancelHandler := workflow.WithCancel(ctx)
		selector := workflow.NewSelector(ctx)
		var activityErr error

		var futures []workflow.Future

		// 遍历组中所有节点，一起启动
		for _, nodeId := range group {
			// 跳过不在子图中的节点
			if !markedNodes[nodeId] {
				state.SkippedNodeIds[nodeId] = struct{}{}
				continue
			}
			// 如果节点存在于skippedNodeIds里面，则跳过
			if _, ok := state.SkippedNodeIds[nodeId]; ok {
				continue
			}

			// 更新运行状态
			state.CurrentNodeIds = append(state.CurrentNodeIds, nodeId)
			delete(state.NotExecutedNodeIds, nodeId)

			// 根据id获取node对象
			node := getNodeByID(dslWorkflow.Nodes, nodeId)

			// 不同的activity类型有不同的处理方式，wait不方便在activity中运行
			switch node.Activity {
			case "Wait":
				var params dslDefinition.WaitParams
				paramBytes, err := json.Marshal(node.ActivityParameters)
				if err != nil {
					logger.Error("Failed to marshal activity parameters", "error", err)
					return nil, err
				}
				err = json.Unmarshal(paramBytes, &params)
				if err != nil {
					logger.Error("Failed to unmarshal activity parameters", "error", err)
					return nil, err
				}

				switch params.Resume {
				case dslDefinition.Interval:
					interval := time.Duration(params.IntervalMillis) * time.Millisecond
					logger.Info("Sleeping for interval", "interval", interval)
					err := workflow.Sleep(childCtx, interval)
					if err != nil {
						logger.Error("Sleep failed", "error", err)
						return nil, err
					}
				case dslDefinition.FixedTime:
					//fmt.Println(params.FixedTime)
					fixedTime, err := time.Parse(time.DateTime, params.FixedTime)
					if err != nil {
						logger.Error("Invalid fixed time format", "error", err)
						return nil, err
					}
					duration := time.Until(fixedTime)
					if duration > 0 {
						logger.Info("Sleeping until fixed time", "fixedTime", fixedTime)
						err = workflow.Sleep(childCtx, duration)
						if err != nil {
							logger.Error("Sleep failed", "error", err)
							return nil, err
						}
					}
				case dslDefinition.WebhookSignal:
					signalChannel := workflow.GetSignalChannel(childCtx, params.WebhookSignal)
					logger.Info("Waiting for webhook signal", "signal", params.WebhookSignal)
					signalChannel.Receive(childCtx, nil)
					logger.Info("Received webhook signal, proceeding with workflow")
				default:
					logger.Error("Unknown wait condition type", "type", params.Resume)
					return nil, fmt.Errorf("unknown wait condition type: %s", params.Resume)
				}
				// 更新节点状态，wait类型中webhook可能需要把请求信息存储到变量里面
				state.CompletedNodeIds[nodeId] = struct{}{}
				state.CurrentNodeIds = removeNode(state.CurrentNodeIds, nodeId)
			case "Switch":
				// 处理Switch节点
				// 根据表达式不同的结果，跳转到对应的nodeId
				var params dslDefinition.SwitchParams
				paramBytes, err := json.Marshal(node.ActivityParameters)
				if err != nil {
					// todo: 丰富日志细节
					logger.Error("Failed to marshal activity parameters", "error", err)
					return nil, err
				}
				err = json.Unmarshal(paramBytes, &params)
				if err != nil {
					logger.Error("Failed to unmarshal activity parameters", "error", err)
					return nil, err
				}
				// todo: 应该检查case里面的value是否都是本次workflow的NodeId
				nextNodeId, skipNodeIds, err := handleSwitchNode(params.Case, &shareVariables)
				if err != nil {
					logger.Error("Failed to handle switch node", "error", err)
					return nil, err
				}
				if nextNodeId != "" {
					// 到这里就是如何让nextNodeId去运行的问题了
					// 把当前表达式里面不满足条件的放在skippedNodeIds里面，这样下次到了不满足条件的节点就会跳过了

					for _, v := range skipNodeIds {
						state.SkippedNodeIds[v] = struct{}{}
					}
				}
				logger.Info(fmt.Sprintf("Switched node: %s, SkipedNodeIds: %v", nextNodeId, state.SkippedNodeIds))
				// 结束，把switch这个节点状态更新
				state.CompletedNodeIds[nodeId] = struct{}{}
				state.CurrentNodeIds = removeNode(state.CurrentNodeIds, nextNodeId)

			default:

				// 根据node里面的activity和activityVersion描述，找到对应的activity和activity的参数
				activityName := fmt.Sprintf("%sV%s", node.Activity, node.ActivityVersion)
				// 根据activityName，找到参数结构并填充实例
				params, err := activities.GetActivityParams(activityName, node.ActivityParameters)
				if err != nil {
					return nil, err
				}

				retryPolicy := &temporal.RetryPolicy{}
				// 如果配置了重试规则
				if node.RetryOnFail {
					retryPolicy = &temporal.RetryPolicy{
						InitialInterval:    time.Duration(node.WaitBetweenTries) * time.Millisecond, // 重试间隔
						BackoffCoefficient: 1.0,                                                     // 重试间隔系数
						MaximumAttempts:    int32(node.MaxTries),                                    // 最大重试次数
					}
				} else {
					// 设置默认的重试策略为1，即不重试
					retryPolicy = &temporal.RetryPolicy{
						InitialInterval:    time.Second,
						BackoffCoefficient: 2.0,
						MaximumInterval:    time.Second * 100,
						MaximumAttempts:    1, // Set to 1 to disable retries
					}
				}

				// 根据node上的配置信息填充activityOptions，然后执行具体的activity
				ao := workflow.ActivityOptions{
					// 这里可以使用指定任务队列来指定运行activity到特定主机/主机组上
					StartToCloseTimeout: time.Minute * 5, // activity超时时间(必填)
					//StartToCloseTimeout: time.Second * node.timeout, // activity超时时间
					RetryPolicy: retryPolicy, // 重试策略
					ActivityID:  nodeId,      // 设置activityId为nodeId
				}

				childCtx = workflow.WithActivityOptions(childCtx, ao)

				future := executeAsync(childCtx, activityName, params)

				// 这里收集信息
				//fmt.Println(workflow.GetInfo(childCtx).WorkflowExecution.)

				selector.AddFuture(future, func(f workflow.Future) {
					// 完成activity的回调函数
					var result interface{}
					err2 := f.Get(ctx, &result)
					// 无论运行状态是正常还是异常，这个节点到这里都结束了，更新节点状态
					state.CompletedNodeIds[nodeId] = struct{}{}
					state.CurrentNodeIds = removeNode(state.CurrentNodeIds, nodeId)

					if err2 != nil {
						switch node.OnError {
						case dslDefinition.ContinueErrorOutput:
							shareVariables.Store(nodeId, map[string]interface{}{
								"result": result,
								"error":  extractActivityError(err2),
							})
							logger.Warn(fmt.Sprintf("Continuing with err on nodeId: %s, result: %v, err: %v", node.Id, result, err2))
						case dslDefinition.ContinueRegularOutput:
							shareVariables.Store(nodeId, map[string]interface{}{
								"result": extractActivityError(err2),
								"error":  extractActivityError(err2),
							})
							logger.Warn(fmt.Sprintf("Continuing with regular on nodeId: %s, result: %v, err: %v", node.Id, result, err2))
						case dslDefinition.StopWorkflow:
							logger.Error(fmt.Sprintf("Stop workflow execution on nodeId:%s, failed: %v", node.Id, err2))
							cancelHandler()
							activityErr = err2
							return
						default:
							activityErr = fmt.Errorf("unknown failure type: %v", node.OnError)
							return
						}
					} else {
						shareVariables.Store(nodeId, map[string]interface{}{
							"result": result,
							"error":  err2,
						})
					}
					//// 把完成的节点状态更新
					//state.CompletedNodeIds[nodeId] = struct{}{}
					//state.CurrentNodeIds = removeNode(state.CurrentNodeIds, nodeId)

					// 转换结果为json格式
					byteResult, err := json.Marshal(result)
					if err != nil {
						logger.Error(fmt.Sprintf("Marshal future failed: %v", err))
					}

					logger.Info(fmt.Sprintf("Execution Group %v %v", group, string(byteResult)))
				})

				//future := workflow.ExecuteActivity(ctx, activityName, params)
				futures = append(futures, future)
			}

		}

		// 等待所有并行activity完成
		//for i, future := range futures {
		//	var result interface{}
		//	// 完成后，要把每个节点的输出结果记录到共享变量里面，key为节点id，value为节点输出，供下一层节点作为入参使用
		//	if err := future.Get(ctx, &result); err != nil {
		//		logger.Error(fmt.Sprintf("Get future failed: %v , result: %v", err, result))
		//
		//		// 在这里配置activity失败时策略
		//		nodeId := group[i]
		//		node := getNodeByID(dslWorkflow.Nodes, nodeId)
		//
		//		switch node.OnError {
		//		case ContinueErrorOutput:
		//			logger.Warn(fmt.Sprintf("Continuing with err on nodeId: %s, result: %v, err: %v", node.Id, result, err))
		//		case ContinueRegularOutput:
		//			logger.Warn(fmt.Sprintf("Continuing with regular on nodeId: %s, result: %v, err: %v", node.Id, result, result))
		//		case StopWorkflow:
		//			logger.Error(fmt.Sprintf("Stop workflow execution on nodeId:%s, failed: %v", node.Id, err))
		//			return nil, err
		//		default:
		//			return nil, fmt.Errorf("unknown failure type: %v", node.OnError)
		//		}
		//	}
		//	// 转换结果为json格式
		//	byteResult, err := json.Marshal(result)
		//	if err != nil {
		//		logger.Error(fmt.Sprintf("Marshal future failed: %v", err))
		//	}
		//
		//	logger.Info(fmt.Sprintf("Execution Group %v %v", group, string(byteResult)))
		//}

		// 等待完成
		for i := 0; i < len(futures); i++ {
			// 等待每个分支完成
			selector.Select(ctx)
			// 如果有分支执行失败，就退出等待
			if activityErr != nil {

				// 最后打印一下共享变量
				// 转换结果为json格式
				// 将 sync.Map 转换为常规 map
				sharedMap := convertSyncMapToMap(&shareVariables)

				// 打印共享变量
				byteResult, _ := json.Marshal(sharedMap)
				//if err != nil {
				//	logger.Error(fmt.Sprintf("Marshal future failed: %v", err))
				//}
				logger.Info(fmt.Sprintf("DSL Workflow ERROR %v", string(byteResult)))

				return nil, activityErr
			}
		}

		if activityErr != nil {
			logger.Error(fmt.Sprintf("Execute Group: %v, has error: %v", group, activityErr))
		}
	}

	// 将 sync.Map 转换为常规 map
	logger.Info(fmt.Sprintf("completed, %v", &shareVariables))
	sharedMap := convertSyncMapToMap(&shareVariables)

	// 打印共享变量
	byteResult, _ := json.Marshal(sharedMap)
	logger.Info(fmt.Sprintf("DSL Workflow completed %v", string(byteResult)))

	return nil, nil
}

// 根据 ID 获取节点
func getNodeByID(nodes []dslDefinition.Node, id string) dslDefinition.Node {
	for _, node := range nodes {
		if node.Id == id {
			return node
		}
	}
	return dslDefinition.Node{}
}

// getTopLevelNodes 获取顶级节点（没有target指向的节点）
func getTopLevelNodes(nodes []dslDefinition.Node, edges []dslDefinition.Edge) map[string]dslDefinition.Node {
	inDegree := make(map[string]int)
	topLevelNodes := make(map[string]dslDefinition.Node)

	// 初始化入度
	for _, node := range nodes {
		inDegree[node.Id] = 0
	}

	// 计算每个节点的入度
	for _, edge := range edges {
		inDegree[edge.Target]++
	}

	// 找出入度为0的节点
	for _, node := range nodes {
		if inDegree[node.Id] == 0 {
			topLevelNodes[node.Id] = node
		}
	}

	return topLevelNodes
}

// executeAsync
//
//	@Description: 异步执行
//	@Author cplinux98 2024-06-18 15:57:54
//	@param ctx
//	@param activity
//	@param params
//	@return workflow.Future
func executeAsync(ctx workflow.Context, activity interface{}, params ...interface{}) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		var result interface{}
		err := workflow.ExecuteActivity(ctx, activity, params...).Get(ctx, &result)
		settable.Set(result, err)
	})
	return future
}

// extractActivityError
//
//	@Description: 从ActivityError信息中读取原始错误
//	@Author cplinux98 2024-06-18 17:20:40
//	@param err
//	@return error
func extractActivityError(err error) string {
	// Extract the original error message from the Temporal error
	var temporalErr *temporal.ActivityError
	if errors.As(err, &temporalErr) {
		return fmt.Sprintf("%s", temporalErr.Unwrap())
	}
	return fmt.Sprintf("%s", err)
}

func convertSyncMapToMap(syncMap *sync.Map) map[string]interface{} {
	regularMap := make(map[string]interface{})
	syncMap.Range(func(key, value interface{}) bool {
		regularMap[key.(string)] = value
		return true
	})
	return regularMap
}

// 删除 slice 中的某个元素
func removeNode(slice []string, nodeID string) []string {
	for i, id := range slice {
		if id == nodeID {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

// 根据规则进行判断，返回下一节点id和跳过的ids
func handleSwitchNode(cases map[string]string, shareVariables *sync.Map) (string, []string, error) {
	env := make(map[string]interface{})
	shareVariables.Range(func(key, value any) bool {
		env[key.(string)] = value
		return true
	})
	skipNodeIds := make([]string, 0)
	var returnNodeId string

	for exprStr, targetNodeId := range cases {
		program, err := expr.Compile(exprStr, expr.Env(env))
		if err != nil {
			return "", skipNodeIds, fmt.Errorf("failed to compile expression %s: %v", exprStr, err)
		}

		output, err := expr.Run(program, env)
		if err != nil {
			return "", skipNodeIds, fmt.Errorf("failed to run expression %s: %v", exprStr, err)
		}
		if result, ok := output.(bool); ok && result {
			returnNodeId = targetNodeId
		} else {
			skipNodeIds = append(skipNodeIds, targetNodeId)
		}
	}

	return returnNodeId, skipNodeIds, nil
}
