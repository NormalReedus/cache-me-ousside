* Caching af HEAD (se TODO i config.go) <--
  * SE OM IKKE DET KAN LADE SIG GØRE BARE AT PREFIXE ALLE CACHED ENDPOINTS MED METODEN
    * f.eks. `GET:/posts/123`
    * Hvis det kan lade sig gøre, så ignorér alle overstregede punkter, da man ikke behøver flere entry maps
  * FØRST SKRIV TESTS
  * ~~Kræver at cachens `entries` er et map af methods (GET og HEAD) som så er maps af endpoints med de gemte værdier~~
    * ~~Dette kræver at `cache.Size()` returner length af entries["HEAD"] + entries["GET"]~~
    * ~~Det kræver også at `cache.CachedEndpoints()` returner et map med HEAD og GET som er arrays af keys~~
  * ~~`cache.Get` skal skrives om til også at tage imod en `method` så den kan finde entry i det korrekte map af entries~~
  * ~~`cache.Set` skal også tage imod en `method` ligesom med ovenstående Get~~
  * ~~`Entry` i linked list skal også have en `method` så den ved hvilket map den er gemt, når man f.eks. skal evict~~
  * ~~`cache.Bust` skal også tage imod en method, så den ved hvor den skal finde en entry, men det er okay hvis den tager imod én enkelt method og så variadic keys, for Bust bør altid blive kaldt fra en controller, der self. kun er kaldt med én bestemt method, og på den måde vil Bust aldrig blive kaldt med entries med forskellige methods~~
* Testing
  * ~~Tests til routeren: https://dev.to/koddr/go-fiber-by-examples-testing-the-application-1ldf~~
  * Se på code coverage om der er dele, der ikke bliver testet
* Lav et config flag, der sætter default bust routes
  * Eller eventuelt bare gør, så alle unsafe methods buster HELE cachen (se om der er forskel på unsafe og busting etc)
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
* Print konfigurationen af cachen når den kører
* README
  * Beskriv alle flags / config props + hvordan man laver multiple vals (f.eks. flere cache:GET hvor man skal gentage flaget)
  * Beskriv hvordan route params kan bustes med : og hvordan det altid bliver parsed før regex
  * Beskriv hvordan dette er lavet til en almindelig REST api og derfor ikke kan garantere at virke med andre slags API, dvs. at det bygger på safe og unsafe http metoder, hvoraf f.eks. kun GET og HEAD er cacheable og man kan definere entries på baggrund af den route, der er brugt til at requeste entries, da REST er bygget sådan, at routes er lig med ressourcer
    * F.eks. virker det her ikke med GraphQL
  * Beskriv at man både kan køre det som binary, der er en server microservice, men man også kan go get pakken og så importere `cache` for selv at bruge den