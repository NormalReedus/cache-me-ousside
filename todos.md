* busting
* dynamiske route selectors
  * f.eks. skal alle queries (/xyz?something=himom) caches under hele OriginalURL, men busting på /xyz skal sandsynligvis kunne buste alle queries der kunne være cached til /xyz
* Find ud af hvordan man manipulerer hostname, hvis man ikke vil bruge endnu en reverse proxy
* Lav alle logs inde i logger
  * f.eks. fejl ved busting etc
* lav alle fejl med rød evt