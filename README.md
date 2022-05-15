# cache-me-ousside
## Your favorite Least Recently Used cache
A simple LRU cache that can be used as a proxy or reverse proxy, meaning that any request to this service will automatically be passed on to your specified REST API without the REST API having to change anything to accomodate the LRU cache.

As opposed to most LRU caches, `cache-me-ousside` allows you to specify exactly which cache entries to bust and when. That means that one single `POST` request will no longer clear your whole cache when it doesn't need to.

![cache-me-ousside demo](img/cache-me-ousside.gif)

## Contents
- [Specs](#specs)
- [Details](#details)

## What is it?
`cache-me-ousside` is a server that will proxy all requests to any REST API and cache results of the configured routes in memory, so you can serve the results faster on the next request without having to do database queries or other expensive operations multiple times.

`cache-me-ousside` is a Least Recently Used cache, which means that when the cache is at capacity, the least recently accessed cache entries will be removed first (the FIFO principle). You can configure the cache capacity to be either a fixed number of entries or a memory limit (coming soon).

What makes this cache different from other LRU caches is that you can specify exactly which entries to remove when data on your API is updated. Do you have separate data, that in no way influence each other? Normally, an unsafe HTTP request to your API (such as POST or PUT) will remove all entries from your cache, but perhaps you only need POST requests that update your `todos` to remove your cached `todos`, so that you don't have to repopulate your cache with your `posts` again. 

### Package `cache` (public API is a WIP)
* Man kan også bruge cachen som package til sin egen LRU Cache


## Installation
To install the `cache-me-ousside` binary, you have two options. Both of these allow you to run the `cache-me-ousside` command from anywhere on your computer.

### NPM (coming soon)
You will need to have [NPM](https://www.npmjs.com/package/npm "NPM package") installed on your computer to use this command.

```sh
npm i -g cache-me-ousside
```

### Go
You will need to have [Go](https://go.dev/dl/ "Go download page") installed on your computer to use this command.

```sh
go install github.com/magnus-bb/cache-me-ousside
```

## Usage
To run start the 

* Hvordan man kører det
  * Start med at beskrive hvordan man kører, derefter config i næste afsnit
  * Læg en default (ikke example) config med i projektet, der virker ligesom alm. cache hvor alt bustes
* tilføj code block med output fra `--help`

## Configuration
* Beskriv de 3 forskellige måder at config
  * JSON5 (anbefalet), cli, env
* Vis et eksempel med alle features / eksempel / default features af hver slags måde og beskriv syntaks og overwrites
  * cli overskriver fil, men tjek hvordan det er med env
* Referér til All features for at se hver config man kan sætte

## All features
* Beskrivelse af alle de forskellige måder man kan bruge programmet med hver sin undertitel (###)
* Hver feature skal have en beskrivelse af i rækkefølge:
  * Hvad det bruges til (hvilket problem det løser)
  * Hvordan det konfigureres i fil, cli, env
  * Eksempel på output eller effekt eller lignende (hvis applicable)
  * Caveats / bugs / todos / ting man skal være opmærksom på
* Beskriv features, der endnu ikke er færdige og markér med TODO eller WIP, så man i Contents kan se oversigt over nuværende og kommende features


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


