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
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/contexts"
)

const DeployTpl = `#!/bin/bash
# Configure local files
export KUBECONFIG="./kube_config"
export AWS_CONFIG_FILE="./aws_config"
export AWS_SHARED_CREDENTIALS_FILE="./aws_credentials"
touch $KUBECONFIG
chmod 777 $KUBECONFIG
touch $AWS_CONFIG_FILE
chmod 777 $AWS_CONFIG_FILE
touch $AWS_SHARED_CREDENTIALS_FILE
chmod 777 $AWS_SHARED_CREDENTIALS_FILE
# Pre execution
{{ .PreExecution }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Expanding variables
{{ range .Expand -}}
filename=$(uuidgen)
if [ $MACHINE_OS == "OSX" ]; then
	echo ${{ .Variable }} | base64 --decode >> $NITROBIN/{{ .Name }}.tmp && envsubst < $NITROBIN/{{ .Name }}.tmp > $NITROBIN/{{ .Name }}.env && rm $NITROBIN/{{ .Name }}.tmp
else
	echo ${{ .Variable }} | base64 -di >> $NITROBIN/{{ .Name }}.tmp && envsubst < $NITROBIN/{{ .Name }}.tmp > $NITROBIN/{{ .Name }}.env && rm $NITROBIN/{{ .Name }}.tmp
fi
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ if eq .Type "environment" -}}
source $NITROBIN/{{ .Name }}.env && export $(cut -d= -f1 $NITROBIN/{{ .Name }}.env)
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm -f $NITROBIN/{{ .Name -}}.env
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ end -}}
{{ end -}}
# Environment configuration
aws configure set aws_access_key_id $NITRO_PIPELINES_VARIABLES_SOURCE_AWS_ACCESS_KEY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws configure set aws_secret_access_key $NITRO_PIPELINES_VARIABLES_SOURCE_AWS_SECRET_ACCESS
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr get-login-password --region $NITRO_PIPELINES_VARIABLES_SOURCE_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Pull docker images
{{ range .Images -}}
docker pull $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ end -}}
# Push docker images
aws configure set aws_access_key_id $NITRO_PIPELINES_VARIABLES_TARGET_AWS_ACCESS_KEY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws configure set aws_secret_access_key $NITRO_PIPELINES_VARIABLES_TARGET_AWS_SECRET_ACCESS
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr get-login-password --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ range .Images -}}
aws ecr create-repository --no-cli-pager --repository-name {{ .TargetImageName }} --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION || true
docker tag $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }}:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker tag $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }}:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }}:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ end -}}
# EKS Deployment
aws eks --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION update-kubeconfig --name $NITRO_PIPELINES_VARIABLES_TARGET_AWS_EKS_CLUSTER_NAME
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Pre deployment
{{ .PreDeployment }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
helm upgrade --install $NITRO_PIPELINES_VARIABLES_TARGET_HELM_RELEASE_NAME "$NITRO_PIPELINES_VARIABLES_TARGET_HELM_CHART_CODE_PATH/chart/$NITRO_PIPELINES_VARIABLES_TARGET_HELM_CHART_NAME" --set environment={{ .Environment }} --set infrastructure.docker_registry=$NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY --set app.tag=$NITRO_PIPELINES_BUILD_NUMBER {{ .HelmArgs }} -n $NITRO_PIPELINES_VARIABLES_TARGET_HELM_NAMESPACE
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Post deployment
{{ .PostDeployment }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm "$NITROBIN/aws_credentials"
rm "$NITROBIN/aws_config"
rm "$NITROBIN/kube_config"
# Cleaning expanded variables
{{ range .Expand -}}
{{ if eq .Type "file" -}}
rm -f $NITROBIN/{{ .Name -}}.env
{{ end -}}
{{ end -}}
# Post execution
{{ .PostExecution }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
`

func ExecuteDeploy(deployCtx *contexts.DeployContext) error {
	tmpl, _ := template.New("DEPLOY").Parse(DeployTpl)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, deployCtx); err != nil {
		fmt.Print(buffer.String())
		return err
	}
	if strings.ToUpper(os.Getenv("DRY_RUN")) == "TRUE" {
		fmt.Println(buffer.String())
	} else {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		if err := saveToFile(exPath + "/nitro-deploy.sh", buffer.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
