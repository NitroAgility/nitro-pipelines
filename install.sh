#!/bin/bash
git clone https://github.com/NitroAgility/nitro-pipelines.git
cd ./nitro-pipelines && go install . && go build -o nitro && mv ./nitro ../nitro
curl https://raw.githubusercontent.com/NitroAgility/nitro-pipelines/main/bitbucket/pipe/nitro-pipe -o ./nitro-pipe
chmod +x ./nitro && chmod +x ./nitro-pipe