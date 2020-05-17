FROM registry.access.redhat.com/ubi8/ubi-init:latest

LABEL maintainer="lzuccarelli@tfd.ie"

RUN dnf remove -y subscription-manager

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
COPY bin/* uid_entrypoint.sh build/microservice /go/

RUN rpm -iv /go/librdkafka-0.11.5-1.el8.x86_64.rpm && rm -rf /go/lib*

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 0755 "$GOPATH"
WORKDIR $GOPATH

USER 1001

ENTRYPOINT [ "./uid_entrypoint.sh" ]

# This will change depending on each microservice entry point
CMD ["./microservice"]
