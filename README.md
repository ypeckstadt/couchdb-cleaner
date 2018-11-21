# CouchDB cleaner and compaction

Apache CouchDB is a good database and easy to use but it tends to get big in size quickly. Luckily CouchDB has a lot of cleanup and compaction commands that cleans  up old data and reduce fragmentation of the database.  This can have a serious big impact on the database. Disk usage can go from 100 GB to a single GB after every command is run.

Now CouchDB offers the functionality to do this all automatically by adding the following settings to the CouchDB  `/etc/local.ini` file.

```
;This is a default rule for all databases.
;When database fragmentation (unused versions) reaches 30% of the total
;file size, the database will be compacted.
[compactions]
_default = [{db_fragmentation, "30%"}, {view_fragmentation, "30%"}]
;Optional compaction default that will only allow compactions from 11PM to 4AM
;_default = [{db_fragmentation, "30%"}, {view_fragmentation, "30%"}, {from, "23:00"}, {to, "04:00"}]

;Database compaction settings.
;Databases will be checked every 300s (5min)
;Databases less than 256K in size will not be compacted
[compaction_daemon]
check_interval = 300
min_file_size = 256000
```

There are a few problems with this though:

- Hyperledger Fabric uses their own build couchDB docker image and until 1.2 the local.ini file does not contain this settings.  So you need to map your own local.ini file in docker-compose. In the 1.3 build they have added the settings.
- While in the 1.3 build they have added the settings, (un)luck has it that couchDB 2.2 for some reason ignores the settings and does not do the cleanup and defragmentation. It is a bug. Older Fabric CouchDB builds run on couchDB 2.1 and that version does pickup the settings
- Even when the setup is right and working, the fragmentation works but it is not doing the cleanup.  It seems that some things just can only be triggered manually.



Now the good thing is that couchDB offers the futon API which is just an series of endpoints we can use to trigger all of these actions manually. These actions are:

- GET http://couchdb:port/_all_dbs: gets a list of all the available database
- GET http://couchdb:port/<database>/_design_docs: retries a list of views per database 
- POST http://couchdb:port/<database>/_compact: defrags and cleans up a database
- POST http://couchdb:port/<database>/_compact/<viewID>: compacts a view in a database



### Go CLI application

This CLI application runs the following steps on a specified interval:

- Get all the database for the specified couchDB instances, it can be multiple
- Compact the database
- Get a list of views for this database
- Compact each view in the database



The configuration for the application is done with ENV variables.

for example

```
COUCHDB_CLEANER_URLS="http://127.0.0.1:5984"
COUCHDB_CLEANER_COMPACT_INTERVAL_MS=5000
```



You can run the application by running `make run-cleaner`

To build a docker image run `make docker-build-docker`

Add the following to the docker-compose file to make it part of a Hyperledger Fabric network.

```
  couchdb-cleaner:
    container_name: couchdb-cleaner
    image: couchdb-cleaner
    environment:
      COUCHDB_CLEANER_URLS: http://couchdb:5984
      COUCHDB_CLEANER_COMPACT_INTERVAL_MS: 5000
    depends_on:
      - couchdb
```





