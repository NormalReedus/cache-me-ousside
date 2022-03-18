# cache-me-ousside
### Your favorite Least Recently Used cache microservice 
A simple LRU cache that can be used as a reverse proxy, meaning that any request to this microservice will automatically be passed on to your specified REST API without the REST API having to change anything to accomodate the LRU cache

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


