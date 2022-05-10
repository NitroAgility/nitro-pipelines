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
	"strings"
	"text/template"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/contexts"
)

const buildTpl = `#!/bin/bash
echo step 1
# Expanding variables
{{ range .Expand -}}
echo ${{ .Variable }} | base64 --decode >> ./{{ .Name }}.tmp && envsubst < ./{{ .Name }}.tmp > ./{{ .Name }}.env && rm ./{{ .Name }}.tmp
echo step 2
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
echo step 3
{{ if eq .Type "environment" -}}
[[ ! -f  ./{{ .Name }}.env ]] && exit 1
if [ -s diff.txt ]; then
	source ./{{ .Name }}.env && export $(cut -d= -f1 ./{{ .Name }}.env)
else
	echo "File ./{{ .Name }}.env is empty"
fi
echo step 4
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
echo step 5
rm -f ./{{ .Name -}}.env
echo step 6
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
echo step 7
{{ end -}}
{{ end -}}
echo step 8
# Environment configuration
aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr create-repository --repository-name {{ .ImageName }} --region $NITRO_PIPELINES_TARGET_AWS_REGION || true
# Docker build
docker build -t {{ .ImageName }}:latest {{ .DockerArgs }} -f {{ .Dockerfile }} .
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Docker push
docker tag {{ .ImageName }}:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker tag {{ .ImageName }}:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Cleaning expanded variables
{{ range .Expand -}}
{{ if eq .Type "file" -}}
rm -f ./{{ .Name -}}.env
{{ end -}}
{{ end -}}
`

func ExecuteBuild(buildCtx *contexts.BuildContext) (error) {
    tmpl, _ :=  template.New("BUILD").Parse(buildTpl)
	var buffer bytes.Buffer
	if err := tmpl.Execute(&buffer, buildCtx); err != nil {
        fmt.Print(buffer.String())
		return err
	}
	if strings.ToUpper(os.Getenv("DRY_RUN")) == "TRUE" {
		fmt.Println(buffer.String())
	} else {
		fileName := fmt.Sprintf("./nitro-%s-build.sh",buildCtx.Name)
		if err := saveToFile(fileName, buffer.Bytes()); err != nil {
			return err
		}
	}
    return nil
}
