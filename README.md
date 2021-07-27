[![Docker Repository on Quay](https://quay.io/repository/marian/lanuv-nrw-water-level-api/status "Docker Repository on Quay")](https://quay.io/repository/marian/lanuv-nrw-water-level-api)

A little proxy service for water level data presented by [HYGON](https://luadb.lds.nrw.de/LUA/hygon/pegel.php?rohdaten=ja) of [LANUV NRW](https://www.lanuv.nrw.de/).

To run the service in development:

```nohighlight
go run main.go
```

### API

The following endpoints are available:

#### /

The root endpoint gives an array of station names.

#### /:station

Returns the time and value of the last measurement for the station.
