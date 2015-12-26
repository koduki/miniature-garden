#!/bin/bash

cd $(dirname $0)
captainhook -listen-addr=0.0.0.0:8080 -echo -configdir ./captainhook
