package istiocsr

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/openshift/cert-manager-operator/api/operator/v1alpha1"
	"github.com/openshift/cert-manager-operator/pkg/controller/istiocsr/fakes"
)

func TestUpdateStatus_PreservesConditions(t *testing.T) {
	tests := []struct {
		name                string
		existingConditions  []metav1.Condition
		newConditions       []metav1.Condition
		expectedConditions  []string
		wantErr             bool
	}{
		{
			name: "preserves existing conditions when updating with no conditions",
			existingConditions: []metav1.Condition{
				{
					Type:   v1alpha1.Ready,
					Status: metav1.ConditionTrue,
					Reason: v1alpha1.ReasonReady,
				},
				{
					Type:   v1alpha1.Degraded,
					Status: metav1.ConditionFalse,
					Reason: v1alpha1.ReasonReady,
				},
			},
			newConditions:      []metav1.Condition{},
			expectedConditions: []string{v1alpha1.Ready, v1alpha1.Degraded},
		},
		{
			name: "merges new condition with existing conditions",
			existingConditions: []metav1.Condition{
				{
					Type:   v1alpha1.Ready,
					Status: metav1.ConditionTrue,
					Reason: v1alpha1.ReasonReady,
				},
			},
			newConditions: []metav1.Condition{
				{
					Type:   v1alpha1.Degraded,
					Status: metav1.ConditionFalse,
					Reason: v1alpha1.ReasonReady,
				},
			},
			expectedConditions: []string{v1alpha1.Ready, v1alpha1.Degraded},
		},
		{
			name: "overrides existing condition with same type",
			existingConditions: []metav1.Condition{
				{
					Type:    v1alpha1.Ready,
					Status:  metav1.ConditionTrue,
					Reason:  v1alpha1.ReasonReady,
					Message: "old message",
				},
			},
			newConditions: []metav1.Condition{
				{
					Type:    v1alpha1.Ready,
					Status:  metav1.ConditionFalse,
					Reason:  v1alpha1.ReasonInProgress,
					Message: "new message",
				},
			},
			expectedConditions: []string{v1alpha1.Ready},
		},
		{
			name:               "sets new conditions when none exist",
			existingConditions: []metav1.Condition{},
			newConditions: []metav1.Condition{
				{
					Type:   v1alpha1.Ready,
					Status: metav1.ConditionTrue,
					Reason: v1alpha1.ReasonReady,
				},
			},
			expectedConditions: []string{v1alpha1.Ready},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := testReconciler(t)
			mock := &fakes.FakeCtrlClient{}

			// Setup mock to return an IstioCSR with existing conditions
			mock.GetCalls(func(ctx context.Context, key types.NamespacedName, obj client.Object) error {
				if istiocsr, ok := obj.(*v1alpha1.IstioCSR); ok {
					istiocsr.Name = testResourcesName
					istiocsr.Namespace = testIstioCSRNamespace
					istiocsr.Status.Conditions = tt.existingConditions
				}
				return nil
			})

			var capturedStatus *v1alpha1.IstioCSRStatus
			mock.StatusUpdateCalls(func(ctx context.Context, obj client.Object, option ...client.SubResourceUpdateOption) error {
				if istiocsr, ok := obj.(*v1alpha1.IstioCSR); ok {
					// Capture the status being updated
					capturedStatus = &v1alpha1.IstioCSRStatus{}
					istiocsr.Status.DeepCopyInto(capturedStatus)
				}
				return nil
			})

			r.ctrlClient = mock

			// Create IstioCSR object with new conditions
			istiocsr := testIstioCSR()
			istiocsr.Status.Conditions = tt.newConditions

			// Call updateStatus
			err := r.updateStatus(context.Background(), istiocsr)
			if (err != nil) != tt.wantErr {
				t.Errorf("updateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify conditions
			if capturedStatus == nil {
				t.Error("StatusUpdate was not called")
				return
			}

			// Check that all expected conditions are present
			for _, expectedType := range tt.expectedConditions {
				found := false
				for _, condition := range capturedStatus.Conditions {
					if condition.Type == expectedType {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected condition type %s not found in status", expectedType)
				}
			}

			// For override case, verify the new values
			if tt.name == "overrides existing condition with same type" {
				for _, condition := range capturedStatus.Conditions {
					if condition.Type == v1alpha1.Ready {
						if condition.Status != metav1.ConditionFalse {
							t.Errorf("Expected Ready condition to be overridden to False, got %v", condition.Status)
						}
						if condition.Reason != v1alpha1.ReasonInProgress {
							t.Errorf("Expected Ready condition reason to be overridden to %s, got %s", v1alpha1.ReasonInProgress, condition.Reason)
						}
					}
				}
			}
		})
	}
}
