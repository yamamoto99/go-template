FROM mcr.microsoft.com/devcontainers/go:2-1.26-trixie

ARG MIGRATE_VERSION=v4.17.1

RUN ARCH="$(dpkg --print-architecture)" \
 && curl -fsSL "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-${ARCH}.tar.gz" \
    | tar -xz -C /usr/local/bin migrate \
 && chmod +x /usr/local/bin/migrate \
 && migrate -version
