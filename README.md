# nitro-pipelines

Enhanced pipelines which include the Nitro tools.

## How to build

```bash
tag=83 && docker build -t nitroagility/nitro-bitbucket-pipelines:$tag -f ./bitbucket/Dockerfile . && docker push nitroagility/nitro-bitbucket-pipelines:$tag
```
