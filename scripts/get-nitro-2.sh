#!/bin/bash
git clone https://github.com/NitroAgility/nitro-pipelines.git
pushd ./nitro-pipelines > /dev/null
git checkout v2.0
go install . && go build -o nitro 
mv ./nitro ../nitro
mv ./scripts/nitro-pipe ../nitro-pipe
popd > /dev/null
chmod +x ./nitro && chmod +x ./nitro-pipe
rm -rf ./nitro-pipelines