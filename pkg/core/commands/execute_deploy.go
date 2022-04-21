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

const DeployTpl = `#!/bin/bash
# Pre execution
{{ .PreExecution }}
# Environment configuration
aws configure set aws_access_key_id $NITRO_PIPELINES_SOURCE_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_SOURCE_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_SOURCE_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_SOURCE_DOCKER_REGISTRY
# Pull docker images
{{ range .Images -}}
docker pull $NITRO_PIPELINES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
{{ end -}}
# Push docker images
aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
{{ range .Images -}}
aws ecr create-repository --repository-name {{ .TargetImageName }} --region $NITRO_PIPELINES_TARGET_AWS_REGION || true
docker tag $NITRO_PIPELINES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }}:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
docker tag $NITRO_PIPELINES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }}:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:latest
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:latest
{{ end -}}
# EKS Deployment
aws eks --region $NITRO_PIPELINES_TARGET_AWS_REGION update-kubeconfig --name $NITRO_PIPELINES_TARGET_AWS_EKS_CLUSTER_NAME
# Pre deployment
{{ .PreDeployment }}
helm upgrade --install $NITRO_PIPELINES_TARGET_HELM_RELEASE_NAME "$NITRO_PIPELINES_TARGET_HELM_CHART_CODE_PATH/chart/$NITRO_PIPELINES_TARGET_HELM_CHART_NAME" --set environment={{ .Environment }} --set infrastructure.domain="$NITRO_PIPELINES_DOMAIN" --set infrastructure.docker_registry=$NITRO_PIPELINES_TARGET_DOCKER_REGISTRY --set app.tag=$NITRO_PIPELINES_BUILD_NUMBER {{ .HelmArgs }} -n $NITRO_PIPELINES_TARGET_HELM_NAMESPACE
# Post deployment
{{ .PostDeployment }}
# Post execution
{{ .PostExecution }}
`

func ExecuteDeploy(deployCtx *contexts.DeployContext) (error) {
    tmpl, _ :=  template.New("DEPLOY").Parse(DeployTpl)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, deployCtx); err != nil {
        fmt.Print(buffer.String())
		return err
	}
    fmt.Print(buffer.String())
    return nil
}
