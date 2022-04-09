* Testing
  * Hvordan tester man om CLI laver korrekt config ud fra args, uden at den kører serveren?
    * Eventuelt kan man selv definere en Action eller lignende, så man kan adskille det hele, dvs at man måske bare kan få CLI til at return en config og så manuelt starte serveren i main, i stedet for at CLI også kører serveren?
  * Tests til routeren: https://dev.to/koddr/go-fiber-by-examples-testing-the-application-1ldf
  * Se om andre tests kan omskrives til at bruge assert-pakken (ligesom i linket)
  * Se på code coverage om der er dele, der ikke bliver testet
* Lav et config flag, der sætter default bust routes
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