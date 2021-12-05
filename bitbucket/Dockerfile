FROM ubuntu

COPY ./requirements.txt /requirements.txt

COPY ./bitbucket/pipe /
COPY ./LICENSE /

# Set environment variables.
ENV HOME /root

# Define working directory.
WORKDIR /root

# Set permissions on files
RUN chmod +x /nitro-pipe && chmod +x /nitro-pipe-expand-base64-env

# Install required tools
RUN apt-get update && apt install -y software-properties-common && add-apt-repository ppa:deadsnakes/ppa 
RUN apt-get update && apt-get install -y unzip \
    curl \
    uuid-runtime \
    gettext-base \
    python3.9 \
    python3-pip \
    python3.8-venv

# Install Dependencies
RUN pip install -r /requirements.txt && rm /requirements.txt

# Install AWS Cli
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && unzip awscliv2.zip && ./aws/install && rm awscliv2.zip && rm -rf awscliv2
RUN curl -o aws-iam-authenticator https://amazon-eks.s3.us-west-2.amazonaws.com/1.16.8/2020-04-16/bin/linux/amd64/aws-iam-authenticator && chmod +x ./aws-iam-authenticator && mkdir -p $HOME/bin && cp ./aws-iam-authenticator $HOME/bin/aws-iam-authenticator && export PATH=$PATH:$HOME/bin && rm -rf ./aws-iam-authenticator

# Install Kubectl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && chmod +x ./kubectl && mv ./kubectl /usr/local/bin

#Install Helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 && chmod 700 get_helm.sh && ./get_helm.sh && rm ./get_helm.sh 

#Set environment path
RUN  echo 'export PATH=$PATH:$HOME/bin' >> ~/.bashrc

# Define default command.
CMD ["bash"]