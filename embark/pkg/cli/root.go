package cli

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/containerum/containerum/embark/help"

	"github.com/huandu/xstrings"

	"github.com/common-nighthawk/go-figure"
	"github.com/octago/sflags/gen/gpflag"
	"github.com/spf13/cobra"
)

type RootFlags struct{}

func Root() *cobra.Command {
	var flags RootFlags
	var cmd = &cobra.Command{
		Use: "embark",
		Run: func(cmd *cobra.Command, args []string) {
			var font = selectFont()
			var greetings = figure.NewFigure("Embark to containeum!", font, false).String()
			fmt.Print(greetings)
			fmt.Println(xstrings.Center(logo, width(greetings), ""))
		},
	}
	cmd.AddCommand(
		Install(nil),
		Download())
	help.FillHelps(cmd)
	if err := gpflag.ParseTo(&flags, cmd.PersistentFlags()); err != nil {
		panic(err)
	}
	return cmd
}

const logo = `
      ___________________
     / -----------------/| 
    / /              / / | 
   / /_____________ / / ||
  | ________________ |  ||
  | |      ^       | |  || 
  | |   ^ / \      | |  ||
  | |  / |   |     | |  ||
  | | |  |   |     | |  ||
  | | |   \ /  /\  | |  ||
  | |  \ / v   ||  | |  //
  | |   v      \/  | | //
  | |______________| |//
  \__________________//`

func width(text string) int {
	var width = 0
	for _, line := range strings.Split(text, "\n") {
		if len(line) > width {
			width = len(line)
		}
	}
	return width
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

func selectFont() string {
	var i = rnd.Intn(len(fonts))
	return fonts[i]
}

var fonts = []string{
	"3x5",
	"5lineoblique",
	"acrobatic",
	"avatar",
	"banner",
	"banner3-D",
	"banner3",
	"banner4",
	"barbwire",
	"basic",
	"bell",
	"big",
	"binary",
	"block",
	"bubble",
	"bulbhead",
	"calgphy2",
	"caligraphy",
	"chunky",
	"coinstak",
	"colossal",
	"computer",
	"contessa",
	"contrast",
	"cosmic",
	"cosmike",
	"cricket",
	"cursive",
	"cyberlarge",
	"cybermedium",
	"cybersmall",
	"diamond",
	"digital",
	"doh",
	"doom",
	"dotmatrix",
	"drpepper",
	"eftifont",
	"eftipiti",
	"eftirobot",
	"eftitalic",
	"eftiwall",
	"eftiwater",
	"epic",
	"fender",
	"fourtops",
	"fuzzy",
	"gothic",
	"graffiti",
	"hollywood",
	"isometric1",
	"isometric2",
	"isometric3",
	"isometric4",
	"italic",
	"jazmine",
	"kban",
	"larry3d",
	"lcd",
	"letters",
	"lockergnome",
	"madrid",
	"maxfour",
	"mike",
	"mini",
	"mnemonic",
	"morse",
	"moscow",
	"nancyj-fancy",
	"nancyj-underlined",
	"nancyj",
	"nipples",
	"ntgreek",
	"o8",
	"ogre",
	"pawp",
	"peaks",
	"pebbles",
	"pepper",
	"pyramid",
	"rectangles",
	"relief2",
	"rev",
	"roman",
	"rot13",
	"rounded",
	"rowancap",
	"rozzo",
	"runic",
	"runyc",
	"sblood",
	"script",
	"serifcap",
	"shadow",
	"short",
	"slant",
	"slide",
	"slscript",
	"small",
	"smscript",
	"smshadow",
	"smslant",
	"smtengwar",
	"speed",
	"stampatello",
	"standard",
	"starwars",
	"stellar",
	"stop",
	"straight",
	"tanja",
	"term",
	"thick",
	"thin",
	"threepoint",
	"ticks",
	"tinker-toy",
	"tombstone",
	"trek",
	"tsalagi",
	"twopoint",
	"univers",
	"usaflag",
	"wavy",
	"weird",
}
