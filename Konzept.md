
Ablauf:
1. Liste aller Politiker von Wikimedia abrufen

2. Warten, bis eine Meldung beim [Änderungsstream](https://www.mediawiki.org/wiki/API:Recent_changes_stream) eintritt
   1. Vergleichen, ob der bearbeitete Artikel ein Politiker ist
   2. Wenn nein, dann abbrechen
   3. Eventuell Beleidigungen etc. rausfiltern, also Diffs entfernen, die offensichtlich innerhalb weniger Minuten zurückgerollt werden würden
   4. Änderung abrufen und als Diff-View auf ein Bild rendern
   5. Bild auf Twitter posten
      1. Der Tweet muss einen eindeutigen Hashtag für jeden Politiker enthalten 
