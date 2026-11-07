package embed

import "testing"

func TestParse(t *testing.T) {
	if n := len(DnD5e().Races); n == 0 {
		t.Fatal("no races")
	}
	if DnD5e().Alignments[0].Desc == "" || Shadowdark().Backgrounds[0].Desc == "" {
		t.Fatal("no descriptions")
	}
	if DnD5e().AbilityDesc["str"] == "" || Shadowdark().AbilityDesc["cha"] == "" {
		t.Fatal("no ability descriptions")
	}
	for _, c := range Shadowdark().Classes {
		if len(c.Talents) != 11 {
			t.Fatalf("%s: %d talents", c.Name, len(c.Talents))
		}
	}
}
