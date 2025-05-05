package v1alpha1

import (
	"context"
	"fmt"
	argoprojv1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/dana-team/argaas-rbac-validator/internal/webhook/clusterclient"
	"github.com/dana-team/argaas-rbac-validator/internal/webhook/config_fetcher"
	_ "github.com/dana-team/argaas-rbac-validator/internal/webhook/config_fetcher"
	"github.com/dana-team/argaas-rbac-validator/internal/webhook/permissions"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var applicationlog = logf.Log.WithName("application-webhook")

// SetupApplicationWebhookWithManager registers the webhook with the manager.
func SetupApplicationWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&argoprojv1alpha1.Application{}).
		WithValidator(&ApplicationValidator{Client: mgr.GetClient()}).
		Complete()
}

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch
// +kubebuilder:webhook:path=/validate-argoproj-io-v1alpha1-application,mutating=false,failurePolicy=fail,sideEffects=None,groups=argoproj.argoproj.io,resources=applications,verbs=create;update,versions=v1alpha1,name=vapplication-v1alpha1.kb.io,admissionReviewVersions=v1

type ApplicationValidator struct {
	client.Client
}

var _ webhook.CustomValidator = &ApplicationValidator{}

func (v *ApplicationValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	app, ok := obj.(*argoprojv1alpha1.Application)
	if !ok {
		return nil, fmt.Errorf("expected Application but got %T", obj)
	}
	applicationlog.Info("Validating Application Create", "name", app.GetName())
	return nil, validateApplication(ctx, v.Client, app)
}

func (v *ApplicationValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	app, ok := newObj.(*argoprojv1alpha1.Application)
	if !ok {
		return nil, fmt.Errorf("expected Application but got %T", newObj)
	}
	applicationlog.Info("Validating Application Update", "name", app.GetName())
	return nil, validateApplication(ctx, v.Client, app)
}

func (v *ApplicationValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateApplication(ctx context.Context, k8sClient client.Client, app *argoprojv1alpha1.Application) error {
	destNamespace := app.Spec.Destination.Namespace
	destServer := app.Spec.Destination.Server
	appNamespace := app.GetNamespace()

	if destNamespace == "" || destServer == "" {
		return fmt.Errorf("destination namespace and server must be specified")
	}

	admins, err := config_fetcher.FetchAPIAdmins(ctx, k8sClient, appNamespace)
	if err != nil {
		return fmt.Errorf("failed to fetch api-admins: %w", err)
	}

	token, err := config_fetcher.FetchClusterToken(destServer)
	if err != nil {
		return fmt.Errorf("failed to fetch cluster token: %w", err)
	}

	clusterClient, err := clusterclient.BuildClusterClient(destServer, token)
	if err != nil {
		return fmt.Errorf("failed to build client: %w", err)
	}

	for _, admin := range admins {
		allowed, err := permissions.CheckNamespaceAdmin(ctx, clusterClient, admin, destNamespace)
		if err != nil {
			return fmt.Errorf("error checking access for user %s: %w", admin, err)
		}
		if !allowed {
			return fmt.Errorf("user %s does not have admin access to namespace %s in cluster %s", admin, destNamespace, destServer)
		}
	}

	return nil
}
