package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	db "github.com/Matltin/simple-bank/db/sqlc"
	"github.com/Matltin/simple-bank/util"
	"github.com/hibiken/asynq"
)

const TaskSendVerifyEmail = "task:send_verify_email"

type PayloadSendVerfyEmail struct {
	Username string `json:"username"`
}

func (distributor *RedisTaskDistributor) DistributTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerfyEmail, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)

	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	// TODO: wirte code writen in 11:45 in 001

	fmt.Println(info)

	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerfyEmail

	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("faild to unmarshal payload: %w", asynq.SkipRetry)
	}

	user, err := processor.store.GetUser(ctx, payload.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user doesn't exist: %w", asynq.SkipRetry)
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	verifyEmail, err := processor.store.CreateVerifyEmail(ctx, db.CreateVerifyEmailParams{
		Username:   user.Username,
		Email:      user.Email,
		SecretCode: util.RandomString(32),
	})

	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	verifyUrl := fmt.Sprintf("http://simple-bank.org?id=%d&secret_code=&s", verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, user.FullName, verifyUrl)
	to := []string{user.Email}

	err = processor.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	fmt.Println(user)
	return nil
}
