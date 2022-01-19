
# Copyright (c) 2022 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

########################################################################
##    docker build --no-cache --target certs -t vela-email:certs .    ##
########################################################################

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

#########################################################
##    docker build --no-cache -t vela-email:local .    ##
#########################################################

FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY release/vela-email /bin/vela-email

ENTRYPOINT [ "/bin/vela-email" ]
