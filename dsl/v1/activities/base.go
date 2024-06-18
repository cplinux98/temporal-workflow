package activities

import (
	"encoding/json"
	"github.com/pkg/errors"
	"reflect"
)

// 活动参数类型映射
var ActivityParamsMap = map[string]interface{}{
	"ShellCommandV1": ShellCommandParams{},
	// 添加其他 Activity 的参数映射
}

func GetActivityParams(activityName string, rawParams interface{}) (interface{}, error) {
	paramType, exists := ActivityParamsMap[activityName]
	if !exists {
		return nil, errors.New("unknown activity: " + activityName)
	}

	// 将原始参数转换为目标类型
	paramValue := reflect.New(reflect.TypeOf(paramType)).Interface()

	if err := mapToStruct(rawParams, paramValue); err != nil {
		return nil, err
	}

	return paramValue, nil
}

// mapToStruct 将 map 转换为结构体
func mapToStruct(data interface{}, result interface{}) error {
	// 你可以使用 json 序列化和反序列化来进行转换
	// 这里假设 rawParams 是 map[string]interface{}
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, result)
}
