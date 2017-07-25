# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.4.0] - 2017-07-25
### Changed
 - Authentication to the NetScaler now only happens once per scrape; the session cookie is saved and re-used in future requests.  When the scrape finishes, the session is disconnected.  Previously each API request was authenticated individually.

## [1.3.0] - 2017-07-21
### Added
 - Exporting Service Group metrics

## [1.2.0] - 2017-07-21
### Added
 - Added service state; 1 if service is up and 0 for any other state.

## [1.1.0] - 2017-07-20
### Added
 - Added model_id which represents, on the VPX line at least, the maximum licensed throughout of the appliance.  For example a VPX 1000 allows for 1000 MB/s.

## [1.0.0] - 2017-07-15
Initial release