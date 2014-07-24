package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"strconv"
)


type Player struct {
	name string
	position string
	team string
	fpts float64
	posvalue float64
}

func main() {
	
	
	players := make([]Player, 1000)

	
	qbsstring:=getURL("http://www.fantasypros.com/nfl/projections/qb.php?export=xls")
	qbplayernum:=getPos(qbsstring,players,10,0)
	fmt.Printf("%d qbs imported\n", qbplayernum)
	
	rbsstring:=getURL("http://www.fantasypros.com/nfl/projections/rb.php?export=xls")
	rbplayernum:=getPos(rbsstring,players,25,qbplayernum)
	fmt.Printf("%d rbs imported\n", rbplayernum-qbplayernum)

	wrsstring:=getURL("http://www.fantasypros.com/nfl/projections/wr.php?export=xls")
	wrplayernum:=getPos(wrsstring,players,25,rbplayernum)
	fmt.Printf("%d wrs imported\n", wrplayernum-rbplayernum)
	
	tesstring:=getURL("http://www.fantasypros.com/nfl/projections/te.php?export=xls")
	teplayernum:=getPos(tesstring,players,10,wrplayernum)
	fmt.Printf("%d tes imported\n", teplayernum-wrplayernum)
	
	ksstring:=getURL("http://www.fantasypros.com/nfl/projections/k.php?export=xls")
	kplayernum:=getPos(ksstring,players,10,teplayernum)
	fmt.Printf("%d ks imported\n", kplayernum-teplayernum)
	
	sort(players)
	printlist(players)
	
}

func printlist(players []Player) {
	for _, player := range players {
		if (len(player.name)>0) {
			fmt.Printf("%25s %s\n",player.name,player.team)
		}

	}
}
func sort(players []Player) {
	for a:=0; a<len(players); a++ {
		for b:=0; b<len(players); b++ {
			if (len(players[a].name)>0 && len(players[b].name)>0) {
				if (players[a].posvalue>players[b].posvalue) {
					tmpplayera:=players[a]
					players[a]=players[b]
					players[b]=tmpplayera
				}
			}
		}	
	}
}

func toFloat(strfloat string) float64{
	
	retval,err:=strconv.ParseFloat(strings.Trim(strfloat," "),64)
	if err != nil {
		log.Fatal(err)
	}
	return retval
}

func getPos(qbs string,players []Player, numstarting int, playernum int) int{
	var (
		lines []string
		columns []string
		fptscolumn int
	)
	origplayernum:=playernum
	lines=strings.Split(qbs,"\n")
	for _, line := range lines {
		if (strings.HasPrefix(line, "Player Name")) {
			//this is the header line
			columns=strings.Split(line,"\t")
			for index, column := range columns {
				if (strings.HasPrefix(column,"fpts")) {
					fptscolumn=index
				}
			}
		} else if (len(line)>40) {
			fmt.Printf("%d %s\n",playernum,line)
			columns=strings.Split(line,"\t")
			players[playernum].name=columns[0]
			players[playernum].team=columns[1]
			players[playernum].position="QB"
			players[playernum].fpts=toFloat(columns[fptscolumn])
			playernum++
		}
	}
	
	comparitor:=players[origplayernum+numstarting-1].fpts
	fmt.Printf("%.1f\n",comparitor)
	//playernum=origplayernum
	for i:=origplayernum; len(players[i].name)>0; i++ {
			players[i].posvalue = players[i].fpts-comparitor
			fmt.Printf("%s %.1f %.1f\n",players[i].name, players[i].fpts, players[i].posvalue)
	}
	return playernum
	
}

func getURL(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	retbs, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(retbs)
}