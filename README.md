<h1 align="center">
      LogSnare
</h1>

<h4 align="center">A web application playground for testing IDOR, broken access controls, and logging in Go.</h4>

<img src="https://raw.githubusercontent.com/sea-erkin/log-snare/main/web/ui/assets/img/logsnare-gopher.png">

## Overview 

LogSnare is an intentionally vulnerable web application, where your goal is to go from
a basic `gopher` user of the LogSnare company, to the prestigious `acme-admin` of Acme Corporation.

The application, while hosting multiple vulnerabilities, serves as a valuable educational tool. However, the real lesson to be learned here is how to 
**prevent and catch these attacks leveraging proper validation and logging**. 

After logging in to the demo application, in the top navbar you'll see a validation toggle
which allows you to toggle security controls in real-time.

<img src="https://raw.githubusercontent.com/sea-erkin/log-snare/main/web/ui/assets/img/challenge.jpg">

## Getting Started

The easiest way to get started is with docker.
```
docker pull seaerkin/log-snare
docker run -p 127.0.0.1:8080:8080 seaerkin/log-snare
```
You'll receive a username and password to login, have at it from there!

## Catching Attackers with Logging
Insecure Direct Object References (IDORs) fall under the OWASP Top category of "Broken Access Controls". 
These vulnerabilities are some of the most severe as they can allow end-users access to resources they shouldn't be able to access. 

Most people don't realize, IDOR vulnerabilities are some of the best opportunities for logging. This is because in most cases, a user will never
"accidentally" trigger an IDOR. IDOR vulnerabilities are typically achieved when a user asks for resources outside their allowed interface.

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