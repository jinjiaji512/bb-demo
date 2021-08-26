/**
 * @Author: jinjiaji
 * @Description:
 * @File:  game.go
 * @Version: 1.0.0
 * @Date: 2021/8/24 下午5:19
 */

package world

import (
	utils "bb/util"
	"fmt"
	"math/rand"
	"time"
)

type Game struct {
	Status     int // 0 未开始 1 进行中 2 已结束
	Teams      [2]*Team
	GameResult *GameResult
	ShowLog    bool
}

type GameTeamResult struct {
	Team             *Team
	IsWin            bool
	Score            int
	Reb              int
	Steal            int
	Block            int
	Shoot            int
	ShootIn          int
	GamePlayerResult map[*Player]*GamePlayerResult
}

type GameResult struct {
	TeamResult     [2]*GameTeamResult
	gameResultChan chan string
}

type GamePlayerResult struct {
	Player  *Player
	Score   int
	Reb     int
	Steal   int
	Block   int
	Shoot   int
	ShootIn int
}

func (g *Game) Show() {
	fmt.Println("开始播报赛场实况")
	over := false
	for !over {
		select {
		case msg := <-g.GameResult.gameResultChan:
			fmt.Println(msg)
			if msg == "over" {
				over = true
				break
			}
		}
	}

}

func (g *Game) WriteLog(msg string) {
	if g.ShowLog {
		g.GameResult.gameResultChan <- msg
	}
}

func (g *Game) Play() {
	//todo 临时模拟
	g.GameResult = &GameResult{
		gameResultChan: make(chan string),
	}
	if g.ShowLog {
		go g.Show()
		fmt.Printf("比赛即将开始：%s vs %s\n", g.Teams[0].Name, g.Teams[1].Name)
		fmt.Println("双方球员：")
		fmt.Println("name\tposition\tscore\tvs\tname\tposition\tscore")
		for i := 5; i >= 1; i-- {
			p1 := g.Teams[0].CurrGameTeam[Position(i)]
			p2 := g.Teams[1].CurrGameTeam[Position(i)]
			fmt.Printf("%s\t\t%s\t\t%d\t\tvs\t\t%s\t\t%s\t\t%d\n",
				p1.Name, PositionM[p1.Position], int(p1.Score),
				p2.Name, PositionM[p2.Position], int(p2.Score),
			)
		}
	}
	//48分钟，总计140回合
	for i, v := range g.Teams {
		g.GameResult.TeamResult[i] = &GameTeamResult{
			Team:             v,
			GamePlayerResult: map[*Player]*GamePlayerResult{},
		}
		for _, v2 := range v.CurrGameTeam {
			g.GameResult.TeamResult[i].GamePlayerResult[v2] = &GamePlayerResult{
				Player: v2,
			}
		}
	}
	for i := 0; i < 140; i++ {
		attackTeam := g.Teams[i%2]
		defenceTeam := g.Teams[(i+1)%2]

		attackTeamResult := g.GameResult.TeamResult[i%2]
		defenceTeamResult := g.GameResult.TeamResult[(i+1)%2]

		attackPlayer := attackTeam.CurrGameTeam[PG]
		defencePlayer := defenceTeam.CurrGameTeam[PG]

		//当前回合有10次小回合
		rand.Seed(time.Now().UnixNano())
		g.WriteLog(attackTeam.Name + "队持球")
		g.WriteLog(g.Teams[i%2].Name)
		for j := 0; j < 10; j++ {
			attackPlayerResult := attackTeamResult.GamePlayerResult[attackPlayer]
			defencePlayerResult := defenceTeamResult.GamePlayerResult[defencePlayer]

			//是否投篮
			if attackPlayer.TendencyShoot*float64(j) > float64(rand.Int31n(500)) || j == int(rand.Int31n(10)) {
				attackPlayerResult.Shoot++
				attackTeamResult.Shoot++
				g.WriteLog(attackPlayer.Name + "选择投篮，防守他的是" + defencePlayer.Name)
				//盖帽
				if (defencePlayer.SkillBlock-80)*2+2*(defencePlayer.Height-attackPlayer.Height) >
					float64(rand.Int31n(100)) {
					defencePlayerResult.Block++
					defenceTeamResult.Block++
					g.WriteLog(attackPlayer.Name + "投篮被" + defencePlayer.Name + "盖了")
					//篮板,所有人随机
					var attGetReb bool
					for ok := false; !ok; {
						randT := rand.Int31n(2)
						randP := rand.Int31n(5) + 1
						p := g.Teams[randT].CurrGameTeam[Position(randP)]
						if p.SkillRebound > float64(rand.Int31n(100)) {
							if attackTeam == g.Teams[randT] {
								//己方拿到
								attGetReb = true
								attackPlayer = p
								defencePlayer = defenceTeam.CurrGameTeam[p.Position]
								attackTeamResult.Reb++
								attackTeamResult.GamePlayerResult[attackPlayer].Reb++
								g.WriteLog("进攻方的" + p.Name + "拿到篮板")
							} else {
								//对方拿到
								attGetReb = false
								defenceTeamResult.Reb++
								defenceTeamResult.GamePlayerResult[p].Reb++
								g.WriteLog("防守方的" + p.Name + "拿到篮板")
							}
							break
						}
					}
					if attGetReb {
						//继续
						continue
					} else {
						//双方交换
						break
					}
				}

				if (attackPlayer.SkillShoot-defencePlayer.SkillDefence)+40 > float64(rand.Int31n(100)) {
					//进
					attackPlayerResult.ShootIn++
					attackTeamResult.ShootIn++
					attackPlayerResult.Score += 2
					attackTeamResult.Score += 2
					g.WriteLog("球进了！" + attackPlayer.Name + "获得2分")
					break
				} else {
					//不进
					//篮板
					//篮板,从c至pg随机
					var attGetReb bool
					for ok := false; !ok; {
						randT := rand.Int31n(2)
						randP := utils.NormalFloat(5, 2, 1.1, 5.9)
						p := g.Teams[randT].CurrGameTeam[Position(randP)]
						if p.SkillRebound > float64(rand.Int31n(100)) {
							if attackTeam == g.Teams[randT] {
								//己方拿到
								attGetReb = true
								attackPlayer = p
								defencePlayer = defenceTeam.CurrGameTeam[p.Position]
								attackTeamResult.Reb++
								attackTeamResult.GamePlayerResult[attackPlayer].Reb++
								g.WriteLog("球不进，不过己方的" + p.Name + "拿到篮板")
							} else {
								//对方拿到
								attGetReb = false
								defenceTeamResult.Reb++
								defenceTeamResult.GamePlayerResult[p].Reb++
								g.WriteLog("球不进，被对方的" + p.Name + "拿到篮板")
							}
							break
						}
					}
					if attGetReb {
						//继续
						continue
					} else {
						//双方交换
						break
					}
				}
			}

			//是否传球
			if j == 0 || attackPlayer.TendencyPass > float64(rand.Int31n(100)) {
				lastAttackPlayer := attackPlayer
				randP := rand.Int31n(5) + 1
				attackPlayer = attackTeam.CurrGameTeam[Position(randP)]
				defencePlayer = defenceTeam.CurrGameTeam[attackPlayer.Position]
				g.WriteLog(lastAttackPlayer.Name + "传给了" + attackPlayer.Name)
			}

			//是否突破

		}
	}
	var winner string
	tr1, tr2 := g.GameResult.TeamResult[0], g.GameResult.TeamResult[1]
	if tr1.Score > tr2.Score {
		tr1.IsWin = true
		tr1.Team.WinCount++
		tr2.Team.LostCount++
		winner = tr1.Team.Name
	} else {
		tr2.IsWin = true
		tr1.Team.LostCount++
		winner = tr2.Team.Name
	}

	for _, v := range g.GameResult.TeamResult {
		for _, v2 := range v.GamePlayerResult {
			v2.Player.Data.Score += v2.Score
			v2.Player.Data.GameCount++
			v2.Player.Data.Shoot += v2.Shoot
			v2.Player.Data.ShootIn += v2.ShootIn
			v2.Player.Data.Reb += v2.Reb
			v2.Player.Data.Steal += v2.Steal
			v2.Player.Data.Block += v2.Block
			v2.Player.Data.ScorePG = float32(v2.Player.Data.Score) / float32(v2.Player.Data.GameCount)
			v2.Player.Data.FG = float32(v2.Player.Data.ShootIn*100) / float32(v2.Player.Data.Shoot)
		}
	}
	g.WriteLog("over")
	//报告

	if g.ShowLog {
		fmt.Println("game is over, winner is ", winner)
		fmt.Printf("the final score is %s:%d vs %s:%d \n",
			tr1.Team.Name,
			tr1.Score,
			tr2.Team.Name,
			tr2.Score)

		fmt.Printf("得分：%d vs %d \n", tr1.Score, tr2.Score)
		fmt.Printf("命中率：%d%% vs %d%% \n", tr1.ShootIn*100/tr1.Shoot, tr2.ShootIn*100/tr2.Shoot)
		fmt.Printf("篮板：%d vs %d \n", tr1.Reb, tr2.Reb)
		fmt.Printf("盖帽：%d vs %d \n", tr1.Block, tr2.Block)

		for _, v := range g.GameResult.TeamResult {
			fmt.Println(tr1.Team.Name, "球员数据：\n位置\t球员\t得分\t命中率\t篮板\t盖帽")
			for _, v2 := range v.GamePlayerResult {
				fmt.Printf("%s\t%s\t%d\t%d\t%d\t%d\n",
					PositionM[v2.Player.Position],
					v2.Player.Name,
					v2.Score,
					v2.ShootIn*100/v2.Shoot,
					v2.Reb,
					v2.Block,
				)
			}
		}
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}

}
