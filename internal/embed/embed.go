// Package embed holds everything the binary carries with it: tables of the editions, sheet templates, fonts, pages and static files.
package embed

import (
	goembed "embed"
	"io/fs"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed editions fonts sheets templates static
var files goembed.FS

type Skill5e struct {
	Key     string `json:"key" yaml:"key"`
	Name    string `json:"name" yaml:"name"`
	Ability string `json:"ability" yaml:"ability"`
}

type Race5e struct {
	Name      string         `json:"name" yaml:"name"`
	Bonus     map[string]int `json:"bonus" yaml:"bonus"`
	Flex      int            `json:"flex" yaml:"flex"`
	Speed     int            `json:"speed" yaml:"speed"`
	Languages []string       `json:"languages" yaml:"languages"`
	Traits    []string       `json:"traits" yaml:"traits"`
	Syllables []string       `json:"-" yaml:"syllables"`
}

type Weapon5e struct {
	Name    string `json:"name" yaml:"name"`
	Damage  string `json:"damage" yaml:"damage"`
	Ranged  bool   `json:"ranged" yaml:"ranged"`
	Finesse bool   `json:"finesse" yaml:"finesse"`
}

type Class5e struct {
	Name       string     `json:"name" yaml:"name"`
	HitDie     int        `json:"hitDie" yaml:"hitDie"`
	Saves      []string   `json:"saves" yaml:"saves"`
	SkillCount int        `json:"skillCount" yaml:"skillCount"`
	SkillPool  []string   `json:"skillPool" yaml:"skillPool"`
	Armor      string     `json:"armor" yaml:"armor"`
	ArmorAC    string     `json:"-" yaml:"armorAC"`
	Weapons    []Weapon5e `json:"weapons" yaml:"weapons"`
	Equipment  []string   `json:"equipment" yaml:"equipment"`
	Features   []string   `json:"features" yaml:"features"`
	Caster     string     `json:"caster" yaml:"caster"`
	Cantrips   []string   `json:"cantrips" yaml:"cantrips"`
	Spells     []string   `json:"spells" yaml:"spells"`
	NCantrips  int        `json:"nCantrips" yaml:"nCantrips"`
	NSpells    int        `json:"nSpells" yaml:"nSpells"`
}

type Background5e struct {
	Name      string   `json:"name" yaml:"name"`
	Skills    []string `json:"skills" yaml:"skills"`
	Extra     string   `json:"extra" yaml:"extra"`
	Equipment []string `json:"equipment" yaml:"equipment"`
	Feature   string   `json:"feature" yaml:"feature"`
}

type DnD5eData struct {
	Skills      []Skill5e      `json:"skills" yaml:"skills"`
	Races       []Race5e       `json:"races" yaml:"races"`
	Classes     []Class5e      `json:"classes" yaml:"classes"`
	Backgrounds []Background5e `json:"backgrounds" yaml:"backgrounds"`
	Alignments  []string       `json:"alignments" yaml:"alignments"`
	Traits      []string       `json:"traits" yaml:"traits"`
	Ideals      []string       `json:"ideals" yaml:"ideals"`
	Bonds       []string       `json:"bonds" yaml:"bonds"`
	Flaws       []string       `json:"flaws" yaml:"flaws"`
}

type AncestrySD struct {
	Name      string   `json:"name" yaml:"name"`
	Trait     string   `json:"trait" yaml:"trait"`
	Languages string   `json:"languages" yaml:"languages"`
	Syllables []string `json:"-" yaml:"syllables"`
}

type WeaponSD struct {
	Name   string `json:"name" yaml:"name"`
	Damage string `json:"damage" yaml:"damage"`
	Ranged bool   `json:"ranged" yaml:"ranged,omitempty"`
	Cost   int    `json:"cost" yaml:"cost"`
}

type ClassSD struct {
	Name     string         `json:"name" yaml:"name"`
	HitDie   int            `json:"hitDie" yaml:"hitDie"`
	Armor    string         `json:"armor" yaml:"armor"`
	Weapons  string         `json:"weapons" yaml:"weapons"`
	Features []string       `json:"features" yaml:"features"`
	Titles   [3]string      `json:"titles" yaml:"titles"`
	Talents  map[int]string `json:"talents" yaml:"talents"`
	Spells   []string       `json:"spells" yaml:"spells"`
	NSpells  int            `json:"nSpells" yaml:"nSpells"`
	Kit      []WeaponSD     `json:"-" yaml:"kit"`
}

type ShadowdarkData struct {
	Ancestries  []AncestrySD `json:"ancestries" yaml:"ancestries"`
	Classes     []ClassSD    `json:"classes" yaml:"classes"`
	Alignments  []string     `json:"alignments" yaml:"alignments"`
	Deities     []string     `json:"deities" yaml:"deities"`
	Backgrounds []string     `json:"backgrounds" yaml:"backgrounds"`
	Gear        []string     `json:"gear" yaml:"gear"`
}

var (
	dnd5e      = sync.OnceValue(func() *DnD5eData { return parse[DnD5eData]("editions/dnd5e.yaml") })
	shadowdark = sync.OnceValue(func() *ShadowdarkData { return parse[ShadowdarkData]("editions/shadowdark.yaml") })
)

func DnD5e() *DnD5eData           { return dnd5e() }
func Shadowdark() *ShadowdarkData { return shadowdark() }

func Sheet(name string) []byte { return must(files.ReadFile("sheets/" + name)) }
func Font(name string) []byte  { return must(files.ReadFile("fonts/" + name)) }
func Templates() fs.FS         { return sub("templates") }
func Static() fs.FS            { return sub("static") }

func parse[T any](path string) *T {
	var v T
	if err := yaml.Unmarshal(must(files.ReadFile(path)), &v); err != nil {
		panic(path + ": " + err.Error())
	}
	return &v
}

func sub(dir string) fs.FS { return must(fs.Sub(files, dir)) }

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
