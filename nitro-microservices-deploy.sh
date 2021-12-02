#!/bin/bash
echo $PIPE_DEV_ENV_FILE | base64 --decode >> variables.tmp && envsubst < ./variables.tmp > ./variables.env && rm ./variables.tmp
source ./variables.env && export $(cut -d= -f1 ./variables.env)
echo $PIPE_DEV_PLATFORM_ENV_FILE | base64 --decode >> platform.tmp && envsubst < ./platform.tmp > ./platform.env && rm ./platform.tmp
echo $PIPE_DEV_PROTOCOLS_ENV_FILE | base64 --decode >> protocols.tmp && envsubst < ./protocols.tmp > ./protocols.env && rm ./protocols.tmp
echo "Pipeline is about to start..."
aws configure set aws_access_key_id $NITRO_PIPELINES_SOURCE_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_SOURCE_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_SOURCE_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_SOURCE_DOCKER_REGISTRY
docker pull $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/groups-microservice:$NITRO_PIPELINES_BUILD_NUMBER
docker tag $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/groups-microservice:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/groups-microservice-${ENV_TARGET}:$NITRO_PIPELINES_BUILD_NUMBER
docker pull $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/users-microservice:$NITRO_PIPELINES_BUILD_NUMBER
docker tag $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/users-microservice:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/users-microservice-${ENV_TARGET}:$NITRO_PIPELINES_BUILD_NUMBER
aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/groups-microservice-${ENV_TARGET}:$NITRO_PIPELINES_BUILD_NUMBER
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/users-microservice-${ENV_TARGET}:$NITRO_PIPELINES_BUILD_NUMBER
aws eks --region $NITRO_PIPELINES_TARGET_AWS_REGION update-kubeconfig --name $NITRO_PIPELINES_TARGET_AWS_EKS_CLUSTER_NAME
helm upgrade --install $NITRO_PIPELINES_TARGET_HELM_NAMESPACE "$NITRO_PIPELINES_TARGET_HELM_CHART_SOURCE/chart/$NITRO_PIPELINES_TARGET_HELM_CHART_NAME" --set environment=dev --set infrastructure.domain=$NITRO_PIPELINES_DOMAIN --set infrastructure.docker_registry=$NITRO_PIPELINES_TARGET_DOCKER_REGISTRY --set app.tag=$NITRO_PIPELINES_BUILD_NUMBER --set env.platform="$(base64 -w 0 ./platform.env)" --set env.protocols="$(base64 -w 0 ./protocols.env)" -n $NITRO_PIPELINES_TARGET_HELM_NAMESPACE
echo "Pipeline completed..."
rm -f ./variables.env
rm -f ./platform.env
rm -f ./protocols.env