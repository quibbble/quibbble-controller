package sdkserver

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Patch struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

// stale routinely patches the stale label.
func (s *SDKServer) stale() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stale := len(s.players) == 0 && s.updatedAt.Add(15*time.Minute).Before(time.Now())
		patch := []Patch{
			{
				Op:    "replace",
				Path:  "/metadata/labels/stale",
				Value: strconv.FormatBool(stale),
			},
		}
		raw, err := json.Marshal(patch)
		if err != nil {
			continue
		}
		if _, err := s.clientset.CoreV1().Pods(s.namespace).Patch(context.Background(), s.pod, types.JSONPatchType, raw, metav1.PatchOptions{}); err != nil {
			continue
		}
	}
}
