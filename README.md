# StreamOps

Outil d'automatisation de streaming écrit en GO, connecté à Twitch (EventSub) et OBS (WebSocket).
Il écoute les événements twitch et les routent vers des actions OBS.

## MVP
- Afficher une source OBS durant X secondes en réaction à un event Twitch (follow par exemple)
- change de scène via API locale
- historisation des events
