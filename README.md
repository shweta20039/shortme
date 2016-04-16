![](logo.png)  
![](https://img.shields.io/badge/version-1.0.0-blue.svg)
### Introduction
----
ShortMe is a url shortening service written in Golang.  
It is with high performance and scalable.  
ShortMe is ready to be used in production. Have fun with it. :)

### Features
----
* Convert same long urls to different short urls.
* Api support
* Short url black list
    * To avoid some words, like `f**k` and `stupid`
    * To make sure that apis such as `/version` and `/health` will only be
    used as api not short urls or otherwise when requesting `http://127.0.0
    .1:3030/version`, version info will be returned rather the long url
    corresponding to the short url "version".
* Base string config in configuration file
    * **Once this base string is specified, it can not be reconfiged anymore
    otherwise the shortened urls may not be unique and thus may conflict with
     previous ones.**

### Implementation
----
Currently, afaik, there are three ways to implement short url service.
* Hash
    * This way is straightforward. However, every hash function will have a 
    collision when data is large.
* Sample
    * this way may contain collision, too. See example below (This example, 
    in Python, is only used to demonstrate the collision situation.).

    ```python
    >>> import random
    >>> import string
    >>> random.sample('abc', 2)
    ['c', 'a']
    >>> random.sample('abc', 2)
    ['a', 'b']
    >>> random.sample('abc', 2)
    ['c', 'b']
    >>> random.sample('abc', 2)
    ['a', 'b']
    >>> random.sample('abc', 2)
    ['b', 'c']
    >>> random.sample('abc', 2)
    ['b', 'c']
    >>> random.sample('abc', 2)
    ['c', 'a']
    >>>
    ```
* Base
    * Just like converting bytes to base64 ascii, we can convert base10 to base62 
    and then make a map between **0 .. 61** to **a-zA-Z0-9**. At last, we can 
    get a unique string if we can make sure that the integer is unique. 
    So, the URL shortening question transforms into making sure we can get a 
    unique integer. 
    ShortMe Use [the method that Flicker use](http://code.flickr.net/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/) 
    to generate a unique integer. 
    Currently, we only use one backend db to generate sequence. For multiple 
    sequence counter db configuration see [Deploy#Sequence Database]
    (#Sequence Database)

### Api
----
* `/version`
    * `HTTP GET`
    * Version info
    * Example
        * `curl http://127.0.0.1:3030/version`
* `/health`
    * `HTTP GET`
    * Health check
    * Example
        * `curl http://127.0.0.1:3030/health`
* `/short`
    * `HTTP POST`
    * Short the long url
    * Example
        * `curl -X POST -H "Content-Type:application/json" -d "{\"longURL\": \"http://www.google.com\"}" http://127.0.0.1:3030/short`
* `/{a-zA-Z0-9}{1,11}`
    * `HTTP GET`
    * Expand the short url and return a **temporary redirect** HTTP status
    * Example
        * `curl -v http://127.0.0.1:3030/3`

        ```bash
            *   Trying 127.0.0.1...
            * Connected to 127.0.0.1 (127.0.0.1) port 3030 (#0)
            > GET /3 HTTP/1.1
            > Host: 127.0.0.1:3030
            > User-Agent: curl/7.43.0
            > Accept: */*
            >
            < HTTP/1.1 307 Temporary Redirect
            < Location: http://www.google.com
            < Date: Fri, 15 Apr 2016 07:25:24 GMT
            < Content-Length: 0
            < Content-Type: text/plain; charset=utf-8
            <
            * Connection #0 to host 127.0.0.1 left intact
        ```

### Install
----
#### Dependency
----
* Golang
* Mysql

#### Compile
----
```bash
mkdir -p $GOPATH/src/github.com/andyxning
cd $GOPATH/src/github.com/andyxning
git clone https://github.com/andyxning/shortme.git

cd shortme
go get ./...
go build -o shortme main.go
```

#### Database Schema
----
We use two databases. Import the two schemas.
* shortme
    * Store short url info
    * [shortme schema](schema/shortme.sql)
* sequence
    * sequence generator
    * [sequence schema](schema/sequence.sql)

#### Configuration
----
```
[http]
# Listen address
listen = "0.0.0.0:3030"

[sequence_db]
# Mysql sequence generator DSN
dsn = "sequence:sequence@tcp(127.0.0.1:3306)/sequence"

# Mysql connection pool max idle connection
max_idle_conns = 4

# Mysql connection pool max open connection
max_open_conns = 4

[short_db]
# Mysql short service read db DSN
read_dsn = "shortme_w:shortme_w@tcp(127.0.0.1:3306)/shortme"

# Mysql short service write db DSN
write_dsn = "shortme_r:shortme_r@tcp(127.0.0.1:3306)/shortme"

# Mysql connection pool max idle connection
max_idle_conns = 8

# Mysql connection pool max open connection
max_open_conns = 8
```
#### Capacity
----
We use an Mysql `unsigned bigint` type to store the sequence counter. According
 to the [Mysql doc](http://dev.mysql.com/doc/refman/5.7/en/integer-types.html) we can get `18446744073709551616` different integers. However, according to [Golang doc about `LastInsertId`](https://golang.org/pkg/database/sql/driver/#RowsAffected.LastInsertId) the returned auto increment integer can only be `int64` which will make the sequence smaller than `uint64`. Even through, we can still get `9223372036854775808` different integers and this will be large enough for most service.  

Supposing that  we consume `100,000,000` short urls one day, then the 
sequence counter can last for `2 ** 63 / 100000000 / 365 = 252695124` years.

#### Grant
----
After setting up the databases and before running `shortme`, make sure that the corresponding user and password has been granted. After logging in mysql console, run following sql statement:
* `grant insert, delete on sequence.* to 'sequence'@'%' identified by 'sequence'`
* `grant insert on shortme.* to 'shortme_w'@'%' identified by 'shortme_w'`
* `grant select on shortme.* to 'shortme_r'@'%' identified by 'shortme_r'`

#### Run
----
* make sure that `static` directory will be at the same directory as `shortme`
* `./shortme -c config.conf`

### Deploy
----

#### <a name="Sequence Database"></a>Sequence Database
----
In the [Flickr blog](http://code.flickr.net/2010/02/08/ticket-servers-distributed-unique-primary-keys-on-the-cheap/),
Flickr suggests that we can use two databases with one for even sequence and
the other one for odd sequence. This will make sequence generator being more
available in case one database is down and will also spread the load about
generate sequence. After splitting sequence db from one to more, we can use 
[HaProxy](http://www.haproxy.org/) as a reverse proxy and thus more sequence 
databases can be used as one. As for load balance algorithm, i think **round 
robin** is good enough for this situation.

In two databases situation, we should add the following configuration to each
 database configuration file.
* First database
   
```
auto_increment_offset 0
auto_increment_increment 2
```

* Second databse

```
auto_increment_offset 1
auto_increment_increment 2
```

Then each time to generate a sequence counter, we can execute below sql 
statement:  
`replace into sequence(stub) values("sequence")`

In cases we use three databases as sequence counter generator, we should 
insert a record for each table in two databases.
* First database

```
auto_increment_offset 0
auto_increment_increment 3
```

* Second database

```
auto_increment_offset 1
auto_increment_increment 3
```

* Third database

```
auto_increment_offset 2
auto_increment_increment 3
```

Then each time to generate a sequence counter, we can execute below sql 
statement:  
`replace into sequence(stub) values("sequence")`

Ok, i think you get the point. When using `N` databases to generate sequence 
counter, configuration for each database configuration file will just 
like below:

```
for i := range N {
    add "auto_increment_offset i" to config file 
    add "auto_increment_increment N" to config file
}

```
So, sequence generator can be horizontally scalable.

