# Prometheus exporter for Citrix NetScaler
This exporter collects statistics from Citrix NetScaler and makes them available for Prometheus to pull.  As the NetScaler is an appliance it's not recommended to run the exporter directly on it, so it will need to run elsewhere.

## NetScaler configuration
The exporter works via a local NetScaler user account.  It would be preferable to configure a specific user for this which only has permissions to retrieve stats and specific configuration details.

If you lean towards the NetScaler CLI, you want to do something like the following (obviously changing the username as you see fit).

````
# Create a new Command Policy which is only allowed to run the stat command
add system cmdPolicy stat ALLOW (^stat.*|show ns license|show serviceGroup)

# Create a new user.  Disabling externalAuth is important as if it is enabled a user created in AD (or other external source) with the same name could login
add system user stats "password" -externalAuth DISABLED # Change the password to reflect whatever complex password you want

# Bind the local user account to the new Command Policy
bind system user stats stat 100
````

## Usage
You can monitor multiple NetScaler instances by passing in the URL, username, and password as command line flags to the exporter.  If you're running multiple exporters on the same server, you'll also need to change the port that the exporter binds to.

| Flag        | Description                                                                                               | Default Value |
| ----------- | --------------------------------------------------------------------------------------------------------- | ------------- |
| url         | Base URL of the NetScaler management interface.  Normally something like https://mynetscaler.internal.com | none          |
| username    | Username with which to connect to the NetScaler API                                                       | none          |
| password    | Password with which to connect to the NetScaler API                                                       | none          |
| bind_port   | Port to bind the exporter endpoint to                                                                     | 9280          |
| ignore-cert | Ignore certificate errors.  Should be used sparingly and only when you fully trust the endpoint           | false         |


Run the exporter manually using the following command:

````
Citrix-NetScaler-Exporter.exe -url https://mynetscaler.internal.com -username stats -password "my really strong password"
````

This will run the exporter using the default bind port.  If you need to change the port, append the ``-bind_port`` flag to the command.

### Running as a service
Ideally you'll run the exporter as a service.  There are many ways to do that, so it's really up to you.  If you're running it on Windows I would recommend [NSSM](https://nssm.cc/).

## Exported metrics
### NetScaler

| Metric                                 | Metric Type | Unit    |
| -------------------------------------- | ----------- | ------- |
| CPU usage                              | Gauge       | Percent |
| Memory usage                           | Gauge       | Percent |
| Management CPU usage                   | Gauge       | Percent |
| Packet engine CPU usage                | Gauge       | Percent |
| /flash partition usage                 | Gauge       | Percent |
| /var partition usage                   | Gauge       | Percent |
| Total received MB                      | Gauge       | MB      |
| Total transmitted MB                   | Gauge       | MB      |
| HTTP requests                          | Gauge       | None    |
| HTTP responses                         | Gauge       | None    |
| Current client connections             | Gauge       | None    |
| Current established client connections | Gauge       | None    |
| Current server connections             | Gauge       | None    |
| Current established server connections | Gauge       | None    |

### Interfaces
For each interface, the following metrics are retrieved.

| Metric                               | Metric Type | Unit  |
| ------------------------------------ | ----------- | ----- |
| Interface ID                         | N/A         | None  |
| Received bytes                       | Gauge       | Bytes |
| Transmitted bytes                    | Gauge       | Bytes |
| Received packets                     | Gauge       | None  |
| Transmitted packets                  | Gauge       | None  |
| Jumbo packets retrieved              | Gauge       | None  |
| Jumbo packets transmitted            | Gauge       | None  |
| Error packets received               | Gauge       | None  |
| Intrerface alias                     | N/A         | None  |

## Virtual Servers
For each virtual server, the following metrics are retrieved.

| Metric                     | Metric Type | Unit    |
| ---------------------------| ----------- | ------- |
| Name                       | N/A         | None    |
| Waiting requests           | Gauge       | None    |
| Health                     | Gauge       | Percent |
| Inactive services          | Gauge       | None    |
| Active services            | Gauge       | None    |
| Total hits                 | Gauge       | None    |
| Total requests             | Gauge       | None    |
| Total responses            | Gauge       | None    |
| Total request bytes        | Gauge       | Bytes   |
| Total response bytes       | Gauge       | Bytes   |
| Current client connections | Gauge       | None    |
| Current server connections | Gauge       | None    |

## Services
For each service, the following metrics are retrieved.

| Metric                         | Metric Type | Unit    |
| -------------------------------| ----------- | ------- |
| Name                           | N/A         | None    |
| Throughput                     | Counter     | MB      |
| Average time to first byte     | Gauge       | Seconds |
| State                          | Gauge       | None    |
| Total requests                 | Gauge       | None    |
| Total responses                | Gauge       | None    |
| Total request bytes            | Gauge       | Bytes   |
| Total response bytes           | Gauge       | Bytes   |
| Current client connections     | Gauge       | None    |
| Surge count                    | Gauge       | None    |
| Current server connections     | Gauge       | None    |
| Server established connections | Gauge       | None    |
| Current reuse pool             | Gauge       | None    |
| Max clients                    | Gauge       | None    |
| Current load                   | Gauge       | Percent |
| Service hits                   | Counter     | None    |
| Active transactions            | Gauge       | None    |

## Service Groups
For each service group member, the following metrics are retrieved.

| Metric                         | Metric Type | Unit    |
| -------------------------------| ----------- | ------- |
| Average time to first byte     | Gauge       | Seconds |
| State                          | Gauge       | None    |
| Total requests                 | Gauge       | None    |
| Total responses                | Gauge       | None    |
| Total request bytes            | Gauge       | Bytes   |
| Total response bytes           | Gauge       | Bytes   |
| Current client connections     | Gauge       | None    |
| Surge count                    | Gauge       | None    |
| Current server connections     | Gauge       | None    |
| Server established connections | Gauge       | None    |
| Current reuse pool             | Gauge       | None    |
| Max clients                    | Gauge       | None    |

## Licensing

| Metric                         | Metric Type | Unit    |
| -------------------------------| ----------- | ------- |
| Model ID                       | Gauge       | None    |

## Downloading a release
<https://github.com/rokett/Citrix-NetScaler-Exporter/releases>

## Building the executable
All dependencies are version controlled, so building the project is really easy.

1. ``go get github.com/rokett/citrix-netscaler-exporter``.
2. From within the repository directory run ``go build``.
3. Hey presto, you have an executable.

## Dockerfile
A Dockerfile has been setup to create the exporter using golang:alpine3.6

This Dockerfile will create a container that will set the entrypoint as /Citrix-Netscaler-Exporter so you can just pass in the command line options
mentioned about to the container without needing to call the executable
