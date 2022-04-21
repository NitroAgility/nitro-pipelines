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

package main

import (
	"log"
	"os"
	"strings"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/commands"
	"github.com/NitroAgility/nitro-pipelines/pkg/core/contexts"
	"github.com/NitroAgility/nitro-pipelines/pkg/core/models"
)

func main() {
	context := contexts.NewContext()
	msModel, _ := models.LoadMicroservicesFile("./test/microservices.yml", context)
	if len(os.Args) > 1 {
		if strings.ToUpper(os.Args[1]) == "BUILD" {
			if len(os.Args) < 3 {
				log.Fatal("Invalid command.")
				os.Exit(1)
			}
			for _, n := range msModel.Microservices {
				if n.Name == os.Args[2]{
					buildCtx := contexts.NewBuildContext(n.Dockerfile, msModel.Build.BuildArgs, n.Name)
					commands.ExecuteBuild(buildCtx)
					return
				}
			}
		} else if strings.ToUpper(os.Args[1]) == "DEPLOY" { 
			deployCtx := contexts.NewDeployContext("XYZ", msModel.Build.BuildArgs, "ABC")
			commands.ExecuteDeploy(deployCtx)
			return
		}	
	}
	log.Fatal("Invalid command.")
	os.Exit(1)
}
