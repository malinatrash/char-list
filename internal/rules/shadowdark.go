package rules

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"char-list/internal/dice"
)

type ancestrySD struct {
	Name      string   `json:"name"`
	Trait     string   `json:"trait"`
	Languages string   `json:"languages"`
	Syllables []string `json:"-"`
}

var ancestriesSD = []ancestrySD{
	{"Человек", "Честолюбие: дополнительный талант на 1 уровне", "Общий, +1 любой", []string{"Ар", "Бел", "Кор", "Дар", "Эд", "Гал", "Ив", "Мар", "Ро", "Сен", "ин", "ик", "ан", "ель", "ис", "ор"}},
	{"Эльф", "Зоркость: +1 к дальним атакам или заклинаниям", "Общий, эльфийский, сильванский", []string{"Аэ", "Эл", "Ил", "Лаэ", "Сил", "Тар", "Фен", "ион", "иэль", "ор", "уэн", "ас", "анис"}},
	{"Дварф", "Крепость: +2 хита, хиты бросаются с преимуществом", "Общий, дварфийский", []string{"Тор", "Бал", "Дур", "Гим", "Брун", "Кар", "Хар", "ин", "ак", "гар", "дин", "ур", "ок"}},
	{"Полурослик", "Скрытность: раз в день стать невидимым на 3 раунда", "Общий", []string{"Бил", "Мер", "Пип", "Ро", "Сэм", "Лил", "Том", "и", "бо", "до", "ри", "ин", "ка"}},
	{"Полуорк", "Мощь: +1 к атаке и урону в ближнем бою", "Общий, орочий", []string{"Гор", "Кра", "Ург", "Дре", "Хог", "Шам", "Бру", "ак", "ог", "аш", "ук", "мог", "ар"}},
	{"Гоблин", "Чуткость: нельзя застать врасплох", "Общий, гоблинский", []string{"Зик", "Ник", "Гриз", "Снаг", "Мик", "Пог", "Рат", "к", "ли", "ит", "аз", "ок", "ль"}},
}

type weaponSD struct {
	Name, Damage string
	Ranged       bool
	Cost         int
}

type classSD struct {
	Name      string     `json:"name"`
	HitDie    int        `json:"hitDie"`
	Armor     string     `json:"armor"`
	Weapons   string     `json:"weapons"`
	Features  []string   `json:"features"`
	Titles    [3]string  `json:"titles"`
	Talents   []string   `json:"talents"`
	Spells    []string   `json:"spells"`
	NSpells   int        `json:"nSpells"`
	Kit       []weaponSD `json:"-"`
	BuysArmor bool       `json:"-"`
}

var classesSD = []classSD{
	{Name: "Воин", HitDie: 8, Armor: "Любые доспехи и щиты", Weapons: "Любое оружие",
		Features: []string{"Носильщик: +Тел слотов снаряжения (мин. 1)", "Мастерство оружия: +1 к атаке и урону одним оружием, +1/2 уровня", "Стойкость: преимущество на проверки Силы или Ловкости (выбрать)"},
		Titles:   [3]string{"Оруженосец", "Воитель", "Плут"},
		Talents:  []string{"", "", "Мастерство ещё одним типом оружия", "+1 к атакам в ближнем и дальнем бою", "+1 к атакам в ближнем и дальнем бою", "+1 к атакам в ближнем и дальнем бою", "+1 к атакам в ближнем и дальнем бою", "+2 к Силе, Ловкости или Телосложению", "+2 к Силе, Ловкости или Телосложению", "+2 к Силе, Ловкости или Телосложению", "+1 КД в выбранном типе доспеха", "+1 КД в выбранном типе доспеха", "Выбрать любой талант или +2 к характеристикам"},
		Kit:      []weaponSD{{"Длинный меч", "1к8", false, 9}, {"Щит", "", false, 10}, {"Кожаный доспех", "", false, 10}}, BuysArmor: true},
	{Name: "Жрец", HitDie: 6, Armor: "Любые доспехи и щиты", Weapons: "Арбалет, кинжал, булава, длинный меч, посох, боевой молот",
		Features: []string{"Изгнание нежити: заклинание всегда известно", "Божественная магия: проверка Мудрости, КС 10 + круг"},
		Titles:   [3]string{"Послушник", "Искатель", "Посвящённый"},
		Talents:  []string{"", "", "Преимущество при касте одного известного заклинания", "+1 к атакам в ближнем бою или заклинаниями", "+1 к атакам в ближнем бою или заклинаниями", "+1 к атакам в ближнем бою или заклинаниями", "+1 к атакам в ближнем бою или заклинаниями", "+1 к проверкам заклинаний жреца", "+1 к проверкам заклинаний жреца", "+1 к проверкам заклинаний жреца", "+2 к Силе или Мудрости", "+2 к Силе или Мудрости", "Выбрать любой талант или +2 к характеристикам"},
		Spells:   []string{"Лечение ран", "Святое оружие", "Свет", "Защита от зла", "Щит веры"}, NSpells: 2,
		Kit: []weaponSD{{"Булава", "1к6", false, 5}, {"Щит", "", false, 10}, {"Кожаный доспех", "", false, 10}}, BuysArmor: true},
	{Name: "Вор", HitDie: 4, Armor: "Кожаный доспех, мифриловая кольчуга", Weapons: "Дубинка, арбалет, кинжал, короткий лук, короткий меч",
		Features: []string{"Удар в спину: +1 кость урона при преимуществе, +1 за каждые 2 уровня", "Воровство: преимущество на лазание, маскировку, поиск и обезвреживание ловушек, кражу, вскрытие замков"},
		Titles:   [3]string{"Воришка", "Грабитель", "Громила"},
		Talents:  []string{"", "", "Преимущество на броски инициативы", "+1 кость урона Удара в спину", "+1 кость урона Удара в спину", "+1 кость урона Удара в спину", "+2 к Силе, Ловкости или Харизме", "+2 к Силе, Ловкости или Харизме", "+2 к Силе, Ловкости или Харизме", "+2 к Силе, Ловкости или Харизме", "+1 к атакам в ближнем и дальнем бою", "+1 к атакам в ближнем и дальнем бою", "Выбрать любой талант или +2 к характеристикам"},
		Kit:      []weaponSD{{"Короткий меч", "1к6", false, 7}, {"Кинжал", "1к4", false, 1}, {"Кожаный доспех", "", false, 10}}, BuysArmor: true},
	{Name: "Волшебник", HitDie: 4, Armor: "Нет", Weapons: "Кинжал, посох",
		Features: []string{"Тайная магия: проверка Интеллекта, КС 10 + круг", "Изучение магии: можно учить заклинания со свитков"},
		Titles:   [3]string{"Ученик", "Шаман", "Адепт"},
		Talents:  []string{"", "", "Случайный магический предмет выбранного типа", "+2 к Интеллекту или +1 к проверкам заклинаний", "+2 к Интеллекту или +1 к проверкам заклинаний", "+2 к Интеллекту или +1 к проверкам заклинаний", "+2 к Интеллекту или +1 к проверкам заклинаний", "+2 к Интеллекту или +1 к проверкам заклинаний", "Преимущество при касте одного известного заклинания", "Преимущество при касте одного известного заклинания", "Выучить ещё одно заклинание волшебника", "Выучить ещё одно заклинание волшебника", "Выбрать любой талант или +2 к характеристикам"},
		Spells:   []string{"Сигнал тревоги", "Горящие руки", "Очарование личности", "Обнаружение магии", "Падение пёрышком", "Парящий диск", "Запереть врата", "Свет", "Доспехи мага", "Волшебная стрела", "Защита от зла", "Усыпление"}, NSpells: 3,
		Kit: []weaponSD{{"Посох", "1к4", false, 1}, {"Кинжал", "1к4", false, 1}}},
}

var alignmentsSD = []string{"Законный", "Нейтральный", "Хаотичный"}

var deitiesSD = []string{"Святой Террагнис (законный)", "Мадира Завет (законная)", "Гед (нейтральный)", "Орд (нейтральный)", "Потерянные (нейтральные)", "Мемнон (хаотичный)", "Шун Мерзкая (хаотичная)", "Рамлаат (хаотичный)"}

var backgroundsSD = []string{"Беспризорник", "В розыске", "Посвящённый культа", "Гильдия воров", "Изгнанник", "Сирота", "Ученик волшебника", "Ювелир", "Травник", "Варвар", "Наёмник", "Моряк", "Послушник", "Солдат", "Следопыт", "Разведчик", "Менестрель", "Учёный", "Дворянин", "Хирург"}

var gearSD = []string{"Факел", "Кинжал", "Шест", "Короткий лук и 5 стрел", "Верёвка 60 футов", "Фляга масла", "Лом", "Железные шипы (10)", "Кремень и огниво", "Кошка", "Дубинка", "Калтропы"}

type CharSD struct {
	Name       string         `json:"name"`
	Ancestry   string         `json:"ancestry"`
	Class      string         `json:"class"`
	Level      int            `json:"level"`
	XP         int            `json:"xp"`
	Alignment  string         `json:"alignment"`
	Title      string         `json:"title"`
	Background string         `json:"background"`
	Deity      string         `json:"deity"`
	Scores     map[string]int `json:"scores"`
	HP         int            `json:"hp"`
	Talents    string         `json:"talents"`
	Attacks    string         `json:"attacks"`
	Gear       string         `json:"gear"`
	GP         int            `json:"gp"`
	SP         int            `json:"sp"`
	CP         int            `json:"cp"`
}

type DerivedSD struct {
	Mods    map[string]int `json:"mods"`
	AC      int            `json:"ac"`
	ACNote  string         `json:"acNote"`
	Slots   int            `json:"slots"`
	XPNext  int            `json:"xpNext"`
	Title   string         `json:"title"`
	HitDie  string         `json:"hitDie"`
	Attacks []Attack5e     `json:"attacks"`
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
		"ancestries":  ancestriesSD,
		"classes":     classesSD,
		"alignments":  alignmentsSD,
		"deities":     deitiesSD,
		"backgrounds": backgroundsSD,
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

func ancestrySDByName(n string) *ancestrySD {
	for i := range ancestriesSD {
		if ancestriesSD[i].Name == n {
			return &ancestriesSD[i]
		}
	}
	return nil
}

func classSDByName(n string) *classSD {
	for i := range classesSD {
		if classesSD[i].Name == n {
			return &classesSD[i]
		}
	}
	return nil
}

func rollStatsSD() (map[string]int, []dice.Log) {
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

func titleSD(cl *classSD, alignment string) string {
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

func (e Shadowdark) Random(raw json.RawMessage, fields []string) (any, []dice.Log, error) {
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
	if has(fields, "ancestry") || ancestrySDByName(c.Ancestry) == nil {
		r := dice.Roll("Происхождение", 1, 12)
		logs = append(logs, r)
		idx := []int{0, 0, 0, 0, 1, 1, 2, 2, 3, 3, 4, 5}[r.Total-1]
		c.Ancestry = ancestriesSD[idx].Name
	}
	if has(fields, "class") || classSDByName(c.Class) == nil {
		r := dice.Roll("Класс", 1, 4)
		logs = append(logs, r)
		c.Class = classesSD[r.Total-1].Name
	}
	if has(fields, "alignment") {
		r := dice.Roll("Мировоззрение", 1, 6)
		logs = append(logs, r)
		c.Alignment = alignmentsSD[[]int{0, 0, 0, 1, 1, 2}[r.Total-1]]
	}
	if has(fields, "background") {
		r := dice.Roll("Предыстория", 1, 20)
		logs = append(logs, r)
		c.Background = backgroundsSD[r.Total-1]
	}
	if has(fields, "deity") {
		c.Deity = dice.Pick(deitiesSD)
	}
	if has(fields, "name") {
		c.Name = genName(ancestrySDByName(c.Ancestry).Syllables)
	}
	if has(fields, "scores") {
		var l []dice.Log
		c.Scores, l = rollStatsSD()
		logs = append(logs, l...)
	}
	cl := classSDByName(c.Class)
	anc := ancestrySDByName(c.Ancestry)
	if has(fields, "title") {
		c.Title = titleSD(cl, c.Alignment)
	}
	if has(fields, "hp") {
		con := modSD(c.Scores["con"])
		r := dice.Roll("Хиты", 1, cl.HitDie)
		logs = append(logs, r)
		hp := r.Total
		if anc.Name == "Дварф" {
			r2 := dice.Roll("Хиты (преимущество)", 1, cl.HitDie)
			logs = append(logs, r2)
			hp = max(hp, r2.Total) + 2
		}
		c.HP = max(hp+con, 1)
	}
	if has(fields, "talents") {
		n := 1
		if anc.Name == "Человек" {
			n = 2
		}
		var lines []string
		lines = append(lines, "Происхождение — "+anc.Trait)
		lines = append(lines, cl.Features...)
		for i := 0; i < n; i++ {
			r := dice.Roll("Талант", 2, 6)
			logs = append(logs, r)
			lines = append(lines, fmt.Sprintf("Талант (%d): %s", r.Total, cl.Talents[r.Total]))
		}
		if cl.NSpells > 0 {
			lines = append(lines, "Заклинания 1 круга: "+strings.Join(dice.PickN(cl.Spells, cl.NSpells), ", "))
		}
		c.Talents = strings.Join(lines, "\n")
	}
	if has(fields, "gear") {
		g := dice.Roll("Золото", 2, 6)
		logs = append(logs, g)
		c.GP = g.Total * 5
		c.SP, c.CP = 0, 0
		var items []string
		for _, k := range cl.Kit {
			if k.Cost <= c.GP {
				c.GP -= k.Cost
				items = append(items, k.Name)
			}
		}
		n := dice.Roll("Снаряжение", 1, 4)
		logs = append(logs, n)
		for i := 0; i < n.Total; i++ {
			r := dice.Roll("Предмет", 1, 12)
			logs = append(logs, r)
			items = append(items, gearSD[r.Total-1])
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
