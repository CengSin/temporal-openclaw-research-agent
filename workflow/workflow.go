package workflow

import (
	"ai.openclaw.usecase/workflow/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func ChronicleArticleWorkflow(ctx workflow.Context, topic string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    20 * time.Second,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var article string
	logger := workflow.GetLogger(ctx)
	logger.Info("Workflow: 调用 Chronicle 生成文章", "topic", topic)

	if err := workflow.ExecuteActivity(ctx, activity.GenerateArticleWithChronicle, topic).Get(ctx, &article); err != nil {
		logger.Error("Workflow: 生成文章失败", "error", err)
		return "", err
	}

	logger.Info("Workflow 完成", "length", len(article))
	return article, nil
}
