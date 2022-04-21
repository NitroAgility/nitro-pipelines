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

import (
	"bytes"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/contexts"
)

func LoadMicroservicesFile(microservicesFile string, context *contexts.Context) (*MicroservicesModel, error) {
	var microservicesModel = &MicroservicesModel{}
    tmpl, _ := template.ParseFiles(microservicesFile)
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, context); err != nil {
		return nil, err
	}
	err := yaml.Unmarshal(tpl.Bytes(), microservicesModel)
	if err != nil {
		return nil, err
	}
	return microservicesModel, nil
}
