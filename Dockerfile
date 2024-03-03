FROM scratch

EXPOSE 1080
USER nobody

ADD .build/passwd /etc/

COPY rgosocks5 /

ENTRYPOINT ["/rgosocks5"]
