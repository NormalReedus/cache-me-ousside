* Hvis cache bare er et slice / array kunne det sættes for både GET og HEAD?
  * Kan det lade sig gøre med go uden at skulle lave any type?
* Skriv test til TrimInvalidMethods og brug den alle steder, der laves en config
* Test om man kan undlade så mange make() i config.New(), da den allerede kalder make(BustMap)
* Caching af HEAD (se TODO i config.go) <-- fortsæt med at skrive tests inden implementering af resten af ændringer
  * SE OM IKKE DET KAN LADE SIG GØRE BARE AT PREFIXE ALLE CACHED ENDPOINTS MED METODEN
    * f.eks. `GET:/posts/123`
    * Hvis det kan lade sig gøre, så ignorér alle overstregede punkter, da man ikke behøver flere entry maps
    * Steder der skal ændres med prefixmetoden:
      * ✅ config.example.json5
        * ✅ `cache` skal være et map med GET og HEAD
        * ✅ Busting skal enten explicitere `GET:` etc foran busting eller ikke starte med `^` 
      * ✅ cli.go
        * ✅ addToConfig skal bruge cacheHEAD.Value() (udkommenteret nu) og tilføje værdien til c.Cache["HEAD"] i stedet for c.Cache
        * ✅ Ligeledes skal a.CacheGET gemmes i c.Cache["GET"] frem for i c.Cache
      * ✅ cli_test.go
        * ✅ fjern // fra --cache:HEAD i generateArgs helper
        * ✅ brug korrekt GET: etc syntaks i generateArgs helper eller fjern ^ i start af patterns
          * ✅ Opdatér TestFlagParsings expected til at reflektere dette også
          * ✅ Gør det samme med TestConfigFileParsing
      * ✅ testdata/test.config.json5
        * ✅ Skal have samme ændringer som config.example.json5 og sørge for at de matcher det som cli_test skal bruge
      * router.go
        * ✅ setCachingEndpoints skal have et ekstra loop over conf.Cache for at lave en caching middleware per endpoint i hver method
          * ✅ Dvs. ikke bare brug app.Get, gør det samme med app.Head og sørg for at loops er hen over conf.Cache["GET"] og conf.Cache["HEAD"] frem for bare conf.Cache
      * ✅ router_test.go
        * ✅ TestRouteParams skal bruge ny method-syntaks i patterns, så både test med `^GET:/^HEAD:` og uden brugen af ^ i starten af pattern for at ramme alle methods
      * controllers.go
        * readCacheMiddleware
          * cacheKey skal bruge method med `ctx.Method() + ":" + ctx.OriginalURL()` til at finde cached entries
          * logger.CacheRead skal bare logge cacheKey i stedet for ctx.OriginalURL()
        * writeCacheMiddleware
          * cacheKey skal sættes magen til i readCacheMiddleware, men skal gøres som det første i funktionen
          * logger.CacheSkip skal tage imod cacheKey i stedet for ctx.OriginalURL()
          * logget.CacheWrite skal tage imod cacheKey også
      * ✅ config.go
        * ✅ ændr Config.Cache til at være et map af metoder med endpoints
        * ✅ ændr New() til at oprette .Cache som et map[string][]string og brug make() så de initieres tomme
        * ✅ ValidateRequiredProps skal ikke kun tjekke len(conf.Cache) men len af conf.Cache["HEAD"] og "GET"
      * ✅ config/testdata/test.config.json5
        * ✅ test.config.json5 skal have samme ændringer som config.example.json5
      * ✅ config/testdata/missing-filerne
        * ✅ Skal have ny syntaks
      * ✅ config_test.go
        * ✅ TestRequiredProps skal virke med de nye versioner af cache missing
          * ✅ Lige nu panicer LoadJSON ikke, selvom der er filer, der slet ikke har cache.GET og cache.HEAD
          * ✅ Der skal også laves testfiler, hvor disse keys findes, men bare er tomme slices
        * ✅ TestLoadProps skal assert.NotEmpty på config.Cache["HEAD"] og "GET" i stedet for bare config.Cache
        * ~~TestRequiredProps skal tjekke om cache er et map med indhold, ikke bare et array~~
      * ✅cache_test.go
        * ✅ Hele f ilen skal bruge "GET:" og "HEAD:" til alle test entries
        * ✅ TestMatch skal bruge nye syntaks i patterns til også at teste om den kan matche kun én metode og begge metoder vha. `^GET:` og ved at undlade brug af `^`
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
* Print konfigurationen af cachen når den kører
  * Brug Config.String()
* README
  * Beskriv alle flags / config props + hvordan man laver multiple vals (f.eks. flere cache:GET hvor man skal gentage flaget)
  * Beskriv hvordan route params kan bustes med : og hvordan det altid bliver parsed før regex
  * Beskriv hvordan dette er lavet til en almindelig REST api og derfor ikke kan garantere at virke med andre slags API, dvs. at det bygger på safe og unsafe http metoder, hvoraf f.eks. kun GET og HEAD er cacheable og man kan definere entries på baggrund af den route, der er brugt til at requeste entries, da REST er bygget sådan, at routes er lig med ressourcer
    * F.eks. virker det her ikke med GraphQL
  * Beskriv at man både kan køre det som binary, der er en server microservice, men man også kan go get pakken og så importere `cache` for selv at bruge den
  * Beskriv hvordan busting foregår med GET og HEAD i regex
