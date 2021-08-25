/**
 * @Author: jinjiaji
 * @Description:
 * @File:  player.go
 * @Version: 1.0.0
 * @Date: 2020/5/29 4:17 下午
 */

package world

import (
	utils "bb/util"
	"encoding/json"
	"fmt"
	"time"
)

type IPlayer interface {
	Say()
	Wait()
	Move()
	Shoot()
	Steal()
	Block()
	Defence()
}

type Position int

const (
	PG Position = 1
	SG Position = 2
	SF Position = 3
	PF Position = 4
	C  Position = 5
)

type Player struct {
	Name     string
	Birthday time.Time
	Position Position

	//以中心点位0,0
	X                     float64
	Y                     float64
	Speed                 float64
	Direction             float64 //速度方向，角度
	Acceleration          float64 //加速度
	AccelerationDirection float64 //加速度方向

	Height float64
	Weight float64
	Strong float64
	Agile  float64
	Energy float64

	SkillShoot   float64
	SkillInShoot float64
	SkillMove    float64
	SkillPass    float64
	SkillDribble float64

	SkillDefence float64
	SkillSteal   float64
	SkillBlock   float64
	SkillRebound float64

	TendencyShoot   float64
	TendencyDribble float64
	TendencyPass    float64

	TendencyDefence float64
	TendencySteal   float64
	TendencyBlock   float64
	TendencyRebound float64

	Score float64

	Data *PlayerData
}

type PlayerData struct {
	Score   float32
	Shoot   int
	ShootIn int
	FG      float32
	Reb     float32
	Steal   float32
	Block   float32
}

func NewPlayer(position Position) *Player {
	var player Player

	player.Name = utils.GetFullName()
	age := utils.NormalFloat(30, 15, 18, 42)
	player.Birthday = time.Now().Add(-time.Hour * 24 * 365 * time.Duration(age))
	player.Position = position

	radH := utils.NormalFloat(85+float32(position*5), 10, 80+float32(position*5), 90+float32(position*5)) / 100
	radHW := utils.NormalFloat(100, 15, 80, 120) / 100
	radStrong := utils.NormalFloat(float32(30+40*(1+radHW)), 20, 0, 100) / 100
	radAgile := utils.NormalFloat(float32(10+30/(1+radHW)+30/(1+radH)), 20, 40, 100) / 100
	radEnergy := utils.NormalFloat(float32(30+20/(1+radHW)+20/(1+radH)), 20, 40, 100) / 100

	player.Height = utils.Decimal(200*radH, 2)
	player.Weight = utils.Decimal(player.Height*player.Height*radHW/200, 2)
	player.Strong = utils.Decimal(100*radStrong, 2)
	player.Agile = utils.Decimal(100*radAgile, 2)
	player.Energy = utils.Decimal(100*(radEnergy+radAgile+radStrong)/3, 2)

	player.SkillShoot = utils.Decimal(utils.NormalFloat(float32(80-10*(radStrong)-10*(radH)), 20, 0, 100), 2)
	player.SkillInShoot = utils.Decimal(utils.NormalFloat(float32(40+40*(radStrong)+20*(radH)), 20, 0, 100), 2)
	player.SkillPass = utils.Decimal(utils.NormalFloat(float32(1+40*(1+radAgile)+10/(1+radH)), 20, 0, 100), 2)
	player.SkillMove = utils.Decimal(utils.NormalFloat(float32(1+30*(1+radAgile)+20/(1+radH)), 20, 0, 100), 2)
	player.SkillDribble = utils.Decimal(utils.NormalFloat(float32(1+30*(1+radAgile)+20/(1+radH)), 10, 0, 100), 2)

	player.SkillDefence = utils.Decimal(utils.NormalFloat(float32(1+20*(1+radAgile)+20*(1+radStrong)+10*(1+radH)), 20, 0, 100), 2)
	player.SkillSteal = utils.Decimal(utils.NormalFloat(float32(1+40*(1+radAgile)+20/(1+radH)), 10, 0, 100), 2)
	player.SkillBlock = utils.Decimal(utils.NormalFloat(float32(1+25*(1+radAgile)+15*(1+radStrong)+15*(1+radH)), 10, 0, 100), 2)
	player.SkillRebound = utils.Decimal(utils.NormalFloat(float32(1+10*(1+radAgile)+20*(1+radStrong)+30*(1+radH)), 10, 0, 100), 2)

	player.TendencyShoot = utils.Decimal(utils.NormalFloat(float32((player.SkillShoot+player.SkillInShoot)/2), 20, 0, 100), 2)
	player.TendencyPass = utils.Decimal(utils.NormalFloat(float32(player.SkillPass), 20, 0, 100), 2)
	player.TendencyDribble = utils.Decimal(utils.NormalFloat(float32(player.SkillDribble), 20, 0, 100), 2)

	player.TendencyDefence = utils.Decimal(utils.NormalFloat(float32(player.SkillDefence), 20, 0, 100), 2)
	player.TendencySteal = utils.Decimal(utils.NormalFloat(float32(player.SkillSteal), 20, 0, 100), 2)
	player.TendencyBlock = utils.Decimal(utils.NormalFloat(float32(player.SkillBlock), 20, 0, 100), 2)
	player.TendencyRebound = utils.Decimal(utils.NormalFloat(float32(player.SkillRebound), 20, 0, 100), 2)

	player.GenScore()
	return &player
}

func (p *Player) GenScore() {
	p.Score = (p.SkillShoot*p.SkillShoot*p.SkillShoot/30 +
		p.SkillInShoot*p.SkillMove*p.SkillDribble/30 +
		p.SkillInShoot*p.SkillInShoot*p.Strong/30 +
		p.SkillDefence*p.SkillSteal*p.Agile/50 +
		p.SkillDefence*p.SkillBlock*p.Strong/50 +
		p.SkillRebound*p.SkillRebound*p.SkillRebound/50 +
		p.Energy*p.Energy*p.Energy/100) / 1200
	p.Score = utils.Decimal(p.Score, 2)
}

func (p *Player) Show() {
	bs, _ := json.MarshalIndent(p, "", "\t")
	fmt.Println(string(bs))
}

//决策
/**
规则 回合制，进攻方先走

进攻 :
1 投篮，根据
2 传球
3 运球

防守 :
1 抢断
2 盖帽
3 移动
*/
func (p *Player) Think() {

}

func (p *Player) Move() {

}

func (p *Player) Shoot() {

}

func (p *Player) Steal() {

}

func (p *Player) Block() {

}

func (p *Player) Defence() {

}

func (p *Player) Rebound() {

}

type PlayerList []*Player

func (l PlayerList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l PlayerList) Less(i, j int) bool {
	return l[i].Score > l[j].Score
}
func (l PlayerList) Len() int {
	return len(l)
}
