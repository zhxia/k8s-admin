FROM scratch
COPY bin/* /
EXPOSE 9981
ENTRYPOINT ["/k8s-operator","--debug"]