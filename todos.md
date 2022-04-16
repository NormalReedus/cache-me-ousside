* Testing
  * ~~Tests til routeren: https://dev.to/koddr/go-fiber-by-examples-testing-the-application-1ldf~~
  * Se på code coverage om der er dele, der ikke bliver testet
  * Check alle tests om de har expected og actual i rigtig rækkefølge
* Lav noget i bust controllers der kan læse :id osv syntaks <---
  * Det kunne være, at når en route skal bustes, så loader man self. alle regexes man skal matche med. Men den controller, der skal buste tager først lige alle regexes og hardcoder alle :xxx med tilsvarende route params ved at bruge f.eks. ctx.AllParams() eller lignende
  * Fortsæt/tjek at tilføje replaceRouteParams i bust route controlleren lige inden patterns gives til cachen for at matche virker
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
* Caching af HEAD og andre GET-lignende requests (se TODO i config.go)
* Memorybaseret caching limit (se TODO i config.go)
* Print konfigurationen af cachen når den kører
* README
  * Beskriv alle flags / config props + hvordan man laver multiple vals (f.eks. flere cache:GET hvor man skal gentage flaget)