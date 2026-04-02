/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package sbserver

import (
	"context"
	"fmt"

	"github.com/containerd/log"
	runtime "k8s.io/cri-api/pkg/apis/runtime/v1"

	customtransfer "github.com/containerd/containerd/pkg/transfer/custom"
)

// CustomSnapshotterName is the snapshotter name that triggers the custom pull path.
const CustomSnapshotterName = "custom"

// customPullImage bypasses the standard OCI pull pipeline and delegates
// image fetching to the black-box custom transfer service.
func (c *criService) customPullImage(ctx context.Context, imageRef string) (*runtime.PullImageResponse, error) {
	ts := customtransfer.GetInstance()
	if ts == nil {
		return nil, fmt.Errorf("custom pull: transfer service not registered")
	}

	img, err := ts.PullImage(ctx, imageRef)
	if err != nil {
		return nil, fmt.Errorf("custom pull failed for %q: %w", imageRef, err)
	}

	imageID := img.Target.Digest.String()
	log.G(ctx).Infof("custom pull: completed for %q, imageID=%s", imageRef, imageID)
	return &runtime.PullImageResponse{ImageRef: imageID}, nil
}
