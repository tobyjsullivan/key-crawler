FROM debian:stretch

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean
ADD ./build/queuer /var/opt/queuer
CMD /var/opt/queuer
