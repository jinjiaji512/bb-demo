/**
 * @Author: jinjiaji
 * @Description:
 * @File:  game.go
 * @Version: 1.0.0
 * @Date: 2021/8/24 下午5:19
 */

package world

type Game struct {
	Status     int // 0 未开始 1 进行中 2 已结束
	Teams      [2]*Team
	GameResult *GameResult
}

type GameTeamResult struct {
	Team  *Team
	IsWin bool
	Score int
	Reb   int
	Steal int
	Block int
}

type GameResult struct {
	TeamResult [2]*GameTeamResult
}

type GamePlayerResult struct {
	Player *Player
	Score  int
	Reb    int
	Steal  int
	Block  int
}

func (g *Game) Play() {

}
