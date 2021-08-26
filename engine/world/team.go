/**
 * @Author: jinjiaji
 * @Description:
 * @File:  team
 * @Version: 1.0.0
 * @Date: 2020/5/29 4:18 下午
 */

package world

type Team struct {
	Name         string
	Players      []*Player
	WinCount     int
	LostCount    int
	CurrGameTeam map[Position]*Player
}
