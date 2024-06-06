# AdGoBye Blocklist Server
This repository contains the AdGoBye Blocklist callback server. Clients can opt into sending back blocklist misses in
order to notify maintainers that blocklists are out of date.

# Running
The Blocklist callback server is a handler to ship the results off to a database to be analyzed and turned
into statistics. [The Docker compose file](docker/docker-compose.yml) contains the deployment for it.

First, fill out [.secrets-example](docker/configuration/.secrets-example) as `.secrets` with the appropriate values.

Create `/docker/configuration/.grafanaAdminPassword` to the password for the admin you want, 
then create an empty `/docker/configuration/.grafanaServiceCredential` file for it to populate with a persistent credential.

Once you have set up the server, go to your AdGoBye installation and set the following values:
```json
"SendUnmatchedObjectsToDevs": true,
"BlocklistUnmatchedServer": "http://<ServerIP>",
```
AdGoBye will then report any blocklist misses to the server and the server will process it for the database if appropriate.