FROM alpine:3.20
COPY kubelab /usr/bin/kubelab
ENTRYPOINT ["/usr/bin/kubelab"]