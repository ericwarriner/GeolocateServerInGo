# Copyright 2022 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19.4 as builder
WORKDIR /app

# Initialize a new Go module.
RUN go mod init geolocate
RUN go get github.com/gin-gonic/gin
RUN go get github.com/oschwald/geoip2-golang
RUN go get github.com/gin-contrib/cors

# Copy local code to the container image.
COPY ./cmd/geolocate-ip/*.go ./

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -o /geolocate

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM gcr.io/distroless/static-debian11

# Change the working directory.
WORKDIR /

# Copy the binary to the production image from the builder stage.
COPY --from=builder /geolocate /geolocate
COPY ./cmd/geolocate-ip/tools/* /cmd/geolocate-ip/tools/
COPY ./files/* /cmd/geolocate-ip/files/

EXPOSE 8080
# Run the web service on container startup.
USER nonroot:nonroot
ENTRYPOINT ["/geolocate"]