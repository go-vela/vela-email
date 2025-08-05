
# SPDX-License-Identifier: Apache-2.0

########################################################################
##    docker build --no-cache --target certs -t vela-email:certs .    ##
########################################################################

FROM alpine:3.22.1@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1 as certs

RUN apk add --update --no-cache ca-certificates

#########################################################
##    docker build --no-cache -t vela-email:local .    ##
#########################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-email /bin/vela-email

ENTRYPOINT [ "/bin/vela-email" ]
