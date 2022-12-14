FROM gcr.io/buildpacks/gcp/run:v1
USER root
RUN apt-get update -y && \
    apt-get install apt-transport-https ca-certificates gnupg -y && \
    echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] http://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg  add - && \
    apt-get update -y && \
    apt-get install google-cloud-sdk -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/*
USER cnb