#####################
#       PROD        #
#####################

## BUILDER

ARG GO_VERSION

FROM golang:$GO_VERSION as builder

ARG USER
ARG UID

# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# create user - choose UID of host machine user that you want files to be owned by
RUN adduser $USER --no-create-home --disabled-password --disabled-login --uid $UID --quiet

# Copy app files
COPY . .
# Build app
RUN make bin/server && make bin/healthcheck

# PROD SERVER

FROM scratch as production

ARG USER

# Copy built binary from builder
COPY --from=builder /app/bin/server /bin/server
COPY --from=builder /app/bin/healthcheck /bin/healthcheck
COPY --from=builder /etc/passwd /etc/passwd

USER $USER

# Expose port
EXPOSE 80

# Exec built binary
HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD [ "/bin/healthcheck" ]
CMD [ "/bin/server" ]