package workflow

import (
	"ai.openclaw.usecase/workflow/activity"
	"os"
	"strings"
	"testing"
	"time"

	"go.temporal.io/sdk/testsuite"
)

func TestChronicleArticleWorkflow_Integration(t *testing.T) {
	if os.Getenv("OPENCLAW_INTEGRATION") != "1" {
		t.Skip("设置 OPENCLAW_INTEGRATION=1 后执行本地 OpenClaw 集成测试")
	}

	topic := os.Getenv("OPENCLAW_TEST_TOPIC")
	if strings.TrimSpace(topic) == "" {
		topic = "新能源行业周度观察"
	}

	var ts testsuite.WorkflowTestSuite
	env := ts.NewTestWorkflowEnvironment()
	env.SetTestTimeout(30 * time.Second)
	env.RegisterWorkflow(ChronicleArticleWorkflow)
	env.RegisterActivity(activity.GenerateArticleWithChronicle)

	env.ExecuteWorkflow(ChronicleArticleWorkflow, topic)

	if !env.IsWorkflowCompleted() {
		t.Fatal("workflow 未完成")
	}

	if err := env.GetWorkflowError(); err != nil {
		t.Fatalf("workflow 失败: %v", err)
	}

	var article string
	if err := env.GetWorkflowResult(&article); err != nil {
		t.Fatalf("读取 workflow 结果失败: %v", err)
	}

	if strings.TrimSpace(article) == "" {
		t.Fatal("workflow 返回空文章")
	}
}
