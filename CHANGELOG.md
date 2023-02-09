# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [4.6.0] - 2023-02-09
### Changed
- #51 Ensure that idle HTTP client connections are closed after collecting metrics.
- Corrected the import path to reflect the actual repository.

## [4.5.1] - 2020-09-14
### Fixed
 - #43 Ensuring that the logging reflects a custom bind port.  Previously the logging was being setup in `init` which was before the flags are parsed.

## [4.5.0] - 2020-07-25
### Added
 - #41 Added the following AAA metrics.
   - aaa_auth_success
   - aaa_auth_fail
   - aaa_auth_only_http_success
   - aaa_auth_only_http_fail
   - aaa_current_ica_sessions
   - aaa_current_ica_only_connections

## [4.4.0] - 2020-07-02
### Added
 - #37 Added LBVS and GSLB VS state metric.

## [4.3.1] - 2020-07-02
### Fixed
 - #13 Error when trying to retrieve stats for servicegroups which have multiple members, where each member is the same server but on different ports.  We now include the port as a label to avoid duplicates.

## [4.3.0] - 2020-01-24
### Added
 - VPN Virtual Server (NetScaler Gateway) stats.

## [4.2.0] - 2019-12-31
### Changed
 - Minimising the number of API requests needed to retrieve Service Group member stats by getting all service group members for a specific service group in one go.  We still have to query for every single service group, but it cuts down the number of API requests by over 50% on the assumption that each service group has multiple members.
 - Refactored the filenames to better reflect the functions contained within them.

## [4.1.0] - 2019-12-30
### Changed
 - Debug logging is now hidden behind the `--debug` flag.

## [4.0.0] - 2019-07-19
### Added
 - Multi query endpoint.  You no longer need to run a single instance per NetScaler you want to monitor.  The exporter now supports Prometheus targets so only one instance is required.
 - Return a 500 status code, and error message, if the NetScaler target cannot be contacted.

### Changed
 - Refactored exporter code into separate package for maintainability reasons.

### Removed
 - Option to specify a NetScaler to scrape from the command line.  Using Prometheus scrape targets is now the only way.

## [3.2.0] - 2018-12-12
### Added
 - Collection of GSLB Service metrics.
 - Collection of GSLB Virtual Server metrics.
 - Collection of Content Switching Virtual Server metrics.
 - Make now also builds Linux x64 executables.

## [3.1.0] - 2018-05-02
### Added
 - Option to ignore certificate errors.  This can be useful when the NetScaler uses a self-signed certificate which will error without this flag.  It should be used sparingly, and only when you fully trust the endpoint.

## [3.0.0] - 2018-05-01
### Added
 - The following metrics replace their equivalent rate/s.
   - total_received_mb
   - total_transmit_mb
   - http_responses
   - http_requests
   - interfaces_received_bytes
   - interfaces_received_bytes
   - interfaces_received_packets
   - interfaces_transmitted_packets
   - interfaces_jumbo_packets_received
   - interfaces_jumbo_packets_transmitted
   - interfaces_error_packets_received

### Removed
 - The following metrics are not needed as Prometheus can calculate them from their equivalent totals.
   - transmit_mb_per_second
   - received_mb_per_second
   - http_requests_rate
   - http_responses_rate
   - interfaces_received_bytes_per_second
   - interfaces_received_bytes_per_second
   - interfaces_received_packets_per_second
   - interfaces_transmitted_packets_per_second
   - interfaces_jumbo_packets_received_per_second
   - interfaces_jumbo_packets_transmitted_per_second
   - interfaces_error_packets_received_per_second
   - virtual_servers_hits_rate
   - virtual_servers_requests_rate
   - virtual_servers_responses_rate
   - virtual_servers_request_bytes_rate
   - virtual_servers_reponse_bytes_rate
   - service_throughput_rate
   - service_request_rate
   - service_responses_rate
   - service_request_bytes_rate
   - service_response_bytes_rate
   - service_virtual_server_service_hits_rate
   - servicegroup_requests_rate
   - servicegroup_responses_rate
   - servicegroup_request_bytes_rate
   - servicegroup_response_bytes_rate

## [2.0.0] - 2017-10-10
### Changed
 - Log entries are no longer sent to a file.  Instead they are logged to stdout in logfmt format.

## [1.4.1] - 2017-08-22
### Fixed
 - NetScaler API bug meant that trying to retrieve stats from a service group member which used a wildcard port (65535 in API and CLI, * in GUI) resulted in error.  Skipping these members until the bug is resolved.

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
