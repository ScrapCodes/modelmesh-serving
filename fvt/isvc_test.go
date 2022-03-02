// Copyright 2021 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fvt

import (
	"fmt"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("Inference service", func() {
	// starting from the desired state.
	FSpecify("Preparing the cluster for inference service tests", func() {
		// ensure configuration has scale-to-zero disabled
		config := map[string]interface{}{
			// disable scale-to-zero to prevent pods flapping as
			// Predictors are created and deleted
			"scaleToZero": map[string]interface{}{
				"enabled": false,
			},
			// disable the model-mesh bootstrap failure check so
			// that the expected failures for invalid
			// tests do not trigger it
			"internalModelMeshEnvVars": []map[string]interface{}{
				{
					"name":  "BOOTSTRAP_CLEARANCE_PERIOD_MS",
					"value": "0",
				},
			},
			"podsPerRuntime": 1,
		}
		fvtClient.ApplyUserConfigMap(config)

		// ensure that there are no predictors to start
		fvtClient.DeleteAllIsvcs()

		// ensure a stable deploy state
		WaitForStableActiveDeployState()
	})
	FIt("should successfully load a model", func() {
		isvcObject := NewIsvcForFVT("new-format-mm.yaml")
		fmt.Printf("Unstructured object \n%#v\n", isvcObject)
		res := CreateIsvcAndWaitAndExpectReady(isvcObject)
		fmt.Printf("\nResult: %#v\n", res)
		// clean up
		fvtClient.DeleteIsvc(isvcObject.GetName())
	})
})
