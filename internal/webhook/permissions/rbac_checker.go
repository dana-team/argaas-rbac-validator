package permissions

import (
	"context"
	"fmt"
	authv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckNamespaceAdmin checks if the user has admin access to namespace.
func CheckNamespaceAdmin(ctx context.Context, client kubernetes.Interface, user, namespace string) (bool, error) {
	sar := &authv1.SubjectAccessReview{
		Spec: authv1.SubjectAccessReviewSpec{
			User: user,
			ResourceAttributes: &authv1.ResourceAttributes{
				Namespace: namespace,
				Verb:      "*",
				Resource:  "*",
			},
		},
	}

	resp, err := client.AuthorizationV1().SubjectAccessReviews().Create(ctx, sar, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf("SubjectAccessReview failed: %w", err)
	}

	return resp.Status.Allowed, nil
}
