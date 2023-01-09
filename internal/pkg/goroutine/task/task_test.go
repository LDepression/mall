//go:build task
// +build task

package task_test

import (
	"context"
	"testing"
	"time"

	"ttms/internal/pkg/goroutine/task"

	"github.com/stretchr/testify/require"
)

func TestNewTickerTask(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		f    func()
	}{
		{
			name: "good task",
			f: func() {
				nums := 0
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				task.NewTickerTask(task.Task{
					Name:            "good task",
					Ctx:             ctx,
					TaskDuration:    500 * time.Millisecond,
					TimeoutDuration: 1 * time.Second,
					F: func(ctx context.Context) {
						nums++
					},
				})
				time.Sleep(6 * time.Second)
				cancel()
				require.True(t, nums >= 10)
			},
		},
		{
			name: "bad task",
			f: func() {
				badNum := 0
				newNum := 0
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				task.NewTickerTask(task.Task{Name: "badTask", Ctx: ctx, TaskDuration: 500 * time.Millisecond, TimeoutDuration: 500 * time.Millisecond, F: func(ctx context.Context) {
					newNum++
					select {
					case <-time.Tick(time.Second):
						badNum++
					case <-ctx.Done():
						return
					}
				}})
				time.Sleep(6 * time.Second)
				cancel()
				require.Equal(t, 0, badNum)
				require.True(t, newNum > 0)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f()
		})
	}
}
