FROM debian:stretch

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean
ADD ./build/enum-batch-gen /var/opt/enum-batch-gen
CMD /var/opt/enum-batch-gen
