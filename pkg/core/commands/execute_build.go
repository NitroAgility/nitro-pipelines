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

const buildTpl = `#!/bin/bash
# Configure local files
export AWS_CONFIG_FILE=./aws_config
export AWS_SHARED_CREDENTIALS_FILE=./aws_credentials
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
[[ ! -f  $NITROBIN/{{ .Name }}.env ]] && exit 1
if [ -s $NITROBIN/{{ .Name }}.env ]; then
	source $NITROBIN/{{ .Name }}.env && export $(cut -d= -f1 $NITROBIN/{{ .Name }}.env)
else
	echo "File $NITROBIN/{{ .Name }}.env is empty"
fi
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm -f $NITROBIN/{{ .Name -}}.env
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
{{ end -}}
{{ end -}}
# Environment configuration
aws configure set aws_access_key_id $NITRO_PIPELINES_VARIABLES_TARGET_AWS_ACCESS_KEY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws configure set aws_secret_access_key $NITRO_PIPELINES_VARIABLES_TARGET_AWS_SECRET_ACCESS
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
echo $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION
echo $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY
aws ecr get-login-password --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr create-repository --no-cli-pager --repository-name {{ .ImageName }} --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION || true
# Docker build
docker build -t {{ .ImageName }}:latest {{ .DockerArgs }} -f {{ .Dockerfile }} .
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Docker push
docker tag {{ .ImageName }}:latest $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker tag {{ .ImageName }}:latest $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm -f ./aws_config
rm -f ./aws_credentials
# Cleaning expanded variables
{{ range .Expand -}}
{{ if eq .Type "file" -}}
rm -f $NITROBIN/{{ .Name -}}.env
{{ end -}}
{{ end -}}
`

func ExecuteBuild(buildCtx *contexts.BuildContext) error {
	tmpl, _ := template.New("BUILD").Parse(buildTpl)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, buildCtx); err != nil {
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
		fileName := fmt.Sprintf(exPath + "/nitro-%s-build.sh", buildCtx.Name)
		if err := saveToFile(fileName, buffer.Bytes()); err != nil {
			return err
		}
	}
	return nil
}
