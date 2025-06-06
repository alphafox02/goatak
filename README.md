# GoATAK Fork Goals- free ATAK/CivTAK server

This fork of GoATAK is meant to evaluate the server aspect only as a map for CoT. Please see the original GoATAK for anything else. 

## GoATAK server features

* v1 (XML) and v2 (protobuf) CoT protocol support
* certificate enrollment (v1 and v2) support
* mission packages management
* datasync / missions basic support
* user management with cli tool
* video feeds management
* visibility scopes for users (devices can communicate and see each other within one scope only)
* default preferences and maps provisioning to connected devices
* ability to log all cot's and cli utility to view cot's log and convert it to json or gpx

you can run it with docker,
using `docker run -p 8088:8088 -p 8080:8080 -p 8999:8999 ghcr.io/kdudkov/goatak_server:latest`

## Client features

* v1 (XML) and v2 (protobuf) CoT protocol support
* SSL connection support, tested with [FreeTakServer](https://github.com/FreeTAKTeam/FreeTakServer)
  , [Argustak](https://argustak.com/) and [urpc.info](https://urpc.info/)
* web-ui, ideal for big screen situation awareness center usage
* unit track - your target unit is always in the center of map
* RedX tool - to measure distance and bearing
* Digital Pointer - send DP position to all other contacts
* Add and edit units on map
* ability to log all cot's and cli utility to view cot's log and convert it to json or gpx

## Libraries used

* [Leaflet](https://leafletjs.com/)
* [Milsymbol](https://github.com/spatialillusions/milsymbol)

[![CI](https://github.com/kdudkov/goatak/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/kdudkov/goatak/actions/workflows/main.yml)

[By me a beer 🍺](https://buymeacoffee.com/kdudkov)
