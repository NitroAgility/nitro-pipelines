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

type PromotionExpandContext struct {
	Variable string
	Type	string
	Name	string
}

type PromotionImageContext struct {
	Name			string
	SourceImageName string
	TargetImageName	string
}

type PromotionContext struct {
	Name			string
	Environment 	string
	PreExecution	string
	PostExecution	string
	PrePromotion	string
	PostPromotion	string
	Expand			[]PromotionExpandContext
	Images			[]PromotionImageContext
	Strategy		string
}

// Creational functions

func NewPromotionContext(microservicesFile string, name string)  (*PromotionContext, error) {
	msModel, err := loadMicroservicesFile(microservicesFile)
	if err != nil {
		return nil, err
	}
	envSource := strings.ToUpper(os.Getenv("ENV_SOURCE"))
	envTarget := strings.ToUpper(os.Getenv("ENV_TARGET"))
	if len(envTarget) == 0 {
		return nil, errors.New("target environment cannot be null")
	}
	if len(envSource) == 0 {
		envSource = "BUILD"
	}
	strategy := "push"
	env := strings.ToLower(envTarget)
	if val, ok := msModel.Settings.Environment[env]; ok {
		if strings.ToLower(val.PromotionStrategy) == "retag" {
			strategy = "retag"
		}
	}
	context := &PromotionContext {
		Name			: name,
		Environment 	: strings.ToLower(os.Getenv("ENV_TARGET")),
		PreExecution	: buildScript(msModel.Deployments.Default.Scripts.PreExecution),
		PostExecution	: buildScript(msModel.Deployments.Default.Scripts.PostExecution),
		PrePromotion	: buildScript(msModel.Deployments.Default.Scripts.PrePromotion),
		PostPromotion	: buildScript(msModel.Deployments.Default.Scripts.PostPromotion),
		Expand			: []PromotionExpandContext{},
		Images			: []PromotionImageContext{},
		Strategy		: strategy,
	}
	for _, e := range msModel.Deployments.Default.Expand {
		expCtx := PromotionExpandContext {}
		expCtx.Variable = e.Variable
		expCtx.Type 	= e.Type
		expCtx.Name 	= buildFileName(e.Name)
		context.Expand	= append(context.Expand, expCtx)
	}
	for _, m := range msModel.Microservices {
		if name == "" || m.Name == name {
			msCtx := PromotionImageContext {}
			msCtx.Name				= m.Name
			msCtx.SourceImageName 	= buildImageName(m.Name, envSource)
			msCtx.TargetImageName 	= buildImageName(m.Name, envTarget)
			context.Images			= append(context.Images, msCtx)
		}
	}
	return context, nil
}
