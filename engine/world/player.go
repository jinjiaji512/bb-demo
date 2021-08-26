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
	"strconv"
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
	Team     *Team `json:"-"`

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
	Score     int
	ScorePG   float32
	GameCount int
	FG        float32
	Shoot     int
	ShootIn   int
	Reb       int
	Steal     int
	Block     int
}

func NewPlayer(position Position) *Player {
	var player Player
	player.Data = &PlayerData{}

	player.Name = utils.GetFullName()
	age := utils.NormalFloat(30, 15, 18, 42)
	player.Birthday = time.Now().Add(-time.Hour * 24 * 365 * time.Duration(age))
	player.Position = position

	radH := utils.NormalFloat(0, 50, -100, 100) / 100
	radHW := utils.NormalFloat(0, 50, -100, 100) / 100
	radStrong := utils.NormalFloat(float32(60+20*radHW+20*radH), 10, 40, 100) / 100
	radAgile := utils.NormalFloat(float32(60-20*radHW-20*radH), 10, 40, 100) / 100
	radEnergy := utils.NormalFloat(float32(30+40*radStrong+30*radAgile), 10, 40, 100) / 100

	player.Height = utils.Decimal(180+6*float64(position)+15*radH, 2)
	player.Weight = utils.Decimal(player.Height*player.Height*25/10000+20*radHW, 2)
	player.Strong = utils.Decimal(100*radStrong, 2)
	player.Agile = utils.Decimal(100*radAgile, 2)
	player.Energy = utils.Decimal(100*radEnergy, 2)

	player.SkillShoot = utils.Decimal(utils.NormalFloat(float32(20+60*radAgile+40*radH), 20, 0, 100), 2)
	player.SkillInShoot = utils.Decimal(utils.NormalFloat(float32(20+60*radStrong+40*radH), 10, 0, 100), 2)
	player.SkillPass = utils.Decimal(utils.NormalFloat(float32(20+90*radAgile), 20, 0, 100), 2)
	player.SkillMove = utils.Decimal(utils.NormalFloat(float32(20+90*radAgile), 20, 0, 100), 2)
	player.SkillDribble = utils.Decimal(utils.NormalFloat(float32(20+90*radAgile-10*radH), 20, 0, 100), 2)

	player.SkillDefence = utils.Decimal(utils.NormalFloat(float32(40+30*radAgile+30*radStrong+40*radH), 40, 0, 100), 2)
	player.SkillSteal = utils.Decimal(utils.NormalFloat(float32(40+40*radAgile+20*radStrong-40*radH), 40, 0, 100), 2)
	player.SkillBlock = utils.Decimal(utils.NormalFloat(float32(40+30*radAgile+30*radStrong+40*radH), 40, 0, 100), 2)
	player.SkillRebound = utils.Decimal(utils.NormalFloat(float32(40+30*radAgile+30*radStrong+40*radH), 40, 0, 100), 2)

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
	p.Score += p.SkillShoot*p.SkillShoot*p.SkillShoot/40 +
		p.SkillInShoot*p.SkillMove*p.SkillDribble/40 +
		p.SkillInShoot*p.Strong*p.Strong/30 +
		p.SkillDefence*p.SkillSteal*p.Agile/30 +
		p.SkillDefence*p.SkillBlock*p.Strong/30 +
		p.SkillRebound*p.SkillRebound*p.SkillRebound/50 +
		p.Energy*p.Energy*p.Energy/50 +
		p.SkillPass*p.SkillPass*p.SkillPass/50 +
		p.SkillMove*p.SkillDribble*p.Agile/50

	if p.Position >= PF {
		p.Score += p.SkillRebound * p.SkillRebound * p.SkillRebound / 50
		p.Score += p.SkillInShoot * p.SkillInShoot * p.SkillInShoot / 30
	} else {
		p.Score += p.SkillShoot * p.SkillShoot * p.SkillShoot / 30
		p.Score += p.SkillInShoot * p.SkillMove * p.SkillDribble / 40
	}
	p.Score /= 3000
	p.Score += 40
	p.Score = utils.Decimal(p.Score, 2)
}

var PositionM = map[Position]string{1: "控球后卫", 2: "得分后卫", 3: "小前锋", 4: "大前锋", 5: "中锋"}

func (p *Player) Show() {
	bs, _ := json.MarshalIndent(p, "", "\t")
	fmt.Println(string(bs))

	comment := "球员档案：\n姓名：" + p.Name + "\n今年" + strconv.Itoa(time.Now().Year()-p.Birthday.Year()) + "岁,司职" +
		PositionM[p.Position] + "，评分：" + strconv.Itoa(int(p.Score)) + "\n个人标签："

	if p.Score > 95 {
		comment += "超级巨星、"
	} else if p.Score > 90 {
		comment += "明星球员、"
	} else if p.Score > 80 {
		comment += "首发球员、"
	} else if p.Score > 75 {
		comment += "角色球员、"
	} else {
		comment += "板凳球员、"
	}
	if p.Score > 90 && time.Now().Year()-p.Birthday.Year() < 22 {
		comment += "天之骄子、"
	}
	if p.Score > 85 && time.Now().Year()-p.Birthday.Year() < 22 {
		comment += "潜力无限、"
	}
	if p.SkillShoot > 95 {
		if p.SkillMove > 80 && p.SkillDribble > 80 {
			comment += "持球神射手、"
		} else {
			comment += "接球投篮射手、"
		}
	}
	if p.SkillInShoot > 90 && p.SkillMove > 80 && p.SkillDribble > 80 {
		comment += "持球突破精英、"
	}
	if p.SkillInShoot > 90 && p.Position >= PF {
		comment += "内线得分手、"
	}
	if p.SkillPass > 90 && p.SkillMove > 90 && p.SkillDribble > 90 && p.Position <= SF {
		comment += "持球组织精英、"
	}
	if p.SkillRebound > 95 {
		comment += "精英篮板手、"
	}
	if p.SkillDefence > 90 && p.SkillBlock > 90 && p.Position >= PF {
		comment += "护框精英、"
	}
	if p.SkillDefence > 90 && p.SkillSteal > 90 {
		comment += "持球防守精英、"
	}
	comment += "\n身体天赋："
	if p.Height-(180+6*float64(p.Position)) > 10 {
		comment += "身高有绝对优势，"
	} else if p.Height-(180+6*float64(p.Position)) > 5 {
		comment += "身高比较出众，"
	} else if p.Height-(180+6*float64(p.Position)) < -5 {
		comment += "身高比较出众，"
	} else if p.Height-(180+6*float64(p.Position)) < -10 {
		comment += "身高有很大劣势，"
	} else {
		comment += "身高不占优势，"
	}

	if p.Strong > 90 {
		comment += "身体极其强壮，"
		if p.Agile > 80 {
			comment += "并且还十分灵敏，"
		} else if p.Agile < 50 {
			comment += "但是太过迟钝，"
		}
	} else if p.Strong > 80 {
		comment += "身体十分强壮，"
		if p.Agile > 80 {
			comment += "并且还十分灵敏，"
		} else if p.Agile < 50 {
			comment += "但是太过迟钝，"
		}
	} else if p.Strong < 40 {
		comment += "身体十分瘦弱，"
		if p.Agile > 80 {
			comment += "不过十分灵敏，"
		} else if p.Agile < 50 {
			comment += "动作也太过迟钝，"
		}
	} else {
		comment += "力量条件一般，"
		if p.Agile > 80 {
			comment += "不过十分灵敏，"
		} else if p.Agile < 50 {
			comment += "动作也太过迟钝，"
		}
	}

	if p.Energy > 90 {
		comment += "体能非常充沛。"
	} else if p.Energy > 85 {
		comment += "体能出色。"
	} else if p.Energy < 50 {
		comment += "体能是很大问题。"
	} else {
		comment += "体能一般。"
	}

	comment += "\n技术能力："
	if p.SkillShoot > 95 {
		comment += "射术极其精准，"
		if p.SkillDribble > 90 && p.SkillInShoot > 90 && p.SkillMove > 90 {
			comment += "并且有极强的持球突破能力。"
		} else if p.SkillDribble > 70 && p.SkillInShoot > 80 && p.SkillMove > 70 {
			comment += "并且有很强的持球突破能力。"
		} else if p.SkillDribble > 80 {
			comment += "而且兼具持球进攻能力。"
		} else if p.SkillDribble < 60 || p.SkillMove < 60 || p.SkillInShoot < 60 {
			comment += "但是仅限于定点接球投篮。"
		} else {
			comment += "但是持球进攻能力一般。"
		}
	} else if p.SkillShoot > 90 {
		comment += "出色的射手，"
		if p.SkillDribble > 90 && p.SkillInShoot > 90 && p.SkillMove > 90 {
			comment += "并且有极强的持球突破能力。"
		} else if p.SkillDribble > 80 && p.SkillInShoot > 80 && p.SkillMove > 80 {
			comment += "并且有很强的持球突破能力。"
		} else if p.SkillDribble > 90 {
			comment += "而且兼具持球进攻能力。"
		} else if p.SkillDribble < 70 {
			comment += "但是仅限于定点接球投篮。"
		} else {
			comment += "但是持球进攻能力一般。"
		}
	} else if p.SkillShoot > 80 {
		comment += "出色的投篮能力，"
		if p.SkillDribble > 90 && p.SkillInShoot > 90 && p.SkillMove > 90 {
			comment += "并且有极强的持球突破能力。"
		} else if p.SkillDribble > 80 && p.SkillInShoot > 80 && p.SkillMove > 80 {
			comment += "并且有很强的持球突破能力。"
		} else if p.SkillDribble > 90 {
			comment += "而且兼具持球进攻能力。"
		} else if p.SkillDribble < 70 {
			comment += "但是仅限于定点接球投篮。"
		} else {
			comment += "但是持球进攻能力一般。"
		}
	} else if p.SkillShoot > 70 {
		comment += "投射能力一般，"
		if p.SkillDribble > 90 && p.SkillInShoot > 90 && p.SkillMove > 90 {
			comment += "但是有极强的持球突破能力。"
		} else if p.SkillDribble > 80 && p.SkillInShoot > 80 && p.SkillMove > 80 {
			comment += "但是有很强的持球突破能力。"
		} else {
			comment += "而且持球进攻能力也一般。"
		}
	} else if p.SkillShoot > 50 {
		comment += "投篮能力十分糟糕，"
		if p.SkillDribble > 90 && p.SkillInShoot > 90 && p.SkillMove > 90 {
			comment += "但是有极强的持球突破能力。"
		} else if p.SkillDribble > 80 && p.SkillInShoot > 80 && p.SkillMove > 80 {
			comment += "但是有很强的持球突破能力。"
		} else {
			comment += "而且持球进攻能力也一般。"
		}
	} else {
		comment += "完全没有投射能力，"
		if p.SkillDribble > 90 && p.SkillInShoot > 90 && p.SkillMove > 90 {
			comment += "但是有极强的持球突破能力。"
		} else if p.SkillDribble > 80 && p.SkillInShoot > 80 && p.SkillMove > 80 {
			comment += "但是有很强的持球突破能力。"
		} else {
			comment += "而且持球进攻能力也一般。"
		}
	}

	if p.Position >= PF {
		if p.SkillRebound > 90 {
			comment += "篮板能力极强。"
		} else if p.SkillRebound > 80 {
			comment += "篮板能力优秀。"
		} else if p.SkillRebound > 60 {
			comment += "篮板能力一般。"
		} else {
			comment += "篮板能力非常差。"
		}

		if p.SkillInShoot > 90 {
			comment += "内线得分能力非常强悍。"
		} else if p.SkillInShoot > 80 {
			comment += "内线得分能力优秀。"
		} else if p.SkillInShoot > 60 {
			comment += "内线得分能力一般。"
		} else {
			comment += "内线得分能力非常差。"
		}

	} else if p.Position <= SF {
		if p.SkillDribble > 90 {
			comment += "控球能力极强，"
		} else if p.SkillDribble > 80 {
			comment += "控球能力优秀，"
		} else if p.SkillDribble > 60 {
			comment += "控球能力一般，"
		} else {
			comment += "控球能力非常差，"
		}

		if p.SkillPass > 90 {
			comment += "组织能力极强。"
		} else if p.SkillDribble > 80 {
			comment += "组织能力优秀。"
		} else if p.SkillDribble > 60 {
			comment += "组织能力一般。"
		} else {
			comment += "组织能力非常差。"
		}
	}

	if p.SkillDefence > 80 && (p.SkillBlock > 80 || p.SkillSteal > 80) {
		comment += "防守极其强悍。"
	} else if p.SkillDefence > 70 && (p.SkillBlock > 70 || p.SkillSteal > 70) {
		comment += "防守能力优秀。"
	} else if p.SkillDefence > 60 && (p.SkillBlock > 60 || p.SkillSteal > 60) {
		comment += "防守能力一般。"
	} else {
		comment += "防守能力非常差。"
	}

	fmt.Println(comment)

	fmt.Println(p.Name, "球员数据：\n位置\t球员\t得分\t命中率\t篮板\t盖帽")
	fmt.Printf("%s\t%s\t%.2f\t%.2f\t%d\t%d\n",
		PositionM[p.Position],
		p.Name,
		p.Data.ScorePG,
		p.Data.FG,
		p.Data.Reb,
		p.Data.Block,
	)

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
	return l[i].Data.Score > l[j].Data.Score
}
func (l PlayerList) Len() int {
	return len(l)
}
