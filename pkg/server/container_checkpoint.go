package server

import (
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
	// opts := []containerd.CheckpointOpts{containerd.WithCheckpointTask, containerd.WithCheckpointImage, containerd.WithCheckpointRuntime, containerd.WithCheckpointRW}
	// if !r.GetOptions().GetLeaveRunning() {
	// 	opts = append(opts, containerd.WithCheckpointTaskExit)
	// }
	task, err := cntr.Container.Task(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to checkpoint container task")
	}
	_, err = task.Checkpoint(ctx, containerd.WithCheckpointImagePath(r.GetOptions().GetCheckpointPath()))
	//img, err := cntr.Container.Checkpoint(ctx, r.GetOptions().GetCheckpointPath(), opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to checkpoint container")
	}

	if !r.GetOptions().GetLeaveRunning() {
		task.Pause(ctx)
	}

	return &runtime.CheckpointContainerResponse{}, nil
}
