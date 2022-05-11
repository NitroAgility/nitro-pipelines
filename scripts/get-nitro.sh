#!/bin/bash
rm -rf ./.nitroci/.binaries
mkdir -p ./.nitroci/.binaries
cd ./.nitroci/.binaries
git clone https://github.com/NitroAgility/nitro-pipelines.git
cd ./nitro-pipelines
git checkout v2.0
go install . && go build -o nitro 
mv ./nitro ../nitro
mv ./scripts/nitro-pipe ../nitro-pipe
cd ..
chmod +x ./nitro && chmod +x ./nitro-pipe
rm -rf ./nitro-pipelines
cd ../../