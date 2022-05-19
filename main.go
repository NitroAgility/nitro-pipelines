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
	"fmt"
	"os"
	"strings"

	"github.com/NitroAgility/nitro-pipelines/pkg/core/commands"
	"github.com/NitroAgility/nitro-pipelines/pkg/core/contexts"
)

func main() {
	msFilePath := os.Getenv("NITRO_PIPELINES_MICROSERVICES_PATH")
	if len(msFilePath) == 0 {
		msFilePath = "microservices.yml"
	}
	if len(os.Args) > 1 {
		if strings.ToUpper(os.Args[1]) == "BUILD" {
			if len(os.Args) > 2 {
				buildCtx, err := contexts.NewBuildContext(msFilePath, os.Args[2])
				if err != nil {
					fmt.Println("An error has occurred whilst executing the command, ", err)
					os.Exit(1)
				}
				commands.ExecuteBuild(buildCtx)
				return
			}
		} else if strings.ToUpper(os.Args[1]) == "PROMOTION" || strings.ToUpper(os.Args[1]) == "DEPLOY" { 
			deployCtx, err := contexts.NewDeployContext(msFilePath)
			if err != nil {
				fmt.Println("An error has occurred whilst executing the command, ", err)
				os.Exit(1)
			}
			if strings.ToUpper(os.Args[1]) == "PROMOTION" {
				commands.ExecutePromotion(deployCtx)
			} else {
				commands.ExecuteDeploy(deployCtx)
			}
			return
		}	
	}
	fmt.Println("Invalid command.")
	os.Exit(1)
}
