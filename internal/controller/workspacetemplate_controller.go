/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	tenantv1alpha2 "kubesphere.io/api/tenant/v1alpha2"
)

func (c *Coordinator) Get() *tenantv1alpha2.WorkspaceTemplateList {

	unstructObj, err := c.Client.Resource(c.Gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	workspaceList := &tenantv1alpha2.WorkspaceTemplateList{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), workspaceList)

	if err != nil {
		panic(err)
	}

	return workspaceList

}

func (c *Coordinator) Update(workspce string, system string) bool {

	unstructObj, err := c.Client.Resource(c.Gvr).Get(context.TODO(), workspce, metav1.GetOptions{})
	if err != nil {
		glog.Error(err.Error())
		return false

	}

	if err := unstructured.SetNestedField(unstructObj.Object, system, "metadata", "annotations", "kubesphere.io/description"); err != nil {
		glog.Error(fmt.Errorf("failed to set max.ingress-bandwidth value: %v", err))
		return false
	}
	_, updateErr := c.Client.Resource(c.Gvr).Update(context.TODO(), unstructObj, metav1.UpdateOptions{})
	if updateErr != nil {
		glog.Error(updateErr.Error())
		return false
	}
	return true

}
