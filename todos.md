* Testing
  * Tjek om alle public metoder i cache.go fungerer, hvis man giver dem forskellige inputs (evt noget 1.18 fuzzing mht endpoints / patterns?)
* Find ud af hvordan man manipulerer hostname, hvis man ikke vil bruge endnu en reverse proxy
* Caching af HEAD og andre GET-lignende requests (se TODO i config.go)
* Memorybaseret caching limit (se TODO i config.go)