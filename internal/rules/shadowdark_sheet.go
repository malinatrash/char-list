package rules

import (
	"fmt"
	"io"
	"strings"

	"char-list/internal/sheet"
)

func renderSD(c CharSD, d DerivedSD, w io.Writer) error {
	doc, err := sheet.New("shadowdark.pdf", 630, 414)
	if err != nil {
		return err
	}
	if err := doc.Page(1); err != nil {
		return err
	}
	doc.Font("bold")
	doc.Fit(224, 50, 112, 13, c.Name)
	doc.Font("serif")
	pos := map[string][2]float64{"str": {45, 95}, "int": {124, 95}, "dex": {45, 144.6}, "wis": {124, 144.6}, "con": {45, 194}, "cha": {124, 194}}
	for _, a := range abilities {
		p := pos[a]
		doc.Center(p[0], p[1], 11, fmt.Sprintf("%d", c.Scores[a]))
		doc.Center(p[0]+33.5, p[1], 11, sign(d.Mods[a]))
	}
	doc.Font("bold")
	doc.Center(62, 240, 22, fmt.Sprintf("%d", c.HP))
	doc.Center(140, 240, 22, fmt.Sprintf("%d", d.AC))
	doc.Font("serif")
	doc.FitCenter(140, 268, 62, 6.5, d.ACNote)
	doc.Center(62, 268, 6.5, "кость "+d.HitDie)

	doc.Fit(252, 88, 86, 10, c.Ancestry)
	doc.Fit(238, 134, 100, 10, c.Class)
	doc.Center(212, 191, 10, fmt.Sprintf("%d", max(c.Level, 1)))
	doc.Center(269.6, 188.3, 9, fmt.Sprintf("%d", c.XP))
	doc.Center(317.4, 188.3, 9, fmt.Sprintf("%d", d.XPNext))
	title := c.Title
	if title == "" {
		title = d.Title
	}
	doc.Fit(228, 225, 110, 10, title)
	doc.Fit(265, 271, 73, 10, c.Alignment)
	doc.Fit(275, 317, 63, 9, c.Background)
	doc.Fit(232, 362, 106, 9, c.Deity)

	attacks := strings.TrimSpace(c.Attacks)
	if attacks == "" {
		var lines []string
		for _, a := range d.Attacks {
			lines = append(lines, fmt.Sprintf("%s: атака %s, урон %s", a.Name, sign(a.Bonus), a.Damage))
		}
		attacks = strings.Join(lines, "\n")
	}
	doc.Box(34, 322, 136, 58, 7.5, 1.25, attacks)
	doc.Box(362, 50, 236, 148, 7.5, 1.3, c.Talents)

	doc.Text(468, 222.3, 8, fmt.Sprintf("%d", c.GP))
	doc.Text(508, 222.3, 8, fmt.Sprintf("%d", c.SP))
	doc.Text(549, 222.3, 8, fmt.Sprintf("%d", c.CP))
	items := strings.Split(strings.TrimSpace(c.Gear), "\n")
	for i, it := range items {
		if i >= 20 || strings.TrimSpace(it) == "" {
			continue
		}
		x, y := 375.0, 244.5+float64(i)*13.4
		if i >= 10 {
			x, y = 463, 244.5+float64(i-10)*13.4
		}
		doc.Fit(x, y, 62, 7, strings.TrimSpace(it))
	}
	doc.Center(564, 258, 9, fmt.Sprintf("%d слотов", d.Slots))
	return doc.Write(w)
}
