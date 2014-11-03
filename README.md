docker-stats
============

Log into the DO instance docker-stats.

# ssh -i ops/keys/docker.key root@104.131.171.105


REST API
========
Definitions
-----------
URI := hostname:4243  ; X being mapped to labs in C1066
uuid := identifies resources
200  := OK
201  := Resource created
400  := Bad Request: syntax error
401  := Unauthorized Request: either authorization header missing or refused
404  := resource not found
406  := Not Acceptable Content Type: identified resource has response entities that don't match accept headers in the request
409  := Conflict: state of resource does not allow requested operation


