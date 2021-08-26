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
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	Teams        TeamList
	GamePlan     []*Game
	TotalPlayers PlayerList
)

func NewWorld() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		var Players PlayerList
		Team := Team{
			Name:         "team-" + strconv.Itoa(i),
			Players:      Players,
			CurrGameTeam: make(map[Position]*Player),
			Data:         &Data{},
		}
		for j := 0; j < 10; j++ {
			for k := 1; k < 6; k++ {
				newPlayer := NewPlayer(Position(k))
				for newPlayer.Score < 70 {
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
	for i := 0; i < 10; i++ {
		GenGamePlan()
	}
	fmt.Println(len(GamePlan))
	BeginGamePlan()

	sort.Sort(&TotalPlayers)
	fmt.Println("球员数据：")
	fmt.Printf("%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s   \n",
		"球员", "位置", "球队", "场均得分", "命中率", "场均篮板", "场均盖帽", "身高", "体重", "总评", "投篮", "内线", "投倾", "篮板", "盖帽", "说明")

	for _, p := range TotalPlayers[0:100] {
		//p.Show()
		fmt.Printf("%-8s\t%-8s\t%-8s\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f\t%-8s   \n",
			p.Name,
			PositionM[p.Position],
			p.Team.Name,
			p.Data.ScorePG,
			p.Data.FG,
			float32(p.Data.Reb)/float32(p.Data.GameCount),
			float32(p.Data.Block)/float32(p.Data.GameCount),

			p.Height,
			p.Weight,
			p.Score,
			p.SkillShoot,
			p.SkillInShoot,
			p.TendencyShoot,
			p.SkillRebound,
			p.SkillBlock,
			strings.ReplaceAll(p.Comment, "\n", ";"),
		)
	}

	sort.Sort(&Teams)
	fmt.Println("球队数据：")
	fmt.Printf("%-8s\t%-8s\t%-8s\t%-8s\t%-8s\t%-8s   \n",
		"胜场", "球队", "场均得分", "命中率", "场均篮板", "场均盖帽")

	for _, p := range Teams {
		//p.Show()
		fmt.Printf("%-8d\t%-8s\t%-8.2f\t%-8.2f\t%-8.2f\t%-8.2f   \n",
			p.Data.WinCount,
			p.Name,
			p.Data.ScorePG,
			p.Data.FG,
			float32(p.Data.Reb)/float32(p.Data.GameCount),
			float32(p.Data.Block)/float32(p.Data.GameCount),
		)
	}

}

func GenGamePlan() {
	for i := 0; i < len(Teams); i++ {
		for j := i + 1; j < len(Teams); j++ {
			v1 := Teams[i]
			v2 := Teams[j]
			gp := Game{
				Status:  0,
				Teams:   [2]*Team{v1, v2},
				ShowLog: false,
			}
			//if v1.Name == "team-24" || v2.Name == "team-24"{
			//	gp.ShowLog = true
			//}
			GamePlan = append(GamePlan, &gp)
		}
	}
}
func BeginGamePlan() {
	for _, v := range GamePlan {
		//time.Sleep(time.Second)
		v.Play()
	}
}
