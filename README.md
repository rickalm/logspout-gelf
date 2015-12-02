# Graylog GELF Module for Logspout
This module allows Logspout to send Docker logs in the GELF format to Graylog via UDP.

## Why

Micha Hausler did an initial module for Logspout to output in Gelf format, but the datamodel he chose was based on gelf ideas. This version of the module outputs in the same format as the docker gelf logger.

The dis-advantange to using the docker-gelf logger is the loss of the local docker log (e.g. docker logs <container>). By using LogSpout to effectively "tail" the log, this creates an additional copy of the log sent in Gelf Format.

Personally i'm using this to output to LogStash for storage in ElasticSearch, but I wanted the ability to flip docker containers between "Json-File" and "Gelf" and not have the log entries stored in different formats. I considered a simple LogStash filter but there were some additional fields missing which caused me to fork the project and decided I didn't want the additional overhead of a filter in logstash when it could be done in GO.

Other changed included reformatting the time into 3339 format, gathering the command arguments  

## Build
To build, you'll need to fork [Logspout](https://github.com/gliderlabs/logspout), add the following code to `modules.go` 

```
_ "github.com/rickalm/logspout-gelf"
```
and run `docker build -t $(whoami)/logspout:gelf`

## Run

```
docker run \
    -v /var/run/docker.sock:/tmp/docker.sock \
    -p 8000:80 \
    $(whoami)/logspout:gelf \
    gelf://<graylog_host>:12201

```


## License
MIT. See [License](LICENSE)
