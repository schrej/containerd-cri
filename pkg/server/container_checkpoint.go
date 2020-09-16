package server

import (
	"syscall"

	"github.com/containerd/containerd"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

func (c *criService) CheckpointContainer(ctx context.Context, r *runtime.CheckpointContainerRequest) (retRes *runtime.CheckpointContainerResponse, retErr error) {
	cntr, err := c.containerStore.Get(r.GetContainerId())
	if err != nil {
		return nil, errors.Wrap(err, "failed to find container")
	}
	task, err := cntr.Container.Task(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to checkpoint container task")
	}
	opts := []containerd.CheckpointTaskOpts{containerd.WithCheckpointImagePath(r.GetOptions().GetCheckpointPath())}
	if !r.GetOptions().LeaveRunning {
		opts = append(opts, containerd.WithCheckpointExit())
	}
	_, err = task.Checkpoint(ctx, opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to checkpoint container")
	}

	if !r.GetOptions().GetLeaveRunning() {
		task.Kill(ctx, syscall.SIGKILL)
	}

	return &runtime.CheckpointContainerResponse{}, nil
}
