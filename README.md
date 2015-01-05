## Anubis

### Purpose

Geospatial service using openstreetmap data.

## Loading Data

#### 1. Install osm2pgsql

##### OSX

    brew install protobuf-c
    brew install osm2pgsql --with-protobuf-c

##### Ubuntu

    sudo apt-get install libprotobuf-c0-dev protobuf-c-compiler
    sudo apt-get install osm2pgsql
    
#### 2. Prepare Postgres

    CREATE EXTENSION postgis;
    CREATE EXTENSION hstore;
    
#### 3. Import

    osm2pgsql -c --slim -d dbname -S default.style --hstore -U postgres -P 5432 planet-latest.osm.pbf