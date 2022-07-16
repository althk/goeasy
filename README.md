#

## GOEasy

A collection of utility packages aimed at reducing/eliminating setup toil for various workflows in golang.

### grpcutils

Provides handy functions that help setting up a new GRPC Server extremely simple.

* Simplified GRPC Server creation including:
  * TLS setup
  * Interceptors
  * Opentelemetry (basic metrics and tracing interceptors)

### sched

Provides a thread-safe task scheduler with an option to choose from different scheduling algorithms.
See [examples directory](sched/examples/) for examples.
