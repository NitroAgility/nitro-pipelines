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

const PromotionTpl = `#!/bin/bash
# Configure local files
export AWS_CONFIG_FILE="$NITROBIN"aws_config
export AWS_SHARED_CREDENTIALS_FILE="$NITROBIN"aws_credentials
# Pre execution
{{ .PreExecution }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Expanding variables
{{ range .Expand -}}
filename=$(uuidgen)
if [ $MACHINE_OS == "OSX" ]; then
	echo ${{ .Variable }} | base64 --decode >> "$NITROBIN"{{ .Name }}.tmp && envsubst < "$NITROBIN"{{ .Name }}.tmp > "$NITROBIN"{{ .Name }}.env && rm "$NITROBIN"{{ .Name }}.tmp
else
	echo ${{ .Variable }} | base64 -di >> "$NITROBIN"{{ .Name }}.tmp && envsubst < "$NITROBIN"{{ .Name }}.tmp > "$NITROBIN"{{ .Name }}.env && rm "$NITROBIN"{{ .Name }}.tmp
fi
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ if eq .Type "environment" -}}
source "$NITROBIN"{{ .Name }}.env && export $(cut -d= -f1 "$NITROBIN"{{ .Name }}.env)
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm -f "$NITROBIN"{{ .Name -}}.env
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
# Pre promotion
{{ .PrePromotion }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{if eq .Strategy "retag" -}}
#Retagging an image
{{ range .Images -}}
MANIFEST=$(aws ecr batch-get-image --repository-name $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/{{ .SourceImageName }} --image-ids imageTag=$NITRO_PIPELINES_BUILD_NUMBER --output json | jq --raw-output --join-output '.images[0].imageManifest')
aws ecr put-image --repository-name $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .TargetImageName }} --image-tag $NITRO_PIPELINES_BUILD_NUMBER --image-manifest "$MANIFEST"
{{ end -}}
{{end -}}
{{if eq .Strategy "push"-}}
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
{{end -}}
# Post promotion
{{ .PostPromotion }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Cleaning expanded variables
{{ range .Expand -}}
{{ if eq .Type "file" -}}
rm -f "$NITROBIN"{{ .Name -}}.env
{{ end -}}
{{ end -}}
# Post execution
{{ .PostExecution }}
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
`

func ExecutePromotion(promotionCtx *contexts.PromotionContext) error {
	tmpl, _ := template.New("PROMOTION").Parse(PromotionTpl)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, promotionCtx); err != nil {
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
		scriptsFoder := os.Getenv("NITRO_PIPELINES_SCRIPTS_FOLDER")
		if scriptsFoder != "" {
			exPath = scriptsFoder
		}
		fileName := fmt.Sprintf(exPath+"/nitro-%s-promote.sh", promotionCtx.Name)
		if err := saveToFile(fileName, buffer.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
