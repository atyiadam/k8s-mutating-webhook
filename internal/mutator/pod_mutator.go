package mutator

import (
	"context"
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
)

func MutatePod(ctx context.Context, admissionReview admissionv1.AdmissionReview) (*admissionv1.AdmissionResponse, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("Mutation cancelled: %v", ctx.Err())
	default:
		patch, err := createPatch()
		if err != nil {
			return nil, err
		}

		patchType := admissionv1.PatchTypeJSONPatch
		return &admissionv1.AdmissionResponse{
			UID:       admissionReview.Request.UID,
			Allowed:   true,
			Patch:     patch,
			PatchType: &patchType,
		}, nil
	}
}

func createPatch() ([]byte, error) {
	patch := []struct {
		Op    string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value"`
	}{
		{
			Op:   "add",
			Path: "/metadata/labels",
			Value: map[string]string{
				"env": "development",
			},
		},
	}
	return json.Marshal(patch)
}
