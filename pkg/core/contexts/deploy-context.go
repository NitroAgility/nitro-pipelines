/*
Copyright 2021 Nitro Agility S.r.l..
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
package contexts

import (
	"errors"
	"os"
	"strings"
)

type DeployExpandContext struct {
	Variable string
	Type	string
	Name	string
}

type DeployImageContext struct {
	SourceImageName string
	TargetImageName	string
}

type DeployContext struct {
	Environment 	string
	PreExecution	string
	PostExecution	string
	PreDeployment	string
	PostDeployment	string
	Expand			[]DeployExpandContext
	Images			[]DeployImageContext
	HelmArgs		string
}

// Creational functions

func NewDeployContext(microservicesFile string)  (*DeployContext, error) {
	msModel, err := loadMicroservicesFile(microservicesFile)
	envSource := strings.ToUpper(os.Getenv("ENV_SOURCE"))
	envTarget := strings.ToUpper(os.Getenv("ENV_TARGET"))
	if len(envTarget) == 0 {
		return nil, errors.New("target environment cannot be null")
	}
	if err != nil {
		return nil, err
	}
	context := &DeployContext {
		Environment 	: strings.ToUpper(os.Getenv("ENV")),
		PreExecution	: buildScript(msModel.Deployments.Default.Scripts.PreExecution),
		PostExecution	: buildScript(msModel.Deployments.Default.Scripts.PostExecution),
		PreDeployment	: buildScript(msModel.Deployments.Default.Scripts.PreDeployment),
		PostDeployment	: buildScript(msModel.Deployments.Default.Scripts.PostDeployment),
		Expand			: []DeployExpandContext{},
		Images			: []DeployImageContext{},
		HelmArgs		: msModel.Deployments.Default.Helm.Parameters,
	}
	for _, e := range msModel.Deployments.Default.Expand {
		expCtx := DeployExpandContext {}
		expCtx.Variable = e.Variable
		expCtx.Type 	= e.Type
		expCtx.Name 	= buildFileName(e.Name)
		context.Expand	= append(context.Expand, expCtx)
	}
	for _, m := range msModel.Microservices {
		msCtx := DeployImageContext {}
		msCtx.SourceImageName 	= buildImageName(m.Name, envSource)
		msCtx.TargetImageName 	= buildImageName(m.Name, envTarget)
		context.Images			= append(context.Images, msCtx)
	}
	return context, nil
}
