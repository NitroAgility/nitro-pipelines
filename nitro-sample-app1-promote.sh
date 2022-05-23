#!/bin/bash
# Configure local files
export AWS_CONFIG_FILE="$NITROBIN"aws_config
export AWS_SHARED_CREDENTIALS_FILE="$NITROBIN"aws_credentials
# Pre execution
#CODE: PRE EXECUTION
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Expanding variables
filename=$(uuidgen)
if [ $MACHINE_OS == "OSX" ]; then
	echo $PIPE_DEV_ENV_FILE | base64 --decode >> "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.tmp && envsubst < "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.tmp > "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.env && rm "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.tmp
else
	echo $PIPE_DEV_ENV_FILE | base64 -di >> "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.tmp && envsubst < "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.tmp > "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.env && rm "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.tmp
fi
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
source "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.env && export $(cut -d= -f1 "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.env)
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm -f "$NITROBIN"ca1a494d-e73a-4ff0-8108-cf8a35ecc03c.env
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
filename=$(uuidgen)
if [ $MACHINE_OS == "OSX" ]; then
	echo $PIPE_DEV_PLATFORM_ENV_FILE | base64 --decode >> "$NITROBIN"platform.tmp && envsubst < "$NITROBIN"platform.tmp > "$NITROBIN"platform.env && rm "$NITROBIN"platform.tmp
else
	echo $PIPE_DEV_PLATFORM_ENV_FILE | base64 -di >> "$NITROBIN"platform.tmp && envsubst < "$NITROBIN"platform.tmp > "$NITROBIN"platform.env && rm "$NITROBIN"platform.tmp
fi
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Environment configuration
aws configure set aws_access_key_id $NITRO_PIPELINES_VARIABLES_SOURCE_AWS_ACCESS_KEY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws configure set aws_secret_access_key $NITRO_PIPELINES_VARIABLES_SOURCE_AWS_SECRET_ACCESS
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr get-login-password --region $NITRO_PIPELINES_VARIABLES_SOURCE_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Pre promotion
#CODE: PRE PROMOTION
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
#Retagging an image
MANIFEST=$(aws ecr batch-get-image --repository-name $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/build-sample-app1 --image-ids imageTag=$NITRO_PIPELINES_BUILD_NUMBER --output json | jq --raw-output --join-output '.images[0].imageManifest')
aws ecr put-image --repository-name $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/dev-sample-app1 --image-tag $NITRO_PIPELINES_BUILD_NUMBER --image-manifest "$MANIFEST"


# Post promotion
#CODE: POST PROMOTION
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Cleaning expanded variables
rm -f "$NITROBIN"platform.env
# Post execution
#CODE: POST EXECUTION
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
