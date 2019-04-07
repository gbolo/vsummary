# github.com/gbolo/vsummary

#
#  BUILD CONTAINER -------------------------------------------------------------
#

FROM gbolo/builder:alpine as builder

COPY . ${GOPATH}/src/github.com/gbolo/vsummary

# Building
RUN   set -xe; \
      SRC_DIR=${GOPATH}/src/github.com/gbolo/vsummary; \
      cd ${SRC_DIR}; \
      mkdir -p /tmp/build && \
      make server && cp -rp bin/vsummary-server /tmp/build/ && \
      cp -rp testdata/sampleconfig/vsummary-config.yaml /tmp/build/

#
#  FINAL BASE CONTAINER --------------------------------------------------------
#

FROM  gbolo/baseos:alpine

# prepare env vars

# prepare homedir
RUN   mkdir -p /opt/vsummary

# Copy in from builder
COPY  --from=builder /tmp/build/ /opt/vsummary

# WORK DIR
WORKDIR /opt/vsummary

# Run as non-privileged user by default
USER  65534

# Inherit gbolo/baseos entrypoint and pass it this argument
CMD     ["/opt/vsummary/vsummary-server"]

