package rules

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"char-list/internal/dice"
	"char-list/internal/embed"
)

var sd = embed.Shadowdark()

type CharSD struct {
	Name       string         `json:"name" yaml:"name"`
	Ancestry   string         `json:"ancestry" yaml:"ancestry"`
	Class      string         `json:"class" yaml:"class"`
	Level      int            `json:"level" yaml:"level"`
	XP         int            `json:"xp" yaml:"xp"`
	Alignment  string         `json:"alignment" yaml:"alignment"`
	Title      string         `json:"title" yaml:"title"`
	Background string         `json:"background" yaml:"background"`
	Deity      string         `json:"deity" yaml:"deity"`
	Scores     map[string]int `json:"scores" yaml:"scores"`
	HP         int            `json:"hp" yaml:"hp"`
	Talents    string         `json:"talents" yaml:"talents"`
	Attacks    string         `json:"attacks" yaml:"attacks"`
	Gear       string         `json:"gear" yaml:"gear"`
	GP         int            `json:"gp" yaml:"gp"`
	SP         int            `json:"sp" yaml:"sp"`
	CP         int            `json:"cp" yaml:"cp"`
}

type DerivedSD struct {
	Mods    map[string]int `json:"mods" yaml:"mods"`
	AC      int            `json:"ac" yaml:"ac"`
	ACNote  string         `json:"acNote" yaml:"acNote"`
	Slots   int            `json:"slots" yaml:"slots"`
	XPNext  int            `json:"xpNext" yaml:"xpNext"`
	Title   string         `json:"title" yaml:"title"`
	HitDie  string         `json:"hitDie" yaml:"hitDie"`
	Attacks []Attack5e     `json:"attacks" yaml:"attacks"`
}

type Shadowdark struct{}

func (Shadowdark) ID() string   { return "shadowdark" }
func (Shadowdark) Name() string { return "Shadowdark" }
func (Shadowdark) Blurb() string {
	return "Мрачная старая школа. Характеристики: 3к6 по порядку, без перестановок. Нет ни одной 14+ — бросаем всё заново."
}
func (Shadowdark) Method() string { return "3к6 по порядку" }

func (Shadowdark) Data() any {
	return map[string]any{
		"abilities":   abilities,
		"abilityRu":   abilityRu,
		"ancestries":  sd.Ancestries,
		"classes":     sd.Classes,
		"alignments":  sd.Alignments,
		"deities":     sd.Deities,
		"backgrounds": sd.Backgrounds,
	}
}

func modSD(score int) int {
	switch {
	case score <= 3:
		return -4
	case score <= 5:
		return -3
	case score <= 7:
		return -2
	case score <= 9:
		return -1
	case score <= 11:
		return 0
	case score <= 13:
		return 1
	case score <= 15:
		return 2
	case score <= 17:
		return 3
	}
	return 4
}

func ancestrySDByName(n string) *embed.AncestrySD {
	for i := range sd.Ancestries {
		if sd.Ancestries[i].Name == n {
			return &sd.Ancestries[i]
		}
	}
	return nil
}

func classSDByName(n string) *embed.ClassSD {
	for i := range sd.Classes {
		if sd.Classes[i].Name == n {
			return &sd.Classes[i]
		}
	}
	return nil
}

func rollStatsSD(method string) (map[string]int, []dice.Log) {
	if method == "4d4" {
		var logs []dice.Log
		scores := map[string]int{}
		for _, a := range abilities {
			l := dice.RollPlus(abilityRu[a], 4, 4, 4)
			logs = append(logs, l)
			scores[a] = l.Total
		}
		return scores, logs
	}
	for attempt := 0; ; attempt++ {
		var logs []dice.Log
		scores := map[string]int{}
		ok := false
		for _, a := range abilities {
			l := dice.Roll(abilityRu[a], 3, 6)
			logs = append(logs, l)
			scores[a] = l.Total
			if l.Total >= 14 {
				ok = true
			}
		}
		if ok || attempt > 20 {
			return scores, logs
		}
	}
}

func titleSD(cl *embed.ClassSD, alignment string) string {
	if cl == nil {
		return ""
	}
	switch alignment {
	case "Законный":
		return cl.Titles[0]
	case "Хаотичный":
		return cl.Titles[2]
	}
	return cl.Titles[1]
}

func (e Shadowdark) Random(raw json.RawMessage, fields []string, method string) (any, []dice.Log, error) {
	c := CharSD{Level: 1, Scores: map[string]int{}}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &c); err != nil {
			return nil, nil, err
		}
	}
	var logs []dice.Log
	if c.Level < 1 {
		c.Level = 1
	}
	if has(fields, "ancestry") || c.Ancestry == "" {
		r := dice.Roll("Происхождение", 1, 12)
		logs = append(logs, r)
		idx := []int{0, 0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 5}[r.Total-1]
		c.Ancestry = sd.Ancestries[idx].Name
	}
	if has(fields, "class") || c.Class == "" {
		r := dice.Roll("Класс", 1, 4)
		logs = append(logs, r)
		c.Class = sd.Classes[r.Total-1].Name
	}
	if has(fields, "alignment") {
		r := dice.Roll("Мировоззрение", 1, 6)
		logs = append(logs, r)
		c.Alignment = sd.Alignments[[]int{0, 0, 0, 1, 1, 2}[r.Total-1]]
	}
	if has(fields, "background") {
		r := dice.Roll("Предыстория", 1, 20)
		logs = append(logs, r)
		c.Background = sd.Backgrounds[r.Total-1]
	}
	if has(fields, "deity") {
		c.Deity = dice.Pick(sd.Deities)
	}
	anc := ancestrySDByName(c.Ancestry)
	cl := classSDByName(c.Class)
	if has(fields, "name") {
		syl := sd.Ancestries[0].Syllables
		if anc != nil {
			syl = anc.Syllables
		}
		c.Name = genName(syl)
	}
	if has(fields, "scores") {
		var l []dice.Log
		c.Scores, l = rollStatsSD(method)
		logs = append(logs, l...)
	}
	if has(fields, "title") {
		c.Title = titleSD(cl, c.Alignment)
	}
	if has(fields, "hp") {
		con := modSD(c.Scores["con"])
		die := 6
		if cl != nil {
			die = cl.HitDie
		}
		r := dice.Roll("Хиты", 1, die)
		logs = append(logs, r)
		hp := r.Total
		if anc != nil && anc.Name == "Дварф" {
			r2 := dice.Roll("Хиты (преимущество)", 1, die)
			logs = append(logs, r2)
			hp = max(hp, r2.Total) + 2
		}
		c.HP = max(hp+con, 1)
	}
	if has(fields, "talents") {
		var lines []string
		if anc != nil {
			lines = append(lines, "Происхождение — "+anc.Trait)
		}
		if cl != nil {
			lines = append(lines, cl.Features...)
			n := 1
			if anc != nil && anc.Name == "Человек" {
				n = 2
			}
			for i := 0; i < n; i++ {
				r := dice.Roll("Талант", 2, 6)
				logs = append(logs, r)
				lines = append(lines, fmt.Sprintf("Талант (%d): %s", r.Total, cl.Talents[r.Total]))
			}
			if cl.NSpells > 0 {
				lines = append(lines, "Заклинания 1 круга: "+strings.Join(dice.PickN(cl.Spells, cl.NSpells), ", "))
			}
		}
		c.Talents = strings.Join(lines, "\n")
	}
	if has(fields, "gear") {
		g := dice.Roll("Золото", 2, 6)
		logs = append(logs, g)
		c.GP = g.Total * 5
		c.SP, c.CP = 0, 0
		var items []string
		if cl != nil {
			for _, k := range cl.Kit {
				if k.Cost <= c.GP {
					c.GP -= k.Cost
					items = append(items, k.Name)
				}
			}
		}
		n := dice.Roll("Снаряжение", 1, 4)
		logs = append(logs, n)
		for i := 0; i < n.Total; i++ {
			r := dice.Roll("Предмет", 1, 12)
			logs = append(logs, r)
			items = append(items, sd.Gear[r.Total-1])
		}
		items = append(items, "Рюкзак", "Рацион (3)")
		c.Gear = strings.Join(items, "\n")
	}
	if has(fields, "attacks") {
		c.Attacks = ""
	}
	return c, logs, nil
}

func (e Shadowdark) derive(c CharSD) DerivedSD {
	d := DerivedSD{Mods: map[string]int{}}
	for _, a := range abilities {
		if c.Scores[a] == 0 {
			c.Scores[a] = 10
		}
		d.Mods[a] = modSD(c.Scores[a])
	}
	cl := classSDByName(c.Class)
	gear := strings.ToLower(c.Gear)
	d.AC = 10 + d.Mods["dex"]
	d.ACNote = "без доспеха"
	switch {
	case strings.Contains(gear, "латы"):
		d.AC = 15
		d.ACNote = "латы"
	case strings.Contains(gear, "кольчуг"):
		d.AC = 13 + d.Mods["dex"]
		d.ACNote = "кольчуга"
	case strings.Contains(gear, "кожан"):
		d.AC = 11 + d.Mods["dex"]
		d.ACNote = "кожаный доспех"
	}
	if strings.Contains(gear, "щит") {
		d.AC += 2
		d.ACNote += " + щит"
	}
	d.Slots = max(c.Scores["str"], 10)
	if cl != nil && cl.Name == "Воин" {
		d.Slots += max(d.Mods["con"], 1)
	}
	lvl := max(c.Level, 1)
	d.XPNext = lvl * 10
	d.Title = titleSD(cl, c.Alignment)
	if cl != nil {
		d.HitDie = fmt.Sprintf("1к%d", cl.HitDie)
	}
	melee := d.Mods["str"]
	ranged := d.Mods["dex"]
	anc := ancestrySDByName(c.Ancestry)
	if anc != nil && anc.Name == "Полуорк" {
		melee++
	}
	if anc != nil && anc.Name == "Эльф" {
		ranged++
	}
	weapons := map[string]string{"длинный меч": "1к8", "булава": "1к6", "короткий меч": "1к6", "кинжал": "1к4", "посох": "1к4", "дубинка": "1к4", "боевой молот": "1к10", "секира": "1к10", "короткий лук": "1к4", "арбалет": "1к6", "копьё": "1к6", "двуручный меч": "1к12"}
	for _, line := range strings.Split(gear, "\n") {
		for name, dmg := range weapons {
			if strings.Contains(line, name) {
				b := melee
				rng := name == "короткий лук" || name == "арбалет"
				if rng {
					b = ranged
				}
				if name == "кинжал" && ranged > melee {
					b = ranged
				}
				d.Attacks = append(d.Attacks, Attack5e{Name: strings.ToUpper(name[:2]) + name[2:], Bonus: b, Damage: fmt.Sprintf("%s %+d", dmg, b)})
			}
		}
	}
	return d
}

func (e Shadowdark) Derive(raw json.RawMessage) (any, error) {
	var c CharSD
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	return e.derive(c), nil
}

func (e Shadowdark) PDF(raw json.RawMessage, w io.Writer) error {
	var c CharSD
	if err := json.Unmarshal(raw, &c); err != nil {
		return err
	}
	return renderSD(c, e.derive(c), w)
}
