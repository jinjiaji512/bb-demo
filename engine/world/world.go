/**
 * @Author: jinjiaji
 * @Description:
 * @File:  world.go
 * @Version: 1.0.0
 * @Date: 2020/5/29 4:02 下午
 */

package world

import (
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
			Name:    "team-" + strconv.Itoa(i),
			Players: Players,
		}
		for i := 0; i < 30; i++ {
			for k := 0; k < 5; k++ {
				newPlayer := NewPlayer(Position(k))
				Team.Players = append(Team.Players, newPlayer)
				TotalPlayers = append(TotalPlayers, newPlayer)
			}
		}
	}

	sort.Sort(&TotalPlayers)

	for _, v := range TotalPlayers[:10] {
		v.Show()
	}

	GenGamePlan()
	BeginGamePlan()

}

func GenGamePlan() {
	//每个队伍之间打2场
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

}
