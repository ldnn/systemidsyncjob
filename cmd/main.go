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

package main

import (

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"fmt"
	"os"
	"time"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"systemidsyncjob/internal/controller"
)

func main() {
	var c controller.Coordinator
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}

	c.Client, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	c.Gvr = schema.GroupVersionResource{
		Group:    "tenant.kubesphere.io",
		Version:  "v1alpha2",
		Resource: "workspacetemplates",
	}
	/*
		c.Systems = map[string]string{
			"YW-DS-SZSL":                           "中移数智商旅",
			"YW-HB-TYZT":                           "通用业务共享能力平台",
			"ZC-IT-YPT":                            "云平台",
			"471455af-0883-41b1-8e71-4c3abb981910": "应用同城双活管控平台",
		}
	*/
	var ok bool
	if c.Systems, ok = controller.GetSystem(); !ok {
		panic(fmt.Errorf("Error: failed to dump cmdb,time:%v", time.Now().Unix()))
	}

	workspaceList := c.Get()
	var id, name string
	for _, d := range workspaceList.Items {
		id = d.Labels["kubesphere.io/it-system-id"]
		name = d.Annotations["kubesphere.io/description"]
		_, ok := c.Systems[id]
		if ok && c.Systems[id] != name {
			glog.Infof("Sync: patch=%v,system=%v\n", d.Name, c.Systems[id])
			c.Update(d.Name, c.Systems[id])
		}

	}
	defer glog.Infoln("Completed: sync system to workspace!!!")
	os.Exit(0)

}
