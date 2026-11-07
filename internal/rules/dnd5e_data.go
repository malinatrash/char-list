package rules

var abilities = []string{"str", "dex", "con", "int", "wis", "cha"}

var abilityRu = map[string]string{"str": "Сила", "dex": "Ловкость", "con": "Телосложение", "int": "Интеллект", "wis": "Мудрость", "cha": "Харизма"}

type skill5e struct {
	Key, Name, Ability string
}

var skills5e = []skill5e{
	{"acrobatics", "Акробатика", "dex"},
	{"investigation", "Анализ", "int"},
	{"athletics", "Атлетика", "str"},
	{"perception", "Внимательность", "wis"},
	{"survival", "Выживание", "wis"},
	{"performance", "Выступление", "cha"},
	{"intimidation", "Запугивание", "cha"},
	{"history", "История", "int"},
	{"sleight", "Ловкость рук", "dex"},
	{"arcana", "Магия", "int"},
	{"medicine", "Медицина", "wis"},
	{"deception", "Обман", "cha"},
	{"nature", "Природа", "int"},
	{"insight", "Проницательность", "wis"},
	{"religion", "Религия", "int"},
	{"stealth", "Скрытность", "dex"},
	{"persuasion", "Убеждение", "cha"},
	{"animal", "Уход за животными", "wis"},
}

type race5e struct {
	Name      string         `json:"name"`
	Bonus     map[string]int `json:"bonus"`
	Speed     int            `json:"speed"`
	Languages []string       `json:"languages"`
	Traits    []string       `json:"traits"`
	Syllables []string       `json:"-"`
	Flex      int            `json:"flex"`
}

var races5e = []race5e{
	{Name: "Человек", Bonus: map[string]int{"str": 1, "dex": 1, "con": 1, "int": 1, "wis": 1, "cha": 1}, Speed: 30, Languages: []string{"Общий", "любой один"}, Traits: []string{"Универсальность: +1 ко всем характеристикам"}, Syllables: []string{"Ал", "Бер", "Гар", "Дан", "Эл", "Ив", "Кас", "Мар", "Ро", "Та", "ин", "ар", "ик", "он", "ель", "ис", "ан"}},
	{Name: "Эльф", Bonus: map[string]int{"dex": 2}, Speed: 30, Languages: []string{"Общий", "Эльфийский"}, Traits: []string{"Тёмное зрение 60 фт", "Обострённые чувства: владение Внимательностью", "Наследие фей: преимущество на спасброски от очарования, иммунитет к усыплению", "Транс: 4 часа вместо сна"}, Syllables: []string{"Аэ", "Гал", "Ил", "Лаэ", "Тар", "Эри", "Фен", "Сил", "ион", "адриэль", "анис", "иэль", "ор", "уэн", "ас"}},
	{Name: "Дварф", Bonus: map[string]int{"con": 2}, Speed: 25, Languages: []string{"Общий", "Дварфийский"}, Traits: []string{"Тёмное зрение 60 фт", "Дварфийская устойчивость: преимущество против яда, сопротивление урону ядом", "Владение боевыми топорами, ручными топорами, лёгкими и боевыми молотами", "Знание камня"}, Syllables: []string{"Тор", "Бал", "Дур", "Гим", "Брун", "Ор", "Кар", "Харг", "ин", "ак", "гар", "дин", "ли", "ур", "ок"}},
	{Name: "Полурослик", Bonus: map[string]int{"dex": 2}, Speed: 25, Languages: []string{"Общий", "Полуросликов"}, Traits: []string{"Везучий: переброс 1 на к20", "Храбрый: преимущество против испуга", "Проворство полуросликов: можно проходить через пространство больших существ"}, Syllables: []string{"Бил", "Мер", "Пип", "Ро", "Сэм", "Фро", "Лил", "и", "бо", "до", "ри", "ин", "ка", "ли"}},
	{Name: "Гном", Bonus: map[string]int{"int": 2}, Speed: 25, Languages: []string{"Общий", "Гномий"}, Traits: []string{"Тёмное зрение 60 фт", "Гномья хитрость: преимущество на спасброски Инт, Мдр, Хар от магии"}, Syllables: []string{"Ал", "Бод", "Дим", "Эл", "Фон", "Гим", "Зук", "Нис", "ли", "бл", "дон", "ин", "ик", "ел", "ус"}},
	{Name: "Полуэльф", Bonus: map[string]int{"cha": 2}, Flex: 2, Speed: 30, Languages: []string{"Общий", "Эльфийский", "любой один"}, Traits: []string{"Тёмное зрение 60 фт", "Наследие фей", "Универсальность: владение двумя навыками на выбор"}, Syllables: []string{"Ар", "Эл", "Кас", "Лир", "Мар", "Таэ", "Вен", "иан", "ис", "ель", "он", "ра", "ен"}},
	{Name: "Полуорк", Bonus: map[string]int{"str": 2, "con": 1}, Speed: 30, Languages: []string{"Общий", "Орочий"}, Traits: []string{"Тёмное зрение 60 фт", "Угрожающий: владение Запугиванием", "Непоколебимая стойкость: раз в день вместо 0 хитов остаётся 1", "Свирепые атаки: лишняя кость урона при крите"}, Syllables: []string{"Гор", "Кра", "Тор", "Ург", "Дре", "Хог", "Шам", "ак", "ог", "аш", "ук", "мог", "ар"}},
	{Name: "Тифлинг", Bonus: map[string]int{"cha": 2, "int": 1}, Speed: 30, Languages: []string{"Общий", "Инфернальный"}, Traits: []string{"Тёмное зрение 60 фт", "Адское сопротивление: сопротивление урону огнём", "Дьявольское наследие: заговор чудотворство"}, Syllables: []string{"Ак", "Дам", "Ка", "Лер", "Мор", "Нем", "Сар", "ис", "ай", "рион", "ат", "ея", "ин", "ос"}},
	{Name: "Драконорождённый", Bonus: map[string]int{"str": 2, "cha": 1}, Speed: 30, Languages: []string{"Общий", "Драконий"}, Traits: []string{"Драконье наследие (выбрать цвет)", "Оружие дыхания 2к6, спасбросок КС 8 + мод Тел + БМ", "Сопротивление урону своего наследия"}, Syllables: []string{"Ар", "Бал", "Дон", "Гар", "Кри", "Мед", "Нар", "Тор", "джан", "асар", "аш", "рах", "ин", "ун", "ил"}},
}

type weapon5e struct {
	Name    string `json:"name"`
	Damage  string `json:"damage"`
	Ranged  bool   `json:"ranged"`
	Finesse bool   `json:"finesse"`
}

type class5e struct {
	Name       string     `json:"name"`
	HitDie     int        `json:"hitDie"`
	Saves      []string   `json:"saves"`
	SkillCount int        `json:"skillCount"`
	SkillPool  []string   `json:"skillPool"`
	Armor      string     `json:"armor"`
	ArmorAC    string     `json:"-"`
	Weapons    []weapon5e `json:"weapons"`
	Equipment  []string   `json:"equipment"`
	Features   []string   `json:"features"`
	Caster     string     `json:"caster"`
	Cantrips   []string   `json:"cantrips"`
	Spells     []string   `json:"spells"`
	NCantrips  int        `json:"nCantrips"`
	NSpells    int        `json:"nSpells"`
}

var allSkillKeys = func() []string {
	var ks []string
	for _, s := range skills5e {
		ks = append(ks, s.Key)
	}
	return ks
}()

var classes5e = []class5e{
	{Name: "Бард", HitDie: 8, Saves: []string{"dex", "cha"}, SkillCount: 3, SkillPool: allSkillKeys, Armor: "Кожаный доспех", ArmorAC: "light11",
		Weapons:   []weapon5e{{Name: "Рапира", Damage: "1к8 колющий", Finesse: true}, {Name: "Кинжал", Damage: "1к4 колющий", Finesse: true}},
		Equipment: []string{"Кожаный доспех", "Рапира", "Кинжал", "Лютня", "Набор артиста"}, Features: []string{"Вдохновение барда (к6)", "Использование заклинаний"},
		Caster: "cha", NCantrips: 2, NSpells: 4, Cantrips: []string{"Злая насмешка", "Малая иллюзия", "Фокусы", "Свет", "Дружба", "Пляшущие огоньки"}, Spells: []string{"Лечение ран", "Очарование личности", "Обнаружение магии", "Волна грома", "Смех Таши", "Маскировка", "Усыпление", "Героизм"}},
	{Name: "Варвар", HitDie: 12, Saves: []string{"str", "con"}, SkillCount: 2, SkillPool: []string{"animal", "athletics", "intimidation", "nature", "perception", "survival"}, Armor: "Без доспеха", ArmorAC: "barb",
		Weapons:   []weapon5e{{Name: "Секира", Damage: "1к12 рубящий"}, {Name: "Ручной топор", Damage: "1к6 рубящий"}, {Name: "Метательное копьё", Damage: "1к6 колющий"}},
		Equipment: []string{"Секира", "Два ручных топора", "Четыре метательных копья", "Набор путешественника"}, Features: []string{"Ярость (2 раза, +2 к урону)", "Защита без доспехов (10 + Лов + Тел)"}},
	{Name: "Воин", HitDie: 10, Saves: []string{"str", "con"}, SkillCount: 2, SkillPool: []string{"acrobatics", "animal", "athletics", "history", "insight", "intimidation", "perception", "survival"}, Armor: "Кольчуга", ArmorAC: "chain16",
		Weapons:   []weapon5e{{Name: "Длинный меч", Damage: "1к8 рубящий"}, {Name: "Лёгкий арбалет", Damage: "1к8 колющий", Ranged: true}},
		Equipment: []string{"Кольчуга", "Длинный меч", "Щит", "Лёгкий арбалет и 20 болтов", "Набор исследователя подземелий"}, Features: []string{"Боевой стиль (выбрать)", "Второе дыхание (1к10 + уровень)"}},
	{Name: "Волшебник", HitDie: 6, Saves: []string{"int", "wis"}, SkillCount: 2, SkillPool: []string{"arcana", "history", "insight", "investigation", "medicine", "religion"}, Armor: "Без доспеха", ArmorAC: "none",
		Weapons:   []weapon5e{{Name: "Боевой посох", Damage: "1к6 дробящий"}, {Name: "Кинжал", Damage: "1к4 колющий", Finesse: true}},
		Equipment: []string{"Боевой посох", "Кинжал", "Мешочек с компонентами", "Книга заклинаний", "Набор учёного"}, Features: []string{"Использование заклинаний", "Магическое восстановление"},
		Caster: "int", NCantrips: 3, NSpells: 6, Cantrips: []string{"Огненный снаряд", "Фокусы", "Свет", "Волшебная рука", "Луч холода", "Электрошок"}, Spells: []string{"Волшебная стрела", "Щит", "Усыпление", "Обнаружение магии", "Доспехи мага", "Волна грома", "Очарование личности", "Туманное облако", "Понимание языков"}},
	{Name: "Друид", HitDie: 8, Saves: []string{"int", "wis"}, SkillCount: 2, SkillPool: []string{"arcana", "animal", "insight", "medicine", "nature", "perception", "religion", "survival"}, Armor: "Кожаный доспех", ArmorAC: "light11",
		Weapons:   []weapon5e{{Name: "Скимитар", Damage: "1к6 рубящий", Finesse: true}, {Name: "Деревянный щит", Damage: "—"}},
		Equipment: []string{"Кожаный доспех", "Деревянный щит", "Скимитар", "Набор исследователя", "Друидический фокус"}, Features: []string{"Друидический язык", "Использование заклинаний"},
		Caster: "wis", NCantrips: 2, NSpells: 3, Cantrips: []string{"Друидизм", "Дубинка", "Сотворение пламени", "Указание"}, Spells: []string{"Лечение ран", "Опутывание", "Дружба с животными", "Волна грома", "Доброе ягоды", "Туманное облако", "Разговор с животными"}},
	{Name: "Жрец", HitDie: 8, Saves: []string{"wis", "cha"}, SkillCount: 2, SkillPool: []string{"history", "insight", "medicine", "persuasion", "religion"}, Armor: "Чешуйчатый доспех", ArmorAC: "scale14",
		Weapons:   []weapon5e{{Name: "Булава", Damage: "1к6 дробящий"}, {Name: "Лёгкий арбалет", Damage: "1к8 колющий", Ranged: true}},
		Equipment: []string{"Чешуйчатый доспех", "Булава", "Щит", "Лёгкий арбалет и 20 болтов", "Священный символ", "Набор священника"}, Features: []string{"Божественный домен (выбрать)", "Использование заклинаний"},
		Caster: "wis", NCantrips: 3, NSpells: 3, Cantrips: []string{"Священное пламя", "Указание", "Сопротивление", "Свет", "Искусство", "Чудотворство"}, Spells: []string{"Лечение ран", "Благословение", "Порча", "Направленный снаряд", "Щит веры", "Обнаружение зла и добра", "Убежище", "Очищение пищи и питья"}},
	{Name: "Колдун", HitDie: 8, Saves: []string{"wis", "cha"}, SkillCount: 2, SkillPool: []string{"arcana", "deception", "history", "intimidation", "investigation", "nature", "religion"}, Armor: "Кожаный доспех", ArmorAC: "light11",
		Weapons:   []weapon5e{{Name: "Кинжал", Damage: "1к4 колющий", Finesse: true}, {Name: "Лёгкий арбалет", Damage: "1к8 колющий", Ranged: true}},
		Equipment: []string{"Кожаный доспех", "Два кинжала", "Лёгкий арбалет и 20 болтов", "Мистический фокус", "Набор учёного"}, Features: []string{"Потусторонний покровитель (выбрать)", "Магия договора"},
		Caster: "cha", NCantrips: 2, NSpells: 2, Cantrips: []string{"Мистический заряд", "Фокусы", "Малая иллюзия", "Дружба", "Ядовитые брызги"}, Spells: []string{"Адское возмездие", "Проклятие", "Доспехи Агатиса", "Очарование личности", "Ведьмин снаряд", "Защита от зла и добра"}},
	{Name: "Монах", HitDie: 8, Saves: []string{"str", "dex"}, SkillCount: 2, SkillPool: []string{"acrobatics", "athletics", "history", "insight", "religion", "stealth"}, Armor: "Без доспеха", ArmorAC: "monk",
		Weapons:   []weapon5e{{Name: "Короткий меч", Damage: "1к6 колющий", Finesse: true}, {Name: "Безоружный удар", Damage: "1к4 дробящий", Finesse: true}, {Name: "Дротик", Damage: "1к4 колющий", Finesse: true}},
		Equipment: []string{"Короткий меч", "10 дротиков", "Набор исследователя подземелий"}, Features: []string{"Защита без доспехов (10 + Лов + Мдр)", "Боевые искусства (к4)"}},
	{Name: "Паладин", HitDie: 10, Saves: []string{"wis", "cha"}, SkillCount: 2, SkillPool: []string{"athletics", "insight", "intimidation", "medicine", "persuasion", "religion"}, Armor: "Кольчуга", ArmorAC: "chain16",
		Weapons:   []weapon5e{{Name: "Длинный меч", Damage: "1к8 рубящий"}, {Name: "Метательное копьё", Damage: "1к6 колющий"}},
		Equipment: []string{"Кольчуга", "Длинный меч", "Щит", "Пять метательных копий", "Священный символ", "Набор священника"}, Features: []string{"Божественное чувство", "Наложение рук (5 хитов)"}},
	{Name: "Плут", HitDie: 8, Saves: []string{"dex", "int"}, SkillCount: 4, SkillPool: []string{"acrobatics", "athletics", "deception", "insight", "intimidation", "investigation", "perception", "performance", "persuasion", "sleight", "stealth"}, Armor: "Кожаный доспех", ArmorAC: "light11",
		Weapons:   []weapon5e{{Name: "Рапира", Damage: "1к8 колющий", Finesse: true}, {Name: "Короткий лук", Damage: "1к6 колющий", Ranged: true}, {Name: "Кинжал", Damage: "1к4 колющий", Finesse: true}},
		Equipment: []string{"Кожаный доспех", "Рапира", "Короткий лук и колчан с 20 стрелами", "Два кинжала", "Воровские инструменты", "Набор взломщика"}, Features: []string{"Компетентность (два навыка ×2)", "Скрытая атака 1к6", "Воровской жаргон"}},
	{Name: "Следопыт", HitDie: 10, Saves: []string{"str", "dex"}, SkillCount: 3, SkillPool: []string{"animal", "athletics", "insight", "investigation", "nature", "perception", "stealth", "survival"}, Armor: "Чешуйчатый доспех", ArmorAC: "scale14",
		Weapons:   []weapon5e{{Name: "Длинный лук", Damage: "1к8 колющий", Ranged: true}, {Name: "Короткий меч", Damage: "1к6 колющий", Finesse: true}},
		Equipment: []string{"Чешуйчатый доспех", "Два коротких меча", "Длинный лук и колчан с 20 стрелами", "Набор исследователя"}, Features: []string{"Избранный враг (выбрать)", "Исследователь природы (выбрать местность)"}},
	{Name: "Чародей", HitDie: 6, Saves: []string{"con", "cha"}, SkillCount: 2, SkillPool: []string{"arcana", "deception", "insight", "intimidation", "persuasion", "religion"}, Armor: "Без доспеха", ArmorAC: "none",
		Weapons:   []weapon5e{{Name: "Лёгкий арбалет", Damage: "1к8 колющий", Ranged: true}, {Name: "Кинжал", Damage: "1к4 колющий", Finesse: true}},
		Equipment: []string{"Лёгкий арбалет и 20 болтов", "Два кинжала", "Магический фокус", "Набор исследователя подземелий"}, Features: []string{"Происхождение чародея (выбрать)", "Использование заклинаний"},
		Caster: "cha", NCantrips: 4, NSpells: 2, Cantrips: []string{"Огненный снаряд", "Фокусы", "Свет", "Волшебная рука", "Луч холода", "Электрошок", "Дружба"}, Spells: []string{"Волшебная стрела", "Щит", "Усыпление", "Хроматический шар", "Обнаружение магии", "Волна грома", "Очарование личности"}},
}

type background5e struct {
	Name      string   `json:"name"`
	Skills    []string `json:"skills"`
	Extra     string   `json:"extra"`
	Equipment []string `json:"equipment"`
	Feature   string   `json:"feature"`
}

var backgrounds5e = []background5e{
	{"Аколит", []string{"insight", "religion"}, "Два любых языка", []string{"Священный символ", "Молитвенник", "5 палочек благовоний", "Одежда", "15 зм"}, "Приют верующих"},
	{"Преступник", []string{"deception", "stealth"}, "Воровские инструменты, игровой набор", []string{"Ломик", "Тёмная одежда с капюшоном", "15 зм"}, "Преступные связи"},
	{"Народный герой", []string{"animal", "survival"}, "Ремесленные инструменты, наземные средства передвижения", []string{"Ремесленные инструменты", "Лопата", "Железный горшок", "Одежда", "10 зм"}, "Деревенское гостеприимство"},
	{"Благородный", []string{"history", "persuasion"}, "Игровой набор, один язык", []string{"Дорогая одежда", "Перстень с печатью", "Свиток родословной", "25 зм"}, "Положение в обществе"},
	{"Мудрец", []string{"arcana", "history"}, "Два любых языка", []string{"Бутылочка чернил", "Перо", "Нож", "Письмо от коллеги", "Одежда", "10 зм"}, "Исследователь"},
	{"Солдат", []string{"athletics", "intimidation"}, "Игровой набор, наземные средства передвижения", []string{"Знак отличия", "Трофей с поверженного врага", "Кости", "Одежда", "10 зм"}, "Военное звание"},
	{"Отшельник", []string{"medicine", "religion"}, "Набор травника, один язык", []string{"Футляр для свитков", "Одеяло", "Одежда", "Набор травника", "5 зм"}, "Открытие"},
	{"Артист", []string{"acrobatics", "performance"}, "Набор для грима, музыкальный инструмент", []string{"Музыкальный инструмент", "Подарок поклонника", "Костюм", "15 зм"}, "Популярность"},
	{"Чужеземец", []string{"athletics", "survival"}, "Музыкальный инструмент, один язык", []string{"Посох", "Капкан", "Трофей с убитого животного", "Одежда путешественника", "10 зм"}, "Странник"},
	{"Гильдейский ремесленник", []string{"insight", "persuasion"}, "Ремесленные инструменты, один язык", []string{"Ремесленные инструменты", "Рекомендательное письмо гильдии", "Одежда путешественника", "15 зм"}, "Членство в гильдии"},
	{"Моряк", []string{"athletics", "perception"}, "Навигационные инструменты, водные средства передвижения", []string{"Дубинка", "50 футов шёлковой верёвки", "Талисман", "Одежда", "10 зм"}, "Проезд на корабле"},
	{"Шарлатан", []string{"deception", "sleight"}, "Набор для грима, набор для фальсификации", []string{"Дорогая одежда", "Набор для грима", "Инструменты мошенника", "15 зм"}, "Фальшивая личность"},
	{"Беспризорник", []string{"sleight", "stealth"}, "Набор для грима, воровские инструменты", []string{"Небольшой нож", "Карта родного города", "Ручная мышь", "Памятная вещь родителей", "Одежда", "10 зм"}, "Городские тайны"},
}

var alignments5e = []string{"Законно-добрый", "Нейтрально-добрый", "Хаотично-добрый", "Законно-нейтральный", "Нейтральный", "Хаотично-нейтральный", "Законно-злой", "Нейтрально-злой", "Хаотично-злой"}

var traits5e = []string{
	"Я всегда ищу повод для доброй драки.", "Я обожаю рассказывать истории о своих подвигах, приукрашивая их.", "Ничто не выведет меня из равновесия.", "Я говорю мало, но слушаю внимательно.",
	"Я сыплю цитатами из древних текстов, часто невпопад.", "Я сужу о людях по тому, как они обращаются со своим оружием.", "Я не могу удержаться, чтобы не взять что-нибудь блестящее.", "Мне неловко в обществе знати.",
	"Я легко завожу друзей, где бы ни оказался.", "Я всегда ем так, словно это мой последний ужин.", "Я предпочитаю решать всё по-своему и ненавижу приказы.", "Я собираю слухи, как другие собирают монеты.",
}

var ideals5e = []string{
	"Свобода. Цепи предназначены для того, чтобы их разбивать. (Хаотичный)", "Традиция. Древние обычаи нужно хранить. (Законный)", "Милосердие. Я помогаю тем, кто не может помочь себе. (Добрый)",
	"Власть. Однажды я буду править. (Злой)", "Знание. Путь к силе лежит через понимание. (Нейтральный)", "Честь. Я никогда не нарушу данного слова. (Законный)", "Перемены. Мы должны помогать переменам, что несёт судьба. (Хаотичный)",
	"Величие. Меня будут помнить. (Любой)", "Народ. Я защищаю тех, с кем живу. (Нейтральный)",
}

var bonds5e = []string{
	"Я отдам жизнь за своих товарищей по оружию.", "Меня ждёт долг, который я должен выплатить.", "Я ищу того, кто разрушил мою семью.", "Мой учитель погиб, и я должен завершить его дело.",
	"Где-то есть артефакт, который принадлежит мне по праву.", "Я защищаю родную деревню, даже когда далеко от неё.", "Моя честь — это моя жизнь.", "Я должен доказать, что достоин своего имени.",
}

var flaws5e = []string{
	"Я не доверяю никому, кроме себя.", "Выпивка и азарт тянут меня на дно.", "Я не могу устоять перед красивым лицом.", "Я слишком уверен в себе и часто лезу на рожон.",
	"Тайна, которую я храню, однажды погубит меня.", "Я презираю трусов и часто говорю об этом вслух.", "Я никогда не признаю своих ошибок.", "Ради золота я готов на многое.",
}
