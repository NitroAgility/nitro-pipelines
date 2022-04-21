#!/bin/bash
aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
aws ecr create-repository --repository-name {{ .Buid.ImageName }} --region $NITRO_PIPELINES_TARGET_AWS_REGION || true
docker build -t {{ .Buid.ImageName }}:latest {{ .Buid.DockerArgs }} -f {{ .Buid.Dockerfile }} .
docker tag {{ .Buid.ImageName }}:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .Buid.ImageName }}:latest
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .Buid.ImageName }}:latest
docker tag {{ .Buid.ImageName }}:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .Buid.ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/{{ .Buid.ImageName }}:$NITRO_PIPELINES_BUILD_NUMBER