
# SPDX-License-Identifier: Apache-2.0

########################################################################
##    docker build --no-cache --target certs -t vela-email:certs .    ##
########################################################################

FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c as certs

RUN apk add --update --no-cache ca-certificates

#########################################################
##    docker build --no-cache -t vela-email:local .    ##
#########################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-email /bin/vela-email

ENTRYPOINT [ "/bin/vela-email" ]
