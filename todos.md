* Testing
  * Fortsæt på at teste GetSet (opdel i flere små, den virker ikke nu)
    * Lav grundig tests af hvordan state bør være efter hver operation, både i entries og i linked list (Get, Set, Bust, Evict)
      * Efter hver operation, assert følgende
        * Kun korrekte keys er i entries map
        * LRU er korrekt
        * MRU er korrekt
        * Alle andre entries er i korrekt rækkefølge
* Find ud af hvordan man manipulerer hostname, hvis man ikke vil bruge endnu en reverse proxy
* Caching af HEAD og andre GET-lignende requests (se TODO i config.go)
* Memorybaseret caching limit (se TODO i config.go)