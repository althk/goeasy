#

## GOEasy

A collection of utility modules aimed at reducing/eliminating setup toil for various workflows in golang.

The primary goals of these modules are to provide:

* Simple APIs
* Declarative paradigm where possible
* Easily pluggable with standard API

Using these modules will result in a simpler, cleaner codebase with little to no learning curve.

### grpcutils

Provides handy functions that help setting up a new GRPC Server extremely simple.

* Simplified GRPC Server creation including:
  * TLS setup
  * Interceptors
  * Opentelemetry (basic metrics and tracing interceptors)

### sched

Provides a thread-safe task scheduler with an option to choose from different scheduling algorithms.
See [examples directory](sched/examples/) for examples.
