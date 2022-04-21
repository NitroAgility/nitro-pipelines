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
package commands

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/contexts"
)

const buildTpl = `#!/bin/bash
aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
aws ecr create-repository --repository-name {{ .ImageName }} --region $NITRO_PIPELINES_TARGET_AWS_REGION || true
docker build -t {{ .ImageName }}:latest {{ .DockerArgs }} -f {{ .Dockerfile }} .
docker tag {{ .ImageName }}:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:latest
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:latest
docker tag {{ .ImageName }}:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
`

func ExecuteBuild(buildCtx *contexts.BuildContext) (error) {
    tmpl, _ :=  template.New("BUILD").Parse(buildTpl)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, buildCtx); err != nil {
        fmt.Print(buffer.String())
		return err
	}
    fmt.Print(buffer.String())
    return nil
}
