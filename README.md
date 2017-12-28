# Location Search

Simple script to learn location based queries with redis.

- Insert location database into redis
  - http://download.geonames.org/export/dump/
  ```sh
  wget http://download.geonames.org/export/dump/US.zip
  ```
- Begin to query them
- Play with goroutines and channels

# How Geohashing works

Get geohash in resolution you need

> **Haversine** formula is a standard way to calculate distance between 2 points on the surface of a sphere.  Earth is not a perfect sphere, so this method may introduce errors of up to .5%.  This is usually an acceptable error rate for social proximity searches.

![GeoHash + Haversine](img/GeoHashing.png?raw=true)

High level process is:
- Select point
- retrieve all points within 8 boxes
- apply haversine formula to results and drop from radius.


- Redis uses Haversine formula, error rate may be up to 0.5%
- Add via name, lat, lon
Sample data:
```
4085315	Ragland Cemetery	Ragland Cemetery		34.67398	-86.82695	S	CMTY	US		AL	083			0	180	178	America/Chicago	2006
```
# Loading Data into Redis
- Pipelining
```
geoadd places lon lat name
```
### Real life query patterns

- Get places near me with activities on specific date/range
- Get places near me with activities

- query for 10 closest races near me, get 10 raceids returned
  - pipeline query for data associated with those 10 ids
  - take all 10 queries and shove them into 1 request, get all results back at once.  can be 100x more efficient or more.

- Loading data:
  - 1 record at a time:  15+ minutes for a 2.3 million records
  - pipelined:  32 seconds.  Wow.


### Failure Recovery

- dump geodata into Redis
- as events get added, push events
- refill every night?

### Benchmarking
Memory usage:
initial empty container:  1.4MB
with geodata: 275MB

```sh
time go run main.go load | redis-cli --pipe
All data transferred. Waiting for the last reply...
Last reply received from server.
errors: 0, replies: 2232904

real	0m32.357s
user	0m7.010s
sys	0m5.994s
```

```bash
#Outside goroutine
real	0m2.443s
user	0m0.848s
sys	0m0.352s
```

# Links
- [Best GeoHashing explanation I've found](https://gis.stackexchange.com/questions/18330/using-geohash-for-proximity-searches)
- [MySQL Geohash](https://dev.mysql.com/doc/refman/5.7/en/spatial-geohash-functions.html)
- [Redis Caching Strategies](https://d0.awsstatic.com/whitepapers/Database/database-caching-strategies-using-redis.pdf)
