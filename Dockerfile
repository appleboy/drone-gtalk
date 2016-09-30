FROM centurylink/ca-certs

ADD drone-gtalk /

ENTRYPOINT ["/drone-gtalk"]
