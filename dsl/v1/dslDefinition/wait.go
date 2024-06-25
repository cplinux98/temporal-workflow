package dslDefinition

type ResumeType string

const (
	Interval      ResumeType = "interval"      // 间隔时间，毫秒
	FixedTime     ResumeType = "fixedTime"     // 固定时间，格式: "2024-06-19 10:10:10"
	WebhookSignal ResumeType = "webhookSignal" // Webhook 的信号名称
)

type WaitParams struct {
	Resume         ResumeType `yaml:"resume"`                   // 恢复类型
	IntervalMillis int64      `yaml:"intervalMillis,omitempty"` // 间隔时间，毫秒
	FixedTime      string     `yaml:"fixedTime,omitempty"`      // 固定时间，格式: "2024-06-19 10:10:10"
	WebhookSignal  string     `yaml:"webhookSignal,omitempty"`  // Webhook 的信号名称
}

type WaitResult struct {
	//Code   int64  `yaml:"code"`   // 退出状态码
	//Stdout string `yaml:"stdout"` // 标准输出
	//Stderr string `yaml:"stderr"` // 标准错误输出
}
