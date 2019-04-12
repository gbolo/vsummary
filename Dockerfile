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
      mkdir -p /tmp/build/bin && \
      make all && cp -rp bin/vsummary-* /tmp/build/bin/ && \
      cp -rp testdata/sampleconfig/vsummary-*.yaml /tmp/build/ && \
      cp -rp www /tmp/build/

#
#  FINAL BASE CONTAINER --------------------------------------------------------
#

FROM  gbolo/baseos:alpine

# prepare env vars
ENV   VSUMMARY_SERVER_STATIC_FILES_DIR=/opt/vsummary/www/static \
      VSUMMARY_SERVER_TEMPLATES_DIR=/opt/vsummary/www/templates

# prepare homedir
RUN   mkdir -p /opt/vsummary

# Copy in from builder
COPY  --from=builder /tmp/build/bin/* /usr/local/bin/
COPY  --from=builder /tmp/build/vsummary-*.yaml /opt/vsummary/
COPY  --from=builder /tmp/build/www /opt/vsummary/www

# fix permissions
RUN   chmod -R 755 /opt/vsummary

# WORK DIR
WORKDIR /opt/vsummary

# Run as non-privileged user by default
USER  65534

# Inherit gbolo/baseos entrypoint and pass it this argument
CMD     ["vsummary-server", "-config", "/opt/vsummary/vsummary-config.yaml"]
