# ShortUrlApp

## Demo
<a href="https://youtu.be/XmaZGW1bTVw"><img src="https://i9.ytimg.com/vi_webp/XmaZGW1bTVw/mq2.webp?sqp=CMjkpaoG-oaymwEmCMACELQB8quKqQMa8AEB-AHUCYAC0AWKAgwIABABGGUgZShlMA8=&rs=AOn4CLCBx6jymnNA06mKkb_20WC1H1xoMw" alt="Demo Video"></a>

## How to build and run
Assuming local GoLang environment is set up already
1. Clone the repo 
2. Build the app with ```Go Build```
3. Run the app with ```./ShortUrlApp```

## Endpoints
**Get Url**
* `GET /url/{short_url}`
* Get LongUrl with ShortUrl and redirect
* #### Example cURL
> ```
>  curl http://localhost:8080/url/cf
> ```

**Add a new url record**
* `POST /url`
* #### Example cURL
> ```
>  curl -X POST http://localhost:8080/url -d '{"short_url":"cf", "long_url":"https://www.cloudflare.com", "expire_at":"2023-11-10T00:00:00.000Z"}'
> ```
> a url with no expiration specified will live forever

**Delete a new url record**
* `DELETE /url`
* the corresponding stats will be deleted at the same time. 
* #### Example cURL
> ```
>  curl -X DELETE http://localhost:8080/url/cf
> ```

**Get Stats with ShortUrl**
* `GET /url/{short_url}/stats`
* #### Example cURL
> ```
>  curl http://localhost:8080/url/cf/stats
> ```

## Design
1. Since data persistence is called out in the requirement, I need to use a database to store the data to survive 
computer restarts. Considering scalability and the nature of the short url system being read-heavy, a cache is introduced 
to reduce the traffic on database and cache-aside pattern is used. For the ease of integration, I selected [SQLite](https://github.com/mattn/go-sqlite3) 
as the database and [TTLcache](https://github.com/jellydator/ttlcache) as the cache.
2. DB schema
> ### URLs Table ###
| Column    | Desc                                   | Type      |Constraint
|-----------|----------------------------------------|-----------|-----------------------
| short_url | short url                              | TEXT      |PK
| long_url  | long url for redirection               | TEXT      |Not Null
| expire_at | indicates when the entry should expire | Timestamp |Not Null
Since short urls are always unique and are used in all requests, thus I'm using short_url as the primary key.

> ### Stats Table ###
| Column    | Desc                                 | Type      |Constraint
|-----------|--------------------------------------|-----------|-----------------------
| id        | auto-increment id                    | TEXT      |PK
| short_url | short url                            | TEXT      |FK
| create_at | indicates when the entry was created | Timestamp |Not Null
Auto-increment id is used as primary key. short_url is the foreign key that refers to short_url in URLs table. 
"ON DELETE CASCADE" is enabled so that if the primary key in URLs table is deleted, the corresponding stats 
event rows will be deleted as well. A composite index is created on short_url and created_at to expedite the 
query in GetCount request.




## To Improve
1. Right now, each get url call will write an event record in DB and each get url stats call is reading the 
number of entries in DB between a time range, which are inefficient, especially when we get high frequency requests. 
A few ideas to consider for improvement are:
   1. store the events and do batch writes once every second or so.
   2. instead of storing a new row for each event record, we can only store and update the count for each short url.
2. For easier distribution and deployment, a container like Docker can be used.
