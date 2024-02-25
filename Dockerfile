FROM alpine:3.19
COPY kubelab /usr/bin/kubelab
ENTRYPOINT ["/usr/bin/kubelab"]