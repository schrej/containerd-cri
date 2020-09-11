package server

import (
	"github.com/containerd/containerd"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// RestoreContainer restores a container from a previously created image.
// It's essentially the same as starting a container with the additon of loading a checkpoint.
func (c *criService) RestoreContainer(ctx context.Context, r *runtime.RestoreContainerRequest) (retRes *runtime.RestoreContainerResponse, retErr error) {

	if err := c.startContainer(ctx, r.GetContainerId(), containerd.WithRestoreImagePath(r.GetOptions().GetCheckpointPath())); err != nil {
		return nil, errors.Wrap(err, "failed to restore container")
	}
	return &runtime.RestoreContainerResponse{}, nil
}
