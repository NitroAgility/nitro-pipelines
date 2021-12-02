import os, sys, stat, yaml
from typing import List, Optional
from pydantic import BaseModel


class Microservice(BaseModel):
    name: str
    dockerfile: str


class ExpandItem(BaseModel):
    variable: str
    type: str


class Build(BaseModel):
    expand: List[ExpandItem]


class ExpandItem1(BaseModel):
    variable: str
    name: Optional[str] = None
    type: str


class Scripts(BaseModel):
    pre_execution: Optional[str] = None
    post_execution: Optional[str] = None


class Helm(BaseModel):
    parameters: str


class Default(BaseModel):
    expand: Optional[List[ExpandItem1]]
    scripts: Optional[Scripts]
    helm: Optional[Helm]


class Deployments(BaseModel):
    default: Default


class Model(BaseModel):
    microservices: List[Microservice]
    build: Build
    deployments: Optional[Deployments]


path = os.getenv('NITRO_PIPELINES_MICROSERVICES_PATH', './microservices.ym')

header_template='#!/bin/bash'

expand_variable_template_create='''echo $@env_var_name@ | base64 --decode >> @file_name@.tmp && envsubst < ./@file_name@.tmp > ./@file_name@.env && rm ./@file_name@.tmp'''
expand_variable_template_destroy='''rm -f ./@file_name@.env'''

load_file_template='''source ./@file_name@.env && export $(cut -d= -f1 ./@file_name@.env)'''

build_template="""aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
docker build -t @registry@:latest --build-arg JFROG_USERNAME=$NITRO_PIPELINES_FROG_USERNAME --build-arg JFROG_PASSWORD=$NITRO_PIPELINES_JFROG_PASSWORD -f @dockerfile@ .
docker tag @registry@:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/@registry@:latest
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/@registry@:latest
docker tag @registry@:latest $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/@registry@:$NITRO_PIPELINES_BUILD_NUMBER
docker push $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY/@registry@:$NITRO_PIPELINES_BUILD_NUMBER"""

deploy_template="""@pre_execution@

aws configure set aws_access_key_id $NITRO_PIPELINES_TARGET_AWS_ACCESS_KEY
aws configure set aws_secret_access_key $NITRO_PIPELINES_TARGET_AWS_SECRET_ACCESS
aws ecr get-login-password --region $NITRO_PIPELINES_TARGET_AWS_REGION | docker login --username AWS --password-stdin $NITRO_PIPELINES_TARGET_DOCKER_REGISTRY
aws eks --region $NITRO_PIPELINES_TARGET_AWS_REGION update-kubeconfig --name $NITRO_PIPELINES_TARGET_AWS_EKS_CLUSTER_NAME
helm upgrade --install $NITRO_PIPELINES_TARGET_HELM_NAMESPACE "$NITRO_PIPELINES_TARGET_HELM_CHART_SOURCE/chart/$NITRO_PIPELINES_TARGET_HELM_CHART_NAME" --set environment=@env@ --set infrastructure.domain=$NITRO_PIPELINES_DOMAIN --set infrastructure.docker_registry=$NITRO_PIPELINES_TARGET_DOCKER_REGISTRY --set app.tag=$NITRO_PIPELINES_BUILD_NUMBER @helm_parameters@ -n $NITRO_PIPELINES_TARGET_HELM_NAMESPACE
@post_execution@"""


if len(sys.argv) > 1:
    arg_operation = sys.argv[1] if sys.argv[1] == 'any' or sys.argv[1] == 'build' or sys.argv[1] == 'deploy' else 'any'
else:
    arg_operation = 'any'
    

if len(sys.argv) > 2:
    arg_ms_name=sys.argv[2]
else:
    arg_ms_name = None


def expand_create_content(model):
    expanded = []
    for expand in model:
        expanded_name = expand.variable
        expanded_name = expanded_name.replace('${ENV}', os.getenv('ENV', ''))
        file_name = expand.name if expand.name is not None else 'variables'
        expanded.append(expand_variable_template_create.replace('@env_var_name@', expanded_name).replace('@file_name@', file_name))
        if expand.type is not None and expand.type.lower() == 'environment':
            expanded.append(load_file_template.replace('@env_var_name@', expanded_name).replace('@file_name@', file_name))
    return expanded


def expand_destroy_content(model):
    expanded = []
    for expand in model:
        expanded_name = expand.variable
        expanded_name = expanded_name.replace('${ENV}', os.getenv('ENV', ''))
        file_name = expand.name if expand.name is not None else 'variables'
        expanded.append(expand_variable_template_destroy.replace('@env_var_name@', expanded_name).replace('@file_name@', file_name))
    return expanded


def create_file(expanded_create: str, expanded_destroy:str, template: str, script_type:str, script_name: str, substitutions:any):
    build_script = f'{header_template}\n{expanded_create}\n{template}\n{expanded_destroy}'
    build_script_content = build_script
    for (key, value) in substitutions:
        build_script_content = build_script_content.replace(key, value)
    build_script_content = build_script_content.replace('@env@', (os.getenv('ENV', '')).lower())
    script_file = f'./nitro-{script_name}-{script_type}.sh'
    build_script_file = open(script_file, "w")
    n = build_script_file.write(build_script_content)
    build_script_file.close()
    st = os.stat(script_file)
    os.chmod(script_file, st.st_mode | stat.S_IEXEC)


with open(path) as file:
    yaml_obj = yaml.load(file, Loader=yaml.FullLoader)
    model = Model(**yaml_obj)
    if arg_operation == 'any' or arg_operation == 'build':
        for microservice in model.microservices:
            ms_name = microservice.name
            if arg_ms_name is not None and arg_ms_name != ms_name:
                continue
            ms_dockerfile = microservice.dockerfile
            if model.build is not None:
                expanded_create = '\n'.join(expand_create_content(model.build.expand))
                expanded_destroy = '\n'.join(expand_destroy_content(model.build.expand))
                substitutions = []
                substitutions.append(('@registry@', ms_name))
                substitutions.append(('@dockerfile@', ms_dockerfile))
                create_file(expanded_create, expanded_destroy, build_template, 'build', ms_name, substitutions)
    if arg_operation == 'any' or arg_operation == 'deploy':
        expanded_create = ''
        if model.deployments is not None and model.deployments.default is not None and model.deployments.default.expand is not None:
            expanded_create = '\n'.join(expand_create_content(model.deployments.default.expand))
            expanded_destroy = '\n'.join(expand_destroy_content(model.deployments.default.expand))
        template = deploy_template
        pre_execution = ''
        post_execution = ''
        if model.deployments.default.scripts is not None:
            if model.deployments.default.scripts.pre_execution is not None:
                pre_execution = model.deployments.default.scripts.pre_execution
            if model.deployments.default.scripts.post_execution is not None:
                post_execution = model.deployments.default.scripts.post_execution
        template = template.replace('@pre_execution@', pre_execution).replace('@post_execution@', post_execution)
        if model.deployments.default.helm is not None and model.deployments.default.helm.parameters is not None:
            template = template.replace('@helm_parameters@', model.deployments.default.helm.parameters)
        substitutions = []
        create_file(expanded_create, expanded_destroy, template, 'deploy', 'microservices', substitutions)
