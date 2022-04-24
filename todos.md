* Hvis cache bare er et slice / array kunne det sættes for både GET og HEAD?
  * Kan det lade sig gøre med go uden at skulle lave any type?
* Se om man kan bruge iota eller den der auto-inkrementérbare const til utils.ToBytes (se Golang Dojo video om const)
* Testing
  * ~~Tests til routeren: https://dev.to/koddr/go-fiber-by-examples-testing-the-application-1ldf~~
  * Se på code coverage om der er dele, der ikke bliver testet
* Lav et config flag, der sætter default bust routes
  * Eller eventuelt bare gør, så alle unsafe methods buster HELE cachen (se om der er forskel på unsafe og busting etc)
    * Skriv tests først der tjekker at alt bare forsvinder
  * For hvert cached endpoint skal de manipulerende metoder buste
  * Det kræver i hvert fald lister af ikke-manipulerende metoder (GET, HEAD etc) og manipulerende (DELETE, PUT, POST etc)
  * Hvordan sørger man for at patterns matcher korrekt?
* Overhold Cache-Control headers
  * https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control
  * https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching
  * https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching#freshness forklarer hvordan exp virker
  * Skip caching hvis API siger, at man ikke bør
  * Implementér expiration som kan overholde header regler
  * Lyt på når API beder om at buste cache
* Find ud af hvordan man manipulerer hostname, hvis man ikke vil bruge endnu en reverse proxy
* Memorybaseret caching limit (se TODO i config.go)
  * Create a method to return the cache capacity
  * if CapacityUnit is set, use utils.ToBytes to convert the capacity to bytes
  * otherwise return the capacity as a number of entries
  * (maybe there should be something that tells whether we use entries og memory)
  * when busting a cache entry, we should then use utils.MemUsage to compare with the capacity
  * when deciding whether to evict, instead of using entries. Using one over the other should be checked with a bool on the config that is initialized in the factory function, so busting knows whether to use memory or entries
* Print konfigurationen af cachen når den kører
  * Brug Config.String()
* README
  * Beskriv alle flags / config props + hvordan man laver multiple vals (f.eks. flere cache:GET hvor man skal gentage flaget)
  * Beskriv hvordan route params kan bustes med : og hvordan det altid bliver parsed før regex
  * Beskriv hvordan dette er lavet til en almindelig REST api og derfor ikke kan garantere at virke med andre slags API, dvs. at det bygger på safe og unsafe http metoder, hvoraf f.eks. kun GET og HEAD er cacheable og man kan definere entries på baggrund af den route, der er brugt til at requeste entries, da REST er bygget sådan, at routes er lig med ressourcer
    * F.eks. virker det her ikke med GraphQL
  * Beskriv at man både kan køre det som binary, der er en server microservice, men man også kan go get pakken og så importere `cache` for selv at bruge den
  * Beskriv hvordan busting foregår med GET og HEAD i regex
