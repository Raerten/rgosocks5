FROM scratch

ARG TARGETPLATFORM

EXPOSE 1080
USER nobody

ADD .build/passwd /etc/

COPY $TARGETPLATFORM/rgosocks5 /

ENTRYPOINT ["/rgosocks5"]
