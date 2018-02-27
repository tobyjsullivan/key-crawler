FROM debian:stretch

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean
ADD ./build/recorder /var/opt/recorder
CMD /var/opt/recorder
