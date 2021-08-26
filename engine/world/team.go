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
	Data         *Data
}

type TeamList []*Team

func (l TeamList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l TeamList) Less(i, j int) bool {
	return l[i].Data.WinCount > l[j].Data.WinCount
}
func (l TeamList) Len() int {
	return len(l)
}
