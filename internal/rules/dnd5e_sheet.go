package rules

import (
	"fmt"
	"io"
	"strings"

	"char-list/internal/sheet"
)

func sign(n int) string { return fmt.Sprintf("%+d", n) }

func render5e(c Char5e, d Derived5e, w io.Writer) error {
	doc, err := sheet.New("dnd5e.pdf", 612, 792)
	if err != nil {
		return err
	}
	if err := doc.Page(1); err != nil {
		return err
	}
	doc.Font("bold")
	doc.Fit(48, 66, 195, 16, c.Name)
	doc.Font("serif")
	lvl := c.Level
	if lvl < 1 {
		lvl = 1
	}
	doc.Fit(270, 53, 100, 9, fmt.Sprintf("%s %d", c.Class, lvl))
	doc.Fit(380, 53, 100, 9, c.Background)
	doc.Fit(490, 53, 90, 9, c.Player)
	doc.Fit(270, 80, 100, 9, c.Race)
	doc.Fit(380, 80, 100, 9, c.Alignment)
	doc.Fit(490, 80, 90, 9, fmt.Sprintf("%d", c.XP))

	tops := map[string]float64{"str": 140, "dex": 213, "con": 286, "int": 358, "wis": 431, "cha": 504}
	doc.Font("bold")
	for _, a := range abilities {
		doc.Center(55, tops[a]+18, 20, fmt.Sprintf("%d", d.Final[a]))
		doc.Center(55, tops[a]+44, 10, sign(d.Mods[a]))
	}
	doc.Font("serif")
	doc.Center(105.4, 166.5, 12, sign(d.Prof))

	for i, a := range abilities {
		y := 208.5 + float64(i)*13.5
		if contains(d.SaveProf, a) {
			doc.Center(104.2, y-0.9, 7, "●")
		}
		doc.Text(113, y, 8, sign(d.Saves[a]))
	}
	for i, s := range d5.Skills {
		y := 323.8 + float64(i)*13.5
		if contains(c.Skills, s.Key) {
			doc.Center(104.2, y-1.1, 7, "●")
		}
		doc.Text(113, y, 8, sign(d.SkillBonus[s.Key]))
	}

	doc.Font("bold")
	doc.Center(246, 147, 16, fmt.Sprintf("%d", d.AC))
	doc.Center(302, 150, 18, sign(d.Initiative))
	doc.Center(356, 150, 18, fmt.Sprintf("%d", d.Speed))
	doc.Font("serif")
	doc.Text(300, 199.8, 9, fmt.Sprintf("%d", d.MaxHP))
	doc.Center(258, 335, 11, d.HitDie)

	for i, a := range d.Attacks {
		if i >= 3 {
			break
		}
		y := 397.25 + float64(i)*20.4
		doc.Fit(228, y, 57, 7.5, a.Name)
		doc.Center(307, y, 8, sign(a.Bonus))
		doc.Fit(329, y, 59, 7.5, a.Damage)
	}
	var spells []string
	if d.Caster != "" {
		spells = append(spells, fmt.Sprintf("Базовая характеристика: %s. КС заклинаний %d, бонус атаки %s. %s", d.Caster, d.SpellDC, sign(d.SpellAtk), d.ArmorNote))
	} else {
		spells = append(spells, d.ArmorNote)
	}
	if strings.TrimSpace(c.Spells) != "" {
		spells = append(spells, c.Spells)
	}
	doc.Lines(232, 452.3, 158, 11, 7.5, 11, strings.Join(spells, "\n"))

	doc.Lines(426, 142.7, 154, 12.3, 7.5, 4, c.Traits)
	doc.Lines(426, 211.7, 154, 12.5, 7.5, 3, c.Ideals)
	doc.Lines(426, 267.2, 154, 12.5, 7.5, 3, c.Bonds)
	doc.Lines(426, 322.7, 154, 12.5, 7.5, 3, c.Flaws)

	doc.Font("bold")
	doc.Center(43, 595, 11, fmt.Sprintf("%d", d.Passive))
	doc.Font("serif")
	doc.Lines(34, 616.5, 168, 11, 7.5, 13, c.Languages)
	doc.Lines(268, 595.5, 120, 11, 7.5, 15, c.Equipment)
	doc.Lines(424, 584.5, 156, 11, 7.5, 16, c.Features)
	return doc.Write(w)
}
