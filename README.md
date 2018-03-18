# oauth2_proxy-token
Provides an oauth2-proxy upstream token generator server for basic-auth clients

## Purpose
Placing online services behind an SSO proxy, such as the oauth2_proxy, helps protect those services from unauthorized users. However, a drawback exists wherein requests to those services outside of a browser are not easily accommodated. This oauth2_proxy-token application solves this problem.

## Tokens
When a user is successfully authenticated to the SSO proxy the user can the navigate to this upstream application via a URL such as http://token.dev.local/. Upon detecting the incoming request, this token application will generate a new secret token and associate it to the user. The user is known because the SSO proxy will be configured to push the username through to this application via an HTTP header. The user and the new secret token are then placed into an htpasswd file, along with a comment containing the expiration date of the token.

## Expired Tokens
An asynchronous maintenance thread will routinely review the htpasswd tokens. Any token that has expired will cause that htpasswd line to be commented out. By commenting the line out rather than deleting it, the token application can later uncomment out the token when the user successfully navigates back to this token app in the browser.

## htpasswd Synchronization
The SSO proxy application is configured to validate basic authentication requests with the htpasswd file. The proxy application must either periodically refresh the htpasswd file in memory, or restart on a frequent basis to keep the basic auth tokens synchronized with this application.

## New Tokens
Occassionally, users will lose their tokens and need to generate new ones. To do this, the user can add a special query parameter `?new` onto the URL, which will forcefully generate a new token and update the htpasswd file. Ex: http://token.dev.local/?new

## Kubernetes
This application was originally designed to run inside of a Kubernetes pod, as a sidecar container to oauth2_proxy. This allows the two containers to share a volume, and thus the htpasswd file is easily exposed to both applications. Because oauth2_proxy does not refresh the htpasswd file from disk, it's necessary to forcibly restart oauth2_proxy on a periodic basis. This can be accomplished via the `timeout` shell command. However, take care to not cause the entire pod to be recycled, otherwise there will be a noticeable downtime of the SSO proxy to users. Rather, place the timeout in a shell while loop that only stops looping when the oauth2_proxy app exits on its own. See the https://github.com/jertel/oauth2_proxy-docker project for a working example. 

## Configuration
A JSON configuration file is required to be provided for this application to start. Below is an example file, with an explanation following:

```
{
  "header.uri": "X-Original-URI",
  "header.username": "X-Auth-Request-User",
  "http.hostport": ":8080",
  "http.path": "/",
  "htpasswd.filepath": "/etc/ssotoken/htpasswd",
  "maintenance.intervalSecs": 60,
  "token.lengthBytes": 32,
  "token.durationHours": 24
}
```

Property                  | Description
---------                 |------------
header.uri                | The HTTP header containing the URI that brought the request to this application
header.username           | The HTTP header containing the authenticated username that submitted the request
http.hostport             | host (or IP address) and port on which to listen for incoming requests (host/IP are optional)
http.path                 | The endpoint URL to monitor for incoming requests
htpasswd.filepath         | Path where the htpasswd file will be created, or opened if already exists
maintenance.intervalSecs  | Number of seconds to wait before checking for expired tokens
token.lengthBytes         | Number of bytes to use when generating the token (before base64 encoding)
token.durationHours       | Number of hours to allow a token to be used before it is expired