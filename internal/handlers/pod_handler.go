package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atyiadam/k8s-mutating-webhook/internal/mutator"
	"github.com/atyiadam/k8s-mutating-webhook/pkg/utils"
)

func MutatePod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.LogError(err, "Error reading body")
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var admissionReview admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &admissionReview); err != nil {
		utils.LogError(err, "Error unmarshaling AdmissionReview")
		http.Error(w, "Error unmarshaling AdmissionReview", http.StatusBadRequest)
		return
	}

	admissionResponse, err := mutator.MutatePod(r.Context(), admissionReview)
	if err != nil {
		utils.LogError(err, "Error mutating pod")
		http.Error(w, "Error mutating pod", http.StatusInternalServerError)
		return
	}

	responseAdmissionReview := admissionv1.AdmissionReview{
		TypeMeta: v1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: admissionResponse,
	}

	responseBytes, err := json.Marshal(responseAdmissionReview)
	if err != nil {
		utils.LogError(err, "Error marshaling response")
		http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseBytes); err != nil {
		utils.LogError(err, "Error writing response")
	}
}
