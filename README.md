[![MIT License][license-shield]][license-url]
[![Issues][issues-shield]][issues-url]

# cache-me-ousside

## Your favorite Least Recently Used cache
A simple LRU cache that can be used as a proxy or reverse proxy, meaning that any request to this service will automatically be passed on to your specified REST API without the REST API having to change anything to accomodate the LRU cache.

As opposed to most LRU caches, `cache-me-ousside` allows you to specify exactly which cache entries to bust and when. That means that one single `POST` request will no longer clear your whole cache when it doesn't need to.

![cache-me-ousside demo](img/cache-me-ousside.gif)

## Contents
  - [About](#about)
  - [Getting started](#getting-started)
    - [Installation](#installation)
      - [NPM (coming soon)](#npm-coming-soon)
      - [Go](#go)
    - [Usage](#usage)
      - [JSON5 configuration file](#json5-configuration-file)
      - [Environment variables](#environment-variables)
      - [CLI flags](#cli-flags)
    - [Go package (coming soon)](#go-package-coming-soon)
      - [Installation](#installation-1)
      - [Usage](#usage-1)
  - [Configuration](#configuration)
    - [Configuration file path](#configuration-file-path)
      - [CLI flag](#cli-flag)
        - [Flags](#flags)
        - [Example](#example)
  - [Roadmap](#roadmap)
  - [Limitations](#limitations)
  - [Specs](#specs)
  - [Details](#details)

## About
`cache-me-ousside` is a server that will proxy all requests to any REST API and cache results of the configured routes in memory, so you can serve the results faster on the next request without having to do database queries or other expensive operations multiple times.

`cache-me-ousside` is a Least Recently Used cache, which means that when the cache is at capacity, the least recently accessed cache entries will be removed first (the FIFO principle). You can configure the cache capacity to be either a fixed number of entries or a memory limit (coming soon).

What makes this cache different from other LRU caches is that you can specify exactly which entries to remove when data on your API is updated. Do you have separate data, that in no way influence each other? Normally, an unsafe HTTP request to your API (such as POST or PUT) will remove all entries from your cache, but perhaps you only need POST requests that update your `todos` to remove your cached `todos`, so that you don't have to repopulate your cache with your `posts` again. 

<p align="right">(<a href="#top">back to top</a>)</p>

## Getting started
It is super simple to get up and running with `cache-me-ousside`!

### Installation
To install the `cache-me-ousside` binary, you have two options. Both of these allow you to run the `cache-me-ousside` command from anywhere on your computer.

#### NPM (coming soon)
You will need to have [NPM](https://www.npmjs.com/package/npm "NPM package") installed on your computer to use this command.

```sh
npm i -g cache-me-ousside
```

#### Go
You will need to have [Go](https://go.dev/dl/ "Go download page") installed on your computer to use this command.

```sh
go install github.com/magnus-bb/cache-me-ousside
```

<p align="right">(<a href="#top">back to top</a>)</p>

### Usage
There are three different ways of configuring `cache-me-ousside`: a JSON5 file, environment variables, or command line flags. The minimal setup only requires you to specify a cache capacity, a proxy API URL, and a list of routes to cache for either HEAD or GET requests.

See the [configuration section](#configuration) for more details on all of the configuration options, or run the command:
```sh
cache-me-ousside --help
```

#### JSON5 configuration file
Using a [JSON5](https://json5.org/ "JSON5 documentation") configuration file, is the recommended method for configuring the cache.

**Configuration**
```json5
{
  capacity: 500,
  apiUrl: 'https://jsonplaceholder.typicode.com/',
  cache:  {
    GET: [ '/posts', '/posts/:id' ],
  },
}
```
**Command line**
```sh
cache-me-ousside --conf /path/to/config.json5
```

#### Environment variables
**Environment / .ENV**
```sh
CAPACITY=500
API_URL=https://jsonplaceholder.typicode.com/
CACHE_GET=/posts,/posts/:id
```
**Command line**
```sh
cache-me-ousside
```

#### CLI flags
**Command line**
```sh
cache-me-ousside --cap 500 --url https://jsonplaceholder.typicode.com/ --cache:GET /posts,/posts/:id
```

<p align="right">(<a href="#top">back to top</a>)</p>

### Go package (coming soon)
It is possible to use this package to implement an LRU cache in your own project. You will get access to the data structures and methods neccesary for an LRU cache, with which you can do whatever you want.

#### Installation
```sh
go get github.com/magnus-bb/cache-me-ousside/cache
```

#### Usage
```go
package main

import github.com/magnus-bb/cache-me-ousside/cache

func main() {
  c := cache.New(500, "")
  c.Set("key", "value")
  c.Get("key")
}
```

See the [documentation](https://pkg.go.dev/github.com/magnus-bb/cache-me-ousside/cache "package cache documentation") for more information on how to use the `cache` package.

<p align="right">(<a href="#top">back to top</a>)</p>

## Configuration
`cache-me-ousside` can be configured with JSON5, environment variables, or CLI flags. The three configuration methods can be used interchangeably (see the [usage section](#usage) for more details on how to use the configuration methods).

If you are mixing configuration methods, the order of precedence is as follows:
1. Command line flags
2. Environment variables
3. JSON5 configuration file

This means that if the same configuration option is specified more than once, the command line flag will be the one used if it is specified, otherwise the environment variable will be used if it is specified, otherwise the JSON5 configuration file will be used.

Use `cache-me-ousside --help` to see a list of all configuration options in the terminal.

### Configuration file path
**Type**: Path (String)

The file path that points to the JSON5 configuration file with options for the cache. This file can be used as an alternative to all other configuration options.

#### CLI flags
`--config` | `--conf` | `--path`

**Example**
```sh
cache-me-ousside --config /path/to/config.json5
```

#### Environment variables
`CONFIG_PATH` | `CONFIG`

**Example**
```sh
CONFIG_PATH=/path/to/config.json5
```

### Cache capacity
**Required**
**Type**: Uint64
**Restrictions**: Must be greater than 0

The cache capacity denotes how much data can be stored in the cache. The capacity can be either a fixed number of entries or a memory limit (coming soon). When the cache is full, the least recently accessed cache entry will be removed.

Regardless of whether you are setting a capacity of a specific number of entries or an amount of memory for the cache, the cache capacity should be set to a number (see [Cache capacity unit](#cache-capacity-unit) for more details on the two modes).

#### CLI flags
`--capacity` | `--cap`

**Example**
```sh
cache-me-ousside --config /path/to/config.json5 --capacity 500
```

#### Environment variables
`CAPACITY`

**Example**
```sh
CAPACITY=500
```

#### JSON5 configuration property
`capacity`

**Example**
```json5
{
  // ...
  capacity: 500,
  // ...
}
```

### Cache capacity unit (coming soon)
**Type**: String
**Options**: `""` | `"b"` | `"kb"` | `"mb"` | `"gb"` | `"tb"`

The cache capacity unit denotes which type of cache limit you want to impose. Leaving this option out or setting it to an empty string will default the cache capacity to use an entry-based cache limit, meaning that the [cache capacity number](#cache-capacity) will represent the exact number of entries that can be stored in the cache. If you set this option to one of the available units, the cache capacity limit will be set to the corresponding number of bytes.

#### CLI flags
`--capacity-unit` | `--cap-unit` | `--cu`

**Example**
```sh
cache-me-ousside --config /path/to/config.json5 --capacity 500 --capacity-unit mb
```

#### Environment variables
`CAPACITY_UNIT`

**Example**
```sh
CAPACITY=500
CAPACITY_UNIT=mb
```

#### JSON5 configuration property
`capacityUnit`

**Example**
```json5
{
  // ...
  capacity: 500,
  capacityUnit: 'mb',
  // ...
}
```

* Hver feature skal have en beskrivelse af i rækkefølge:
  * Hvad det bruges til (hvilket problem det løser)
  * Hvordan det konfigureres i fil, cli, env
  * Eksempel på output eller effekt eller lignende (hvis applicable)
  * Caveats / bugs / todos / ting man skal være opmærksom på
    * trailing slashes bliver fjernet fra api url
    * requests kender forskel på /posts/ og /posts
* Beskriv features, der endnu ikke er færdige og markér med TODO eller WIP

<p align="right">(<a href="#top">back to top</a>)</p>

## Roadmap
* [ ] GraphQL support
* [ ] Cache expiry
* [ ] Respect cache-related headers
* [ ] Package on NPM.org
* [ ] Public API of package `cache`
* [ ] Allow for a cache prop called `both`, that will apply to both GET and HEAD requests

<p align="right">(<a href="#top">back to top</a>)</p>

## Limitations
* Skal være REST (ikke graphql), da entries bliver cached på baggrund af deres route, da det repræsenterer den ressource man tilgår
  * Måske kommer der en graphql version
* Man kan kun cache GET og HEAD
* 

<p align="right">(<a href="#top">back to top</a>)</p>


## Specs
* Basic REST API som per default bare skal reverse proxy alle requests direkte til et angivet endpoint og give svaret tilbage
* Skal have en JSON config, som man peger på stien til med et flag, når man kører servicen
* Config skal definere
  * reverse proxy host som ALLE requests sendes til (required)
  * max antal gemte cache keys (required) eller eventuelt en størrelse (mbs etc) som cachen ikke må overskride
  * endpoints der skal caches
    * skal inddeles i hvilken metode endpointet bruger (nok mest GET, man kan være hvad som helst)
    * man skal kunne definere helt fast endpoints med faste parametre i både route og query, som så bliver cached
    * man skal kunne definere endpoints med variable parametre i både route og query, men som stadig skal caches med en key, der er den helt konkrete requests url og parametre
  * endpoints der skal buste caches
    * skal inddeles efter metoden brugt (ligesom caching endpoints)
    * skal også både kunne definere faste og variable endpoints (ligesom caching endpoints)
    * men hvert endpoint her skal definere hvilke keys (endpoints og parametre) i cachen de buster
      * Det kunne helt basalt være faste strings eller regex den matcher keys på og fjerner fra cache, men det kunne være fedt at gøre, så der også er en syntaks for at buste en række af endpoints (f.eks. alle fra et bestemt endpoint uanset queryparametre) 
* LRUCache skal have et map[string][]byte eller map[string]string til cache, dvs map[endpoint]jsondata
* LRUCache skal have en doubly linked list til at holde styr på LRU, hvor alle node refererer til deres entry i cachen
* Hvis en entry bliver fetchet og den er i cachen skal den lægges forrest i queuen (MRU)
* Hvis en entry ikke er i memory skal den sættes i cache og lægges forrest i queuen (MRU)
  * Størrelse (om det er enheder eller bytes) skal opdateres efter hver operation i cachen
  * Hvis størrelsen angives i bytes skal hver tilføjelse til cache tjekke størrelsen efter operationen, og fjerne LRU key og tjekke igen iterativt indtil størrelsen er under max
* Man skal kunne Clear() alting i LRUCache
* It would be possible to allow for other methods than GET (POST, PUT) to also cache (non-lazily) in the future
  * this would require a mapping to know how to cache the return value of a POST request (e.g. POST to /posts might create something on /posts/slug)
  * it would also require that the specific PUT requests are mapped (if their endpoints are not the same), and that the router calls the LRUCache with instructions to both Bust and set a key
  * POST etc should not be used for caching, since it breaks with the convention of lazy usage defining what stays in cache
    * it also requires the API to return data from POST requests in the exact same format as what would be returned from a GET
    * and that the POST request endpoint is the same as the GET endpoint, which it often is not (e.g., with variable route params)


## Details
* You should trust proxies from your REST API, but it is not strictly necessary
* The cache will only bust cached routes with requests done through the caching microservice. That means, that changing data through any other method will make the cache fall out of sync with the API
* Your API should always return the exact same data when requesting the same url - the cache uses the url to cache the data.
* You can only cache responses from GET requests (so far)


<!-- MARKDOWN LINKS & IMAGES -->
[license-shield]: https://img.shields.io/github/license/magnus-bb/cache-me-ousside.svg?style=for-the-badge
[license-url]: https://github.com/magnus-bb/cache-me-ousside/blob/main/LICENSE
[issues-shield]: https://img.shields.io/github/issues/magnus-bb/cache-me-ousside.svg?style=for-the-badge
[issues-url]: https://github.com/magnus-bb/cache-me-ousside/issues