# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## 0.2.1 - 2018-11-20
### Fixed
- improved error message when you try to launch a workflow without having listened yet
- changes cp command in circle ci
- fixed wait weekday and DayOfMonth methods to wait for the next time if the given time is today.

### Added
- add tests for wait methods

## 0.2.0 - 2018-10-14
### Added
- add integration tests to circle ci.

## 0.2.0 - 2018-10-10
### Added
- unit tests for worklow package

### Changed
- Changed unsafe exported names to indicate that they are not safe

## 0.1.0 - 2018-10-10
### Added
- CircleCI
- Changelog.
