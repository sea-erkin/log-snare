<h1 align="center">
      LogSnare
</h1>

<h4 align="center">A web application playground for testing IDOR, broken access controls, and logging in Golang.</h4>

## Overview 

LogSnare is an intentionally vulnerable web application, where your goal is to go from
a basic `gopher` user of the LogSnare company, to the prestigious `acme-admin` of Acme Corporation.

The application contains several vulnerabilities, however, the real lesson to be learned here is how to 
**prevent and catch these attacks leveraging proper validation and logging**. In the top navbar of the application you'll see a validation toggle
which allows users to enable and disable server-side validation, seeing how the application would react when it is vulnerable,
and when access to objects is validated and secured.

<img src="https://raw.githubusercontent.com/sea-erkin/log-snare/main/web/ui/assets/img/challenge.jpg">

## Getting Started

The easiest way to get started is with docker.
```
docker pull seaerkin/log-snare
docker run -p 127.0.0.1:8080:8080 seaerkin/log-snare
```
You'll receive a username and password to login, have at it from there!

## Catching Attackers with Logging
Insecure Direct Object References (IDORs) fall under the OWASP Top 10 category of "Broken Access Control". 
These vulnerabilities are some of the most severe as they can allow end-users access to resources they shouldn't be able to access. 

What most people don't realize, is IDOR vulnerabilities are some of the best opportunities for logging because in most cases, a user will never
"accidentally" trigger an IDOR. IDOR vulnerabilities are typically achieved when a user asks for resources outside of their allowed interface.

Attackers abuse web applications by asking web servers to return resources the user may not have access too, and this application hopes
to serve as an educational resources on how to fix, prevent, and log these types of security events.

Here are some example log outputs from the application when validation is enabled.

```
{"message":"user is trying to access a company ID that is not theirs","program":"log-snare","version":0.1,"username":"gopher","eventType":"security","securityType":"tamper-certain","eventCategory":"validation","clientIp":"172.17.0.1"}
{"message":"user is trying to enable admin, but they are a basic user","program":"log-snare","version":0.1,"username":"gopher","eventType":"security","securityType":"tamper-certain","eventCategory":"validation","clientIp":"172.17.0.1"}
```

## Checking Application Logs

All logs print to stdout by default, however if you want to view just the application logs containing validation logic, you can do the following:
```
docker exec -it <container_id> /bin/bash
tail -f logsnare.log
```
