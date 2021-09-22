/*
Copyright 2021.

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

package v1alpha1

import (
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var databaselog = logf.Log.WithName("database-resource")

func (r *Database) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-actions-msft-isd-coe-io-v1alpha1-database,mutating=true,failurePolicy=fail,sideEffects=None,groups=actions.msft.isd.coe.io,resources=databases,verbs=create;update,versions=v1alpha1,name=mdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Database{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Database) Default() {
	databaselog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-actions-msft-isd-coe-io-v1alpha1-database,mutating=false,failurePolicy=fail,sideEffects=None,groups=actions.msft.isd.coe.io,resources=databases,verbs=create;update,versions=v1alpha1,name=vdatabase.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Database{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateCreate() error {
	databaselog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateUpdate(old runtime.Object) error {
	databaselog.Info("validate update", "name", r.Name)

	var allErrs field.ErrorList
	// TODO(user): fill in your validation logic upon object update.
	curr := old.(*Database)
	if curr == nil {
		return fmt.Errorf("could not convert runtime.Object to Database")
	}

	if r.Spec.Name != curr.Spec.Name {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("name"), r.Spec.Name, "cannot rename the database"))
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "actions.msft.isd.coe.io", Kind: "Database"},
			r.Name, allErrs)
	}
	if r.Spec.Collation != curr.Spec.Collation {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec").Child("collation"), r.Spec.Collation, "cannot change the collation of the database"))
		return apierrors.NewInvalid(
			schema.GroupKind{Group: "actions.msft.isd.coe.io", Kind: "Database"},
			r.Name, allErrs)
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Database) ValidateDelete() error {
	databaselog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
