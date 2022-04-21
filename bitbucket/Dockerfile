# syntax=docker/dockerfile:1
FROM golang:1.17 AS build
WORKDIR /
COPY . .
RUN go install .
RUN go build -o nitro

FROM amazon/aws-cli:latest  
WORKDIR /
COPY --from=build /nitro ./
COPY --from=build /bitbucket/pipe/nitro-pipe ./
COPY --from=build /bitbucket/pipe/nitro-pipe-expand-base64-env ./
COPY ./LICENSE /
# Set environment variables.
ENV HOME /
RUN yum update -y \
 && yum install -y https://dl.fedoraproject.org/pub/epel/epel-release-latest-7.noarch.rpm \
 && yum install -y \
            git \
            sudo \
            httpd-tools \
            openssl \
            openssl11 \
            tar \
            unzip \
            wget \
 && yum clean all \
 && rm -rf /var/cache/yum
# Set permissions on files
RUN chmod +x /nitro-pipe && chmod +x /nitro-pipe-expand-base64-env
# Install Kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x ./kubectl && mv ./kubectl /usr/local/bin
#Install Helm
RUN curl -sLo get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 \
 && chmod +x get_helm.sh \
 && ./get_helm.sh
#Set environment path
RUN  echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc
# Define default command.
CMD ["bash"]