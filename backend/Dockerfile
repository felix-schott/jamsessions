#####################
#       PROD        #
#####################

## BUILDER

ARG GO_VERSION

FROM golang:$GO_VERSION AS builder

# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy app files
COPY . .
# Build app
RUN make bin/server && make bin/healthcheck

# PROD SERVER

FROM scratch AS production

# Copy built binary from builder
COPY --from=builder /app/bin/server /bin/server
COPY --from=builder /app/bin/healthcheck /bin/healthcheck
COPY --from=builder /etc/passwd /etc/passwd

# Expose port
EXPOSE 80

# Exec built binary
HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD [ "/bin/healthcheck" ]
CMD [ "/bin/server" ]