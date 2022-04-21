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
	"fmt"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/models"
)

type BuildExpandContext struct {
	Variable string
	Type	string
	Name	string
}

type BuildContext struct {
	Expand		[]BuildExpandContext
	Dockerfile 	string
	DockerArgs 	string
	ImageName	string
}

// Creational functions

func NewBuildContext(microservicesFile string, msName string) (*BuildContext, error) {
	msModel, err := loadMicroservicesFile(microservicesFile)
	if err != nil {
		return nil, err
	}
	var microservice *models.Microservices
	for _, m := range msModel.Microservices {
		if m.Name == msName {
			microservice = &m
		}
	}
	if microservice == nil {
		return nil, errors.New("invalid microservice name")
	}
	context := &BuildContext {
		Dockerfile: microservice.Dockerfile,
		DockerArgs: msModel.Build.BuildArgs,
		ImageName: fmt.Sprintf("build-%s", microservice.Name),
	}
	for _, e := range msModel.Deployments.Default.Expand {
		expCtx := BuildExpandContext {}
		expCtx.Variable = e.Variable
		expCtx.Type 	= e.Type
		expCtx.Name 	= buildFileName(e.Name)
		context.Expand	= append(context.Expand, expCtx)
	}
	return context, nil
}
