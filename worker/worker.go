package main

import (
	"ai.openclaw.usecase/namespace"
	"ai.openclaw.usecase/workflow"
	"ai.openclaw.usecase/workflow/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func main() {
	// 1. 创建 Temporal Client
	// 默认连接到本地 localhost:7233
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("无法创建 Temporal Client", err)
	}
	defer c.Close()

	// 2. 启动 Worker
	// Worker 负责监听 Task Queue，并执行具体的 Workflow 和 Activity 代码
	we := worker.New(c, namespace.TaskQueueName, worker.Options{})
	we.RegisterWorkflow(workflow.ChronicleArticleWorkflow)
	we.RegisterActivity(activity.GenerateArticleWithChronicle)

	// 3. 运行worker
	if err := we.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("无法启动 Worker", err)
	}
}
