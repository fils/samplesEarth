# Start from scratch image and add in a precompiled binary
# CGO_ENABLED=0 env go build .
# docker build -t earthcube/p418webui:latest -t earthcube/p418webui:0.0.4 .
# docker run -d -p 9900:9900  opencoredata/ocdweb:0.0.1
FROM scratch

# Add in the static elements (could also mount these from local filesystem) later as the indexes grow
ADD cmd/samplesearth/samplesearth /
ADD ./web /web
ADD ./data /data 
#  NOTE.. content for /uis will be mounted in the compose file via datavolume
# do this as a data mount now.. indexes are large
# ADD ./indexes  /indexes 

# Add our binary
CMD ["/samplesearth"]

# Document that the service listens on this port
EXPOSE 9990
