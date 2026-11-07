package rules

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"char-list/internal/dice"
)

type Char5e struct {
	Name       string         `json:"name"`
	Player     string         `json:"player"`
	Race       string         `json:"race"`
	Class      string         `json:"class"`
	Background string         `json:"background"`
	Alignment  string         `json:"alignment"`
	Level      int            `json:"level"`
	XP         int            `json:"xp"`
	Scores     map[string]int `json:"scores"`
	Skills     []string       `json:"skills"`
	Traits     string         `json:"traits"`
	Ideals     string         `json:"ideals"`
	Bonds      string         `json:"bonds"`
	Flaws      string         `json:"flaws"`
	Languages  string         `json:"languages"`
	Equipment  string         `json:"equipment"`
	Features   string         `json:"features"`
	Spells     string         `json:"spells"`
}

type Derived5e struct {
	Mods       map[string]int `json:"mods"`
	Final      map[string]int `json:"final"`
	Bonus      map[string]int `json:"bonus"`
	Prof       int            `json:"prof"`
	SaveProf   []string       `json:"saveProf"`
	Saves      map[string]int `json:"saves"`
	SkillBonus map[string]int `json:"skillBonus"`
	SkillCount int            `json:"skillCount"`
	SkillPool  []string       `json:"skillPool"`
	MaxHP      int            `json:"maxHP"`
	HitDie     string         `json:"hitDie"`
	AC         int            `json:"ac"`
	ArmorNote  string         `json:"armorNote"`
	Initiative int            `json:"initiative"`
	Speed      int            `json:"speed"`
	Passive    int            `json:"passive"`
	Attacks    []Attack5e     `json:"attacks"`
	SpellDC    int            `json:"spellDC"`
	SpellAtk   int            `json:"spellAtk"`
	Caster     string         `json:"caster"`
}

type Attack5e struct {
	Name   string `json:"name"`
	Bonus  int    `json:"bonus"`
	Damage string `json:"damage"`
}

type DnD5e struct{}

func (DnD5e) ID() string   { return "dnd5e" }
func (DnD5e) Name() string { return "Dungeons & Dragons 5e" }
func (DnD5e) Blurb() string {
	return "Классический лист пятой редакции. Характеристики: 4к6, худший куб отбрасывается, распределение по выбору."
}
func (DnD5e) Method() string { return "4к6, отбросить худший" }

func (DnD5e) Data() any {
	return map[string]any{
		"abilities":   abilities,
		"abilityRu":   abilityRu,
		"skills":      skills5e,
		"races":       races5e,
		"classes":     classes5e,
		"backgrounds": backgrounds5e,
		"alignments":  alignments5e,
	}
}

func race5eByName(n string) *race5e {
	for i := range races5e {
		if races5e[i].Name == n {
			return &races5e[i]
		}
	}
	return nil
}

func class5eByName(n string) *class5e {
	for i := range classes5e {
		if classes5e[i].Name == n {
			return &classes5e[i]
		}
	}
	return nil
}

func background5eByName(n string) *background5e {
	for i := range backgrounds5e {
		if backgrounds5e[i].Name == n {
			return &backgrounds5e[i]
		}
	}
	return nil
}

func rollStats5e() (map[string]int, []dice.Log) {
	var logs []dice.Log
	var totals []int
	for i := range 6 {
		l := dice.DropLowest(fmt.Sprintf("Бросок %d", i+1), 4, 6, 1)
		logs = append(logs, l)
		totals = append(totals, l.Total)
	}
	scores := map[string]int{}
	for i, a := range abilities {
		scores[a] = totals[i]
	}
	return scores, logs
}

func genName(syl []string) string {
	n := 2 + dice.D(2)
	var b strings.Builder
	for range n {
		b.WriteString(syl[dice.D(len(syl))-1])
	}
	s := strings.ToLower(b.String())
	r := []rune(s)
	return strings.ToUpper(string(r[0])) + string(r[1:])
}

func (e DnD5e) Random(raw json.RawMessage, fields []string) (any, []dice.Log, error) {
	c := Char5e{Level: 1, Scores: map[string]int{}}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &c); err != nil {
			return nil, nil, err
		}
	}
	var logs []dice.Log
	if c.Level < 1 {
		c.Level = 1
	}
	if has(fields, "race") || race5eByName(c.Race) == nil {
		c.Race = dice.Pick(races5e).Name
	}
	if has(fields, "class") || class5eByName(c.Class) == nil {
		c.Class = dice.Pick(classes5e).Name
	}
	if has(fields, "background") || background5eByName(c.Background) == nil {
		c.Background = dice.Pick(backgrounds5e).Name
	}
	if has(fields, "alignment") {
		c.Alignment = dice.Pick(alignments5e)
	}
	if has(fields, "name") {
		c.Name = genName(race5eByName(c.Race).Syllables)
	}
	if has(fields, "scores") {
		c.Scores, logs = rollStats5e()
	}
	cl := class5eByName(c.Class)
	bg := background5eByName(c.Background)
	if has(fields, "skills") {
		c.Skills = append([]string{}, bg.Skills...)
		pool := []string{}
		for _, k := range cl.SkillPool {
			if !contains(c.Skills, k) {
				pool = append(pool, k)
			}
		}
		c.Skills = append(c.Skills, dice.PickN(pool, cl.SkillCount)...)
	}
	if has(fields, "traits") {
		c.Traits = dice.Pick(traits5e)
	}
	if has(fields, "ideals") {
		c.Ideals = dice.Pick(ideals5e)
	}
	if has(fields, "bonds") {
		c.Bonds = dice.Pick(bonds5e)
	}
	if has(fields, "flaws") {
		c.Flaws = dice.Pick(flaws5e)
	}
	if has(fields, "languages") {
		r := race5eByName(c.Race)
		c.Languages = strings.Join(r.Languages, ", ")
		if bg.Extra != "" {
			c.Languages += "\n" + bg.Extra
		}
	}
	if has(fields, "equipment") {
		c.Equipment = strings.Join(append(append([]string{}, cl.Equipment...), bg.Equipment...), "\n")
	}
	if has(fields, "features") {
		r := race5eByName(c.Race)
		lines := append([]string{}, r.Traits...)
		lines = append(lines, cl.Features...)
		lines = append(lines, "Предыстория: "+bg.Feature)
		c.Features = strings.Join(lines, "\n")
	}
	if has(fields, "spells") {
		c.Spells = ""
		if cl.Caster != "" {
			c.Spells = "Заговоры: " + strings.Join(dice.PickN(cl.Cantrips, cl.NCantrips), ", ") + "\n1 круг: " + strings.Join(dice.PickN(cl.Spells, cl.NSpells), ", ")
		}
	}
	return c, logs, nil
}

func (e DnD5e) derive(c Char5e) Derived5e {
	d := Derived5e{Mods: map[string]int{}, Final: map[string]int{}, Bonus: map[string]int{}, Saves: map[string]int{}, SkillBonus: map[string]int{}}
	r := race5eByName(c.Race)
	cl := class5eByName(c.Class)
	if c.Level < 1 {
		c.Level = 1
	}
	d.Prof = 2 + (c.Level-1)/4
	for _, a := range abilities {
		b := 0
		if r != nil {
			b = r.Bonus[a]
		}
		d.Bonus[a] = b
		if c.Scores[a] == 0 {
			c.Scores[a] = 10
		}
		d.Final[a] = c.Scores[a] + b
		d.Mods[a] = dice.Mod5e(d.Final[a])
	}
	d.Speed = 30
	if r != nil {
		d.Speed = r.Speed
	}
	for _, a := range abilities {
		d.Saves[a] = d.Mods[a]
	}
	if cl != nil {
		d.SaveProf = cl.Saves
		d.SkillCount = cl.SkillCount
		d.SkillPool = cl.SkillPool
		for _, a := range cl.Saves {
			d.Saves[a] += d.Prof
		}
		d.HitDie = fmt.Sprintf("%dк%d", c.Level, cl.HitDie)
		d.MaxHP = cl.HitDie + d.Mods["con"] + (c.Level-1)*(cl.HitDie/2+1+d.Mods["con"])
		switch cl.ArmorAC {
		case "light11":
			d.AC = 11 + d.Mods["dex"]
			d.ArmorNote = "Кожаный доспех"
		case "chain16":
			d.AC = 18
			d.ArmorNote = "Кольчуга + щит"
		case "scale14":
			d.AC = 14 + min(d.Mods["dex"], 2) + 2
			d.ArmorNote = "Чешуйчатый доспех + щит"
		case "barb":
			d.AC = 10 + d.Mods["dex"] + d.Mods["con"]
			d.ArmorNote = "Без доспеха"
		case "monk":
			d.AC = 10 + d.Mods["dex"] + d.Mods["wis"]
			d.ArmorNote = "Без доспеха"
		default:
			d.AC = 10 + d.Mods["dex"]
			d.ArmorNote = "Без доспеха"
		}
		if cl.Name == "Друид" {
			d.AC += 2
			d.ArmorNote += " + щит"
		}
		for _, w := range cl.Weapons {
			if w.Damage == "—" {
				continue
			}
			m := d.Mods["str"]
			if w.Ranged || (w.Finesse && d.Mods["dex"] > m) {
				m = d.Mods["dex"]
			}
			d.Attacks = append(d.Attacks, Attack5e{Name: w.Name, Bonus: m + d.Prof, Damage: fmt.Sprintf("%s %+d", w.Damage, m)})
		}
		if cl.Caster != "" {
			d.Caster = abilityRu[cl.Caster]
			d.SpellDC = 8 + d.Prof + d.Mods[cl.Caster]
			d.SpellAtk = d.Prof + d.Mods[cl.Caster]
		}
	} else {
		d.AC = 10 + d.Mods["dex"]
		d.MaxHP = 8 + d.Mods["con"]
	}
	for _, s := range skills5e {
		b := d.Mods[s.Ability]
		if contains(c.Skills, s.Key) {
			b += d.Prof
		}
		d.SkillBonus[s.Key] = b
	}
	d.Initiative = d.Mods["dex"]
	d.Passive = 10 + d.SkillBonus["perception"]
	return d
}

func (e DnD5e) Derive(raw json.RawMessage) (any, error) {
	var c Char5e
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}
	return e.derive(c), nil
}

func (e DnD5e) PDF(raw json.RawMessage, w io.Writer) error {
	var c Char5e
	if err := json.Unmarshal(raw, &c); err != nil {
		return err
	}
	return render5e(c, e.derive(c), w)
}
