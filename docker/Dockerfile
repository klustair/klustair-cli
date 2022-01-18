FROM alpine:3.15
ARG VERSION=0.0.1

RUN wget https://github.com/klustair/klustair-cli/releases/download/v${VERSION}/klustair-cli_${VERSION}_Linux_i386.tar.gz

RUN tar -xvzf klustair-cli_${VERSION}_Linux_i386.tar.gz && rm klustair-cli_${VERSION}_Linux_i386.tar.gz

RUN mv klustair-cli /usr/local/bin/klustair-cli