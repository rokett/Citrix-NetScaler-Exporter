# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

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
