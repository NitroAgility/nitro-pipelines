#!/bin/bash
git clone https://github.com/NitroAgility/nitro-pipelines.git
pushd ./nitro-pipelines > /dev/null
go install . && go build -o nitro 
mv ./nitro ../nitro
mv ./nitro ../nitro-pipe
popd > /dev/null
chmod +x ./nitro && chmod +x ./nitro-pipe