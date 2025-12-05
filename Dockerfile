
# SPDX-License-Identifier: Apache-2.0

########################################################################
##    docker build --no-cache --target certs -t vela-email:certs .    ##
########################################################################

FROM alpine:3.23.0@sha256:51183f2cfa6320055da30872f211093f9ff1d3cf06f39a0bdb212314c5dc7375 as certs

RUN apk add --update --no-cache ca-certificates

#########################################################
##    docker build --no-cache -t vela-email:local .    ##
#########################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-email /bin/vela-email

ENTRYPOINT [ "/bin/vela-email" ]
