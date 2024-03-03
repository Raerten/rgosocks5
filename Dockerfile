FROM scratch
USER nonroot
ENTRYPOINT ["/rgosocks5"]
COPY rgosocks5 /
