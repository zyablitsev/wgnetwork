FROM busybox AS build-env
RUN mkdir -p -m 1755 /foo/tmp

FROM scratch
ADD bin/wgn-managercli_linux_amd64 /wgn_managercli
ADD bin/wgnetwork_linux_amd64 /wgnetwork
COPY --from=build-env /tmp /tmp

CMD ["/wgnetwork"]
