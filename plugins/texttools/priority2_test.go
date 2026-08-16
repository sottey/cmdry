package texttools

import "testing"

func TestPriority2TextTools(t *testing.T) {
	if got := RotateText("abcd", 1); got != "bcda" {
		t.Fatalf("rotate = %q", got)
	}
	if got := RotateText("a🙂b", -1); got != "ba🙂" {
		t.Fatalf("unicode rotate = %q", got)
	}
	morse, err := ToMorse("SOS 2")
	if err != nil || morse != "... --- ... / ..---" {
		t.Fatalf("morse=%q err=%v", morse, err)
	}
	censored, count, err := CensorText("Bad badger, bad!", "bad", "[x]")
	if err != nil || censored != "[x] badger, [x]!" || count != 2 {
		t.Fatalf("censored=%q count=%d err=%v", censored, count, err)
	}
	quoted, count, err := QuoteText("Ada\n\nLin", "curly")
	if err != nil || quoted != "“Ada”\n“Lin”" || count != 2 {
		t.Fatalf("quoted=%q count=%d err=%v", quoted, count, err)
	}
}
