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
package models

type MicroservicesModel struct {
	Settings	  Setting `yaml:"settings"`
	Microservices []Microservices `yaml:"microservices"`
	Build         Build           `yaml:"build"`
	Deployments   Deployments     `yaml:"deployments"`
}
type Setting struct {
	Environment map[string]EnvironemntSetting `yaml:"environments"`
}

type EnvironemntSetting struct {
	RepoStrategoy	string `yaml:"repo_strategy"`
	PromotionStrategy	string `yaml:"promotion_strategy"`
}

type Microservices struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
}
type Expand struct {
	Variable string `yaml:"variable"`
	Type     string `yaml:"type"`
	Name     string `yaml:"name,omitempty"`
}
type Build struct {
	BuildArgs   string 		`yaml:"build_args"`
	Expand   	[]Expand 	`yaml:"expand"`
	Registry 	string   	`yaml:"registry"`
}
type DeploymentScripts struct {
	PreExecution   string `yaml:"pre_execution"`
	PrePromotion   string `yaml:"pre_promotion"`
	PostPromotion  string `yaml:"post_promotion"`
	PreDeployment  string `yaml:"pre_deployment"`
	PostDeployment string `yaml:"post_deployment"`
	PostExecution  string `yaml:"post_execution"`
}
type DeploymentHelm struct {
	Parameters string `yaml:"parameters"`
}
type Default struct {
	Expand   []Expand 			`yaml:"expand"`
	Scripts  DeploymentScripts  `yaml:"scripts"`
	Helm     DeploymentHelm     `yaml:"helm"`
}
type Deployments struct {
	Default Default `yaml:"default"`
}
