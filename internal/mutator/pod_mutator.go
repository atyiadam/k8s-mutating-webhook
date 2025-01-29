package mutator

import (
	"encoding/json"

	admissionv1 "k8s.io/api/admission/v1"
)

func MutatePod(admissionReview admissionv1.AdmissionReview) (*admissionv1.AdmissionResponse, error) {
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
