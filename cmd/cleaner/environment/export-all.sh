#!/usr/bin/env bash

source <(sed -E -n 's/[^#]+/export &/ p' ./files/couchdb-cleaner.environment.dev)