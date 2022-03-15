* busting
* dynamiske route selectors
  * f.eks. skal alle queries (/xyz?something=himom) caches under hele OriginalURL, men busting på /xyz skal sandsynligvis kunne buste alle queries der kunne være cached til /xyz
* Find ud af hvordan man manipulere hostname, hvis man ikke vil bruge endnu en reverse proxy