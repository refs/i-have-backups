# Design

Included inside there is:

- 1x HTTP Service
- 1x gRPC Service

gRPC handles internal communication, HTTP is user facing and rely calls to the gRPC service.

# TODO

- [ ] configure jaeger exporter