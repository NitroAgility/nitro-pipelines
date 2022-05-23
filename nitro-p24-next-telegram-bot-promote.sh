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
	echo $PIPE_DEV_ENV_FILE | base64 --decode >> "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.tmp && envsubst < "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.tmp > "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.env && rm "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.tmp
else
	echo $PIPE_DEV_ENV_FILE | base64 -di >> "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.tmp && envsubst < "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.tmp > "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.env && rm "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.tmp
fi
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
source "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.env && export $(cut -d= -f1 "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.env)
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
rm -f "$NITROBIN"20a645d9-0dbe-44aa-84d7-414bce88bbaf.env
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


# Pull docker images
docker pull $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/build-p24-next-telegram-bot:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Push docker images
aws configure set aws_access_key_id $NITRO_PIPELINES_VARIABLES_TARGET_AWS_ACCESS_KEY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws configure set aws_secret_access_key $NITRO_PIPELINES_VARIABLES_TARGET_AWS_SECRET_ACCESS
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr get-login-password --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
aws ecr create-repository --no-cli-pager --repository-name dev-p24-next-telegram-bot --region $NITRO_PIPELINES_VARIABLES_TARGET_AWS_REGION || true
docker tag $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/build-p24-next-telegram-bot:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/dev-p24-next-telegram-bot:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/dev-p24-next-telegram-bot:$NITRO_PIPELINES_BUILD_NUMBER
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker tag $NITRO_PIPELINES_VARIABLES_SOURCE_DOCKER_REGISTRY/build-p24-next-telegram-bot:$NITRO_PIPELINES_BUILD_NUMBER $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/dev-p24-next-telegram-bot:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
docker push $NITRO_PIPELINES_VARIABLES_TARGET_DOCKER_REGISTRY/dev-p24-next-telegram-bot:latest
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi

# Post promotion
#CODE: POST PROMOTION
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
# Cleaning expanded variables
rm -f "$NITROBIN"platform.env
# Post execution
#CODE: POST EXECUTION
exit_code=$? && if [ $exit_code -ne 0 ]; then exit $exit_code; fi
