FROM scratch

ARG VERSION
LABEL maintainer="Ben Cessa <ben@aid.technology>"
LABEL version=${VERSION}

COPY suss-workshop-linux /
COPY ca-roots.crt /etc/ssl/certs/

EXPOSE 9090

ENTRYPOINT ["/suss-workshop-linux"]
