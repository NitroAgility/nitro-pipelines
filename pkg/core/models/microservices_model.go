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
	Microservices []Microservices `yaml:"microservices"`
	Build         Build           `yaml:"build"`
	Deployments   Deployments     `yaml:"deployments"`
}
type Microservices struct {
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
}
type BuildExpand struct {
	Variable string `yaml:"variable"`
	Type     string `yaml:"type"`
}
type Build struct {
	BuildArgs   string `yaml:"build_args"`
	Expand   	[]BuildExpand `yaml:"expand"`
	Registry 	string   `yaml:"registry"`
}
type DeploymentExpand struct {
	Variable string `yaml:"variable"`
	Type     string `yaml:"type"`
	Name     string `yaml:"name,omitempty"`
	FileName string `yaml:"file_name,omitempty"`
}
type DeploymentScripts struct {
	PreExecution   string `yaml:"pre_execution"`
	PreDeployment  string `yaml:"pre_deployment"`
	PostDeployment string `yaml:"post_deployment"`
	PostExecution  string `yaml:"post_execution"`
}
type DeploymentHelm struct {
	Parameters string `yaml:"parameters"`
}
type Default struct {
	Expand   []DeploymentExpand `yaml:"expand"`
	Registry string   `yaml:"registry"`
	Scripts  DeploymentScripts  `yaml:"scripts"`
	Helm     DeploymentHelm     `yaml:"helm"`
}
type Deployments struct {
	Default Default `yaml:"default"`
}
