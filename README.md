# Prometheus exporter for Citrix NetScaler
This exporter collects statistics from Citrix NetScaler and makes them available for Prometheus to pull.  As the NetScaler is an appliance it's not recommended to run the exporter directly on it, so it will need to run elsewhere.

## NetScaler configuration
The exporter works with any local NetScaler user account which has permissions to run the ``stats`` command.  It would be preferable to configure a specific user for this which only has permissions to retrieve stats.

If you lean towards the NetScaler CLI, you want to do something like the following (obviously changing the username as you see fit).

````
# Create a new Command Policy which is only allow to run the stat command
add system cmdPolicy stat ALLOW "^stat.*"

# Create a new user.  Disabling externalAuth is important as if it is enabled a user created in AD with the same name could login
add system user stats "password" -externalAuth DISABLED # Change the password to reflect whatever complex password you want

# Bind the local user account to the new Command Policy
bind system user stats stat 100
````

## Usage
You can monitor multiple NetScaler instances by passing in the URL, username, and password as command line flags to the exporter.  If you're running multiple exporters on the same server, you'll also need to change the port that the exporter binds to.

| Flag      | Description                                                                                               | Default Value |
| --------- | --------------------------------------------------------------------------------------------------------- | ------------- |
| url       | Base URL of the NetScaler management interface.  Normally something like https://mynetscaler.internal.com | none          |
| username  | Username with which to connect to the NetScaler API                                                       | none          |
| password  | Password with which to connect to the NetScaler API                                                       | none          |
| bind_port | Port to bind the exporter endpoint to                                                                     | 9280          |


Run the exporter manually using the following command:

````
Citrix-NetScaler-Exporter.exe -url https://mynetscaler.internal.com -username stats -password "my really strong password"
````

This will run the exporter using the default bind port.  If you need to change the port, append the ``-bind_port`` flag to the command.

### Running as a service
Ideally you'll run the exporter as a service.  There are many ways to do that, so it's really up to you.  If you're running it on Windows I would recommend [NSSM](https://nssm.cc/).

## Downloading a release
https://gitlab.com/rokett/Citrix-NetScaler-Exporter/tags

## Building the executable
All dependencies are version controlled, so building the project is really easy.

1. Clone the repository locally.
2. From within the repository directory run ``go build``.
3. Hey presto, you have an executable.
