/*
Copyright 2022 Robert Nemet.

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

package pkg_webhook

import (
	"context"
	"encoding/json"
	"fmt"

	klockv1 "klock/apis/klock/v1"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var log = logf.Log.WithName("klock-validate-hook")

//+kubebuilder:rbac:groups=klock.rnemet.dev,resources=locks,verbs=get;list;watch
//+kubebuilder:webhook:verbs=delete;update,path=/validate-all,admissionReviewVersions="v1",sideEffects=None,mutating=false,failurePolicy=fail,groups=*,resources=pods;deployments;secrets,versions=*,name=klocks.rnemet.dev

// validator validates resources
type validator struct {
	Client client.Client
}

func NewValidator(c client.Client) admission.Handler {
	return &validator{Client: c}
}

// podValidator admits a pod if a specific annotation exists.
func (v *validator) Handle(ctx context.Context, req admission.Request) admission.Response {

	log.Info("enter webhook")

	lockList, err := v.getLocks(ctx, req)

	if err != nil {
		log.Error(err, "error while getting LockList")
		return admission.Allowed("")
	}

	var target unstructured.Unstructured

	err = json.Unmarshal(req.OldObject.Raw, &target)
	if err != nil {
		log.Error(err, "can not unmarchal OldObject => allow")
		return admission.Allowed("can not unmarchal OldObject => allow")
	}

	labels := target.GetLabels()

	for _, lock := range lockList.Items {
		if v.matchOperation(string(req.Operation), lock.Spec.Operations) {

			for k, v := range lock.Spec.Matcher {
				if labels[k] == v {
					return admission.Denied(fmt.Sprintf("denied, there is a lock: %v", lock.Spec.Matcher))
				}
			}
		}
	}

	return admission.Allowed("")
}

// matchOperation, check if desired operatio is covered by lock. In essence this fn just check if element is in array
func (v *validator) matchOperation(op string, operations []string) bool {
	for _, o := range operations {
		if o == op {
			return true
		}
	}
	return false
}

func (v *validator) getLocks(ctx context.Context, request admission.Request) (lockList klockv1.LockList, err error) {

	opts := client.ListOptions{
		Namespace: request.Namespace,
	}

	if err = v.Client.List(ctx, &lockList, &opts); err != nil {
		log.Error(err, "error while getting LockList")
	}

	log.Info(fmt.Sprintf("locks found %v", lockList))
	return lockList, err
}
