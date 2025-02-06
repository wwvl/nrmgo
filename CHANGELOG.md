# Changelog

<!-- This project adheres to [Semantic Versioning](https://semver.org/). -->

## 1.0.0

### Features

- Add support for bun package manager
- Add `backup` command for configuration file backup and restoration
- Add `config` command for configuration file initialization and viewing
- Add automatic retry mechanism to improve network request stability

### Enhancements

- Optimize registry latency testing algorithm
- Improve readability of error messages
- Optimize command execution concurrency control
- Improve configuration file read/write performance

### Bug Fixes

- Fix path parsing issues on Windows systems
- Fix memory leaks during concurrent requests

## 0.9.1

### Bug Fixes

- Fix configuration file permission check logic

## 0.9.0

### Features

- Add `info` command to view package manager information
  - Display installation status and version information
  - Show current registry and configuration file path
- Add `version` command to view version information

## 0.8.0

### Enhancements

- Improve command-line interaction experience
  - Add colored output support
  - Add progress bar display
  - Enhance error message clarity
  - Add confirmation prompts for dangerous operations

## 0.7.0

### Features

- Add `use` command for smart registry switching
  - Automatically detect and select fastest registry
  - Support direct switching to specified registry
  - Display switching progress bar
- Add `unuse` command to restore default registry

### Enhancements

- Optimize HTTP request performance using connection pool

## 0.6.0

### Features

- Add `add` command for custom registry management
  - Support setting registry name, URL, homepage, and description
  - Complete parameter validation and error prompts
- Add `rm` command for registry removal
  - Support single and batch deletion
  - Built-in registry protection mechanism
- Add `rename` command for renaming custom registries

## 0.5.0

### Features

- Add `ls` command to view registry list
  - Support basic and detailed information display
  - Highlight currently used registry

### Bug Fixes

- Fix configuration file read/write permission issues
- Fix registry URL format validation
- Fix command-line argument parsing errors

## 0.4.0

### Features

- Refactor code structure
- Add package manager detection
- Add registry detection

## 0.3.0

### Features

- Add configuration file support

### Enhancements

- Optimize HTTP requests using connection pool
- Improve concurrency control

## 0.2.1

### Enhancements

- Add built-in registry protection

### Bug Fixes

- Fix configuration file path issues
- Fix memory leaks

## 0.2.0

### Features

- Add support for npm, yarn, and pnpm package managers
- Implement basic registry management functionality
  - Support registry switching
  - Support registry list viewing
  - Support adding/removing/renaming custom registries

### Enhancements

- Improve command-line interface interaction
- Optimize registry switching speed
- Reduce memory usage

## 0.1.1

### Bug Fixes

- Fix registry switching failures
- Fix configuration file path parsing errors

## 0.1.0

### Features

- Implement basic registry switching functionality
- Support viewing built-in registry list
- Multi-platform support
