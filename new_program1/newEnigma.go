/*
References:
1. https://github.com/emedvedev/enigma
I used his enigma emulator. Understood his code and used this program to encode text. My work is to use
algorithm to figure out the configuration of the enigma machine as fast as possible.

2. https://github.com/becgabri/enigma
I actually used this program because  the enigma simulator (github.com/emedvedev/enigma)
has a bug in its code.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	//"text/template"

	//"text/template"

	//"github.com/becgabri/enigma"
	"../../cli"
	"../../enigma"
)

// CLIOpts sets the parameter format for Enigma CLI. It also includes a "help"
// flag and a "condensed" flag telling the program to output plain result.
// Also, this CLI module abuses tags so much it hurts. Oh well. ¯\_(ツ)_/¯
/*
type CLIOpts struct {
	Help      bool `cli:"!h,help" usage:"Show help."`
	Condensed bool `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`

	Rotors    string `cli:"rotors" name:"\"I II III\"" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
	Rings     string `cli:"rings" name:"\"1 1 1\"" usage:"Rotor rings offset: from 1 (default) to 26 for each rotor."`
	Position  string `cli:"position" name:"\"A A A\"" usage:"Starting position of the rotors: from A (default) to Z for each."`
	Plugboard string `cli:"plugboard" name:"\"AB CD\"" usage:"Optional plugboard pairs to scramble the message further."`

	Reflector string `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-thin, C-thin."`
}
*/
// CLIDefaults is used to populate default values in case
// one or more of the parameters aren't set. It is assumed
// that rotor rings and positions will be the same for all
// rotors if not set explicitly, so only one value is stored.
var CLIDefaults = struct {
	Reflector string
	Ring      string
	Position  string
	Rotors    string
	//Plugboard string
}{
	Reflector: "C-thin",
	Ring:      "1 1 1 16",
	Position:  "D A B Q",
	Rotors:    "Gamma VI IV III",
	//Plugboard: ""
}

// SetDefaults sets values for all Enigma parameters that
// were not set explicitly.
// Plugboard is the only parameter that does not require a
// default, since it may not be set, and in some Enigma versions
// there was no plugboard at all.

func SetDefaults(argv *CLIOpts) {
	if argv.Reflector == "" {
		argv.Reflector = CLIDefaults.Reflector
	}
	if len(argv.Rotors) == 0 {
		argv.Rotors = CLIDefaults.Rotors
	}
	loadRings := (len(argv.Rings) == 0)
	loadPosition := (len(argv.Position) == 0)
	if loadRings {
		argv.Rings = CLIDefaults.Ring
	}
	if loadPosition {
		argv.Position = CLIDefaults.Position

	}
	//mark
	// if argv.Plugboard == "" {
	// 	argv.Plugboard = CLIDefaults.Plugboard
	// }
}

func main() {

	var bestRotors = []string{"Gamma", "VI", "IV", "III"}
	var bestPositions = []string{"D", "A", "B", "Q"}
	var bestPlugboard = []string{"MS", "KU", "FY", "AG", "BN", "PQ", "HJ", "DI", "ER", "LW"}
	var ring_array = []int{1, 1, 1, 16}

	cli.SetUsageStyle(cli.DenseManualStyle)
	cli.Run(new(CLIOpts), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*CLIOpts)

		//read the text to process
		textPath := strings.Join(ctx.Args(), " ")
		fi, pathErr := os.Open(textPath)
		if pathErr != nil {
			panic(pathErr)
		}
		defer fi.Close()
		fd, fileErr := ioutil.ReadAll(fi)
		if fileErr != nil {
			panic(fileErr)
		}
		originalPlaintext := string(fd)

		plaintext := enigma.SanitizePlaintext(originalPlaintext)
		fmt.Print("original text:  ")
		fmt.Println(plaintext)

		//set rotor, ring, start position as config
		config := make([]enigma.RotorConfig, len(bestRotors))
		for index, rotor := range bestRotors {
			ring := ring_array[index]
			value := bestPositions[index][0]
			config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
		}

		argv.Reflector = "C-thin"
		e := enigma.NewEnigma(config, argv.Reflector, bestPlugboard)

		//after setting, encode the text
		encoded := e.EncodeString(plaintext)
		fmt.Println("\nfinal result:  ")
		fmt.Print("encoded text:  ")
		fmt.Println(encoded)

		fmt.Print("Rotors:  ")
		fmt.Println(bestRotors)
		fmt.Print("Positions:  ")
		fmt.Println(bestPositions)
		fmt.Print("Plugborad:   ")
		fmt.Println(bestPlugboard)

		return nil
	})

}
