# poliwiki
["Politisches Wikipedia"](https://twitter.com/politischeswiki) ist ein Twitter-Bot, der interessante Änderungen an Einträgen von Politikern auf Twitter postet.


#### Was ist interessant?
Im Moment postet der Bot nur Änderungen, bei denen sich die Länge des Wikipedia-Eintrags um 50 oder mehr Zeichen ändert. Dabei ist es möglich, dass große Änderungen, nach denen die Länge ungefähr gleich ist, ignoriert werden.

#### Wie werden Seiten von Politiker gefunden?
Politiker sind im Sinne des Bots alle WikiData-Objekte, die eine [abgeordnetenwach.de id](https://www.wikidata.org/wiki/Property:P5355) haben.

#### Wie werden Änderungen gefunden?
Wikimedia stellt einen [Stream für Änderungen](https://wikitech.wikimedia.org/wiki/Event_Platform/EventStreams) bereit. Dieser wird vom Bot so gefiltert, dass nur noch Änderungen am deutschen Wikipedia, die nicht von Bots gemacht wurden, betrachtet werden.

Dann findet ein Abgleich mit den Titeln der Seiten zu zuvor abgefragten Politiker statt. Wird hier eine Änderung gefunden, wird der Längenunterschied des Artikels betrachtet. 

Mit 50 oder mehr Zeichen Unterschied macht der Bot einen Screenshot der Seite und postet diesen mit Link und Name des Politiker.


### Vorschläge & Änderungen
Falls du Ideen für Änderungen hast, kannst du sie gerne dem Bot per DM oder direkt hier auf GitHub vorschlagen. Auch gerne gesehen sind Änderungsvorschläge am Code :)

