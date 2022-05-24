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
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func buildScript(script string) string {
	if len(script) > 0 {
		return script
	}
	return "# No script was provided"
}

func buildFileName(fileName string) string {
	if len(fileName) > 0 {
		return fileName
	}
	return  uuid.New().String()
}

func buildImageName(imageName string, env string, repoIncludeEnv bool) string {
	if len(imageName) == 0 || !repoIncludeEnv {
		return imageName
	}
	if len(env) == 0 {
		return strings.ToLower(fmt.Sprintf("build-%s", imageName))
	}
	return strings.ToLower(fmt.Sprintf("%s-%s", env, imageName))
}
