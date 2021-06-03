# jungle-flame

Scratch code for learning golang. Everyone dislikes their jungler. Now you can dislike yours with science.


'Works' mostly out of the box. Just make a api_conf.json and pump in your riot api key
```
{"RiotKey":"$ThisIsWhereMyKeyWouldGo"}
```
Nav to proj and set it up (or just compile it):
go run .\main.go .\champions.go .\summoner.go .\matchByAccount.go .\allMatches.go .\liveGame.go  .\static\*

A few endpoints currently exist:
+ /summoner/:name
  + Shows you the last time this summoner fed. Queries recent matches, looks for negative KDA, outputs stats/date/champ etc

+ /kills/:name
  + Find avg kill local on map for first 15 mins of the game. Needs mapped to testJS func still for mapping. 
    
+ /jkl/:name
  + Used to query live games. Grabs kill locale data on both junglers

+ /test
  + Functional demo of dynamically plotting average kill location of recent jungler's jungle games.
  
