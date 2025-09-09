# Stage: Build
FROM golang:1.25.0 AS builder

ENV OUT=/out
# Fully static binaries (for better portability)
ENV CGO_ENABLED=0

# Prepare image
RUN apt-get update && apt-get install -y upx && mkdir -p $OUT

WORKDIR /src

COPY . .

RUN go mod tidy \
    && ./scripts/build.sh $OUT

# Stage: Export binaries
FROM scratch AS export
COPY --from=builder /out /out/
# Neede for later be able to create and image from the build and copy the files we want within it
CMD ["true"]
