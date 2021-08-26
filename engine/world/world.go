/**
 * @Author: jinjiaji
 * @Description:
 * @File:  world.go
 * @Version: 1.0.0
 * @Date: 2020/5/29 4:02 下午
 */

package world

import (
	"fmt"
	"sort"
	"strconv"
)

var (
	Teams        []*Team
	GamePlan     []*Game
	TotalPlayers PlayerList
)

func NewWorld() {
	for i := 0; i < 10; i++ {
		var Players PlayerList
		Team := Team{
			Name:         "team-" + strconv.Itoa(i),
			Players:      Players,
			CurrGameTeam: make(map[Position]*Player),
		}
		for j := 0; j < 1; j++ {
			for k := 1; k < 6; k++ {
				newPlayer := NewPlayer(Position(k))
				for newPlayer.Score < 60 {
					newPlayer = NewPlayer(Position(k))
				}
				newPlayer.Team = &Team
				Team.Players = append(Team.Players, newPlayer)
				TotalPlayers = append(TotalPlayers, newPlayer)
				if Team.CurrGameTeam[Position(k)] == nil || Team.CurrGameTeam[Position(k)].Score < newPlayer.Score {
					Team.CurrGameTeam[Position(k)] = newPlayer
				}
			}
		}
		Teams = append(Teams, &Team)
	}

	fmt.Println(len(Teams))
	GenGamePlan()
	fmt.Println(len(GamePlan))
	BeginGamePlan()

	sort.Sort(&TotalPlayers)
	for _, v := range TotalPlayers[:10] {
		v.Show()
	}

}

func GenGamePlan() {
	for _, v1 := range Teams {
		for _, v2 := range Teams {
			GamePlan = append(GamePlan, &Game{
				Status: 0,
				Teams:  [2]*Team{v1, v2},
			})
		}
	}
}
func BeginGamePlan() {
	for _, v := range GamePlan {
		//time.Sleep(time.Second)
		v.Play()
	}
}
