# goSDK

This SDK provides utilities and modules that are useful in building golang applications.

Utilities and modules in this SDK should be reusable and their goal is to reduce toil in the other projects that use them.

## Definitions

**Utilities:** a utility is a collection of functions, type definitions, or interfaces within a common domain where each function, type, or interface reduces the amount of code that needs to be written to perform a common task within that domain.  Utilities can be found under the `/util` directory.

**Modules:** a module wraps complex, but reusable functionality and provides a simplified interface for accessing it.  Modules can be found under the `/modules` directory

## Features

The SDK supports the following features:

### Modules

* EventBus

### Utilities

* arrays
* bits
* events
* interfaces
* json
* maps
* stats
* thread

## Testing

Standard golang test harness using `go test`