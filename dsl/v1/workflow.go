package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"sync"
	"temporal-workflow/dsl/v1/activities"
	"time"
)

// SimpleDSLWorkflow 解析DSL中的node并进行运行
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow, start string) ([]byte, error) {
	// 先把初始的共享变量存储起来
	//shareVariables := make(map[string]interface{})
	//shareVariables["variables"] = dslWorkflow.Variables
	//workflowcheck:ignore Only iterates for building another map
	//for k, v := range dslWorkflow.Variables {
	//	shareVariables[k] = v
	//}

	shareVariables := sync.Map{}
	shareVariables.Store("variables", dslWorkflow.Variables)

	// 对拓扑进行排序
	sortedNodes, err := TopSort(dslWorkflow.Nodes, dslWorkflow.Edges)
	if err != nil {
		return nil, err
	}
	// 生成执行组
	executionGroups := DetectParallelExecution(sortedNodes, dslWorkflow.Edges)

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
	markedNodes := MarkNodesFromStart(dslWorkflow.Edges, start)

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
				continue
			}
			// 根据id获取node对象
			node := getNodeByID(dslWorkflow.Nodes, nodeId)

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
				StartToCloseTimeout: time.Minute * 5, // activity超时时间
				//StartToCloseTimeout: time.Second * node.timeout, // activity超时时间
				RetryPolicy: retryPolicy,
			}

			childCtx = workflow.WithActivityOptions(childCtx, ao)

			future := executeAsync(childCtx, activityName, params)
			selector.AddFuture(future, func(f workflow.Future) {
				// 完成activity的回调函数
				var result interface{}
				err2 := f.Get(ctx, &result)
				if err2 != nil {
					switch node.OnError {
					case ContinueErrorOutput:
						shareVariables.Store(nodeId, map[string]interface{}{
							"result": result,
							"error":  extractActivityError(err2),
						})
						logger.Warn(fmt.Sprintf("Continuing with err on nodeId: %s, result: %v, err: %v", node.Id, result, err2))
					case ContinueRegularOutput:
						shareVariables.Store(nodeId, map[string]interface{}{
							"result": extractActivityError(err2),
							"error":  extractActivityError(err2),
						})
						logger.Warn(fmt.Sprintf("Continuing with regular on nodeId: %s, result: %v, err: %v", node.Id, result, err2))
					case StopWorkflow:
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
func getNodeByID(nodes []Node, id string) Node {
	for _, node := range nodes {
		if node.Id == id {
			return node
		}
	}
	return Node{}
}

// getTopLevelNodes 获取顶级节点（没有target指向的节点）
func getTopLevelNodes(nodes []Node, edges []Edge) map[string]Node {
	inDegree := make(map[string]int)
	topLevelNodes := make(map[string]Node)

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
