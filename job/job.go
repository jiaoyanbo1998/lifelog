package job

// Job 定义一个任务接口
type Job interface {
	Name() string      // job的名字
	Run() (any, error) // 执行job
}
