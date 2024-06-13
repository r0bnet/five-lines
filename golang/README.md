Five Lines of Code in Go
========================

## Introduction

This folder contains the [`Five Lines of Code`](https://github.com/thedrlambda/five-lines) project migrated to Go. The following sections describe the differences that were introduced in the Go version.

One of the main differences is that the game is not running in the browser or having an user interface at all. This project should rather highlight the `Five Lines of Code` refactorings in Go instead of being a playable game. However it should mimic the input and game logic as much as possible so that it's theoretically possible to run the game in the future.

One of the main motivations is to explore whether the refactoring rules that are meant to be applied on object-oriented languages (like TypeScript) are valid in Go as well, although it's [not strictly object-oriented](https://go.dev/doc/faq#Is_Go_an_object-oriented_language).

## Differences

### Section 3.1

This is the inital state of the project and some changes were needed so that the project is not only running Go code, but also following the Go conventions. The following changes were made (from top to bottom):

- **imports**:

  The TypeScript version is meant to be executed in the browser and doesn't need any imports, since all functions are part of the global scope. In Go, however, we need to import some packages:
  - `time`: to calculate the sleep time for the game loop
  - `sync`: to synchronize the game loop and the input handling. We need to use a `sync.Mutex` to synchronize access to the global `inputs` slice. The game loop runs in a separate goroutine. JavaScript is single-threaded, so we don't need to worry about synchronization there.
  - `"github.com/eiannone/keyboard"` to handle keyboard inputs. This third-party package offers a simple abstraction over the keyboard input handling, which comes out of the box in JavaScript.
- **enums**:

  Go does not have enums like TypeScript. The closest thing to an enum in Go is a set of `const` declarations, whereas a group of constants is prefixed with a common identifier. Besides that we can make use of the `iota` keyword to create a sequence of related constants. This is default behavior for TypeScript enums.
- **entrypoint**:

  Since TypeScript is transpiled to JavaScript it doesn't need an explicit entrypoint function. It will execute anything that is in the global scope. Go, however, needs an explicit entrypoint function, which is the `main` function. It starts the game loop in a goroutine and handles the input in the main thread. The game loop is in a goroutine because it's an endless loop that would otherwise block the execution. This would prevent us from capturing any keyboard input.

- **stubs**:

  Since there currently is no user interface (like the browser with some HTML elements) we have to stub some global objects (`Canvas` and `CanvasRenderingContext2D`) and write helper methods to mimic their behavior. Luckily two of the four methods are simple procedures with no return value. The other two methods could easily be stubbed.

All the rest of the code is pretty much the same as in the TypeScript version. Even the syntax is really similar and should be easily readable for anyone who is not familiar with Go.
