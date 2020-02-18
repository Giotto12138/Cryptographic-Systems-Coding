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
	"strconv"
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
type CLIOpts struct {
	Help      bool `cli:"!h,help" usage:"Show help."`
	Condensed bool `cli:"c,condensed" name:"false" usage:"Output the result without additional information."`

	Rotors    string `cli:"rotors" name:"\"I II III\"" usage:"Rotor configuration. Supported: I, II, III, IV, V, VI, VII, VIII, Beta, Gamma."`
	Rings     string `cli:"rings" name:"\"1 1 1\"" usage:"Rotor rings offset: from 1 (default) to 26 for each rotor."`
	Position  string `cli:"position" name:"\"A A A\"" usage:"Starting position of the rotors: from A (default) to Z for each."`
	Plugboard string `cli:"plugboard" name:"\"AB CD\"" usage:"Optional plugboard pairs to scramble the message further."`

	Reflector string `cli:"reflector" name:"C" usage:"Reflector. Supported: A, B, C, B-thin, C-thin."`
}

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
	Position:  "A A B Q",
	Rotors:    "Beta II IV III",
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

/*
calculate the Ioc of a text
*/

func IocSum(text string) float64 {
	var l int
	var count int
	var Ioc float64
	var letterTimes [26]int

	l = len(text)

	for i := 'A'; i <= 'Z'; i++ {
		letter := string(i)
		num := strings.Count(text, letter)
		letterTimes[count] = num
		count++
	}

	for _, value := range letterTimes {
		value := float64(value)
		l := float64(l)
		Ioc = Ioc + (value*(value-1))/(l*(l-1))
	}
	Ioc = Ioc * 26

	return Ioc
}

func RemoveRepByMap(slc []string) []string {
	result := []string{}
	tempMap := map[string]byte{} // store different keys
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // after adding to the map, if the length of map changes, no same element
			result = append(result, e)
		}
	}
	return result
}

func findPlugs(maxIoc float64, bestRotors [4]string, ring_array []int, bestPositions [4]string, lastPlugsBoard []string, plaintext string) []string {
	bestPlugboardNew := lastPlugsBoard
	plugsStoreNew := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	var bestPairNew string
	var newPair string

	for firstLetter := 0; firstLetter < len(plugsStoreNew); firstLetter++ { //the first letter of the pair
		tempMaxIoc := maxIoc

		for secondLetter := 0; secondLetter < len(plugsStoreNew); secondLetter++ { //the second letter of the pair

			if firstLetter == secondLetter { // we don't switch the same letter
				continue
			}

			//set rotor, ring, start position as config
			config := make([]enigma.RotorConfig, len(bestRotors))
			for index, rotor := range bestRotors {
				ring := ring_array[index]
				value := bestPositions[index][0]
				config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
			}

			newPair = plugsStoreNew[firstLetter] + plugsStoreNew[secondLetter] //new pair of two letter

			e := enigma.NewEnigma(config, "C-thin", append(bestPlugboardNew, newPair))

			//after setting, encode the text
			encoded := e.EncodeString(plaintext)

			newIoc := IocSum(encoded) //get the current Ioc of the encoded text
			if newIoc > tempMaxIoc {
				tempMaxIoc = newIoc // store the max Ioc
				// store the setting of the max Ioc
				//argv.Plugboard = strings.Join(bestPlugboardNew, " ")
				bestPairNew = newPair
			}
			//fmt.Print("MaxIoc:  ")
			//fmt.Println(tempMaxIoc)
			fmt.Print("bestPlugboardNew:  ")
			fmt.Println(bestPlugboardNew)

		} //the second letter of the pair
		//resolve conflicts
		if len(plugsStoreNew) > 6 {
			for p := 0; p < len(plugsStoreNew); p++ {
				if strings.Contains(bestPairNew, plugsStoreNew[p]) {
					plugsStoreNew = append(plugsStoreNew[:p], plugsStoreNew[p+1:]...)
					p--
				}
			}
		}

		//add a new good pair
		if len(bestPlugboardNew) < 10 {
			bestPlugboardNew = append(bestPlugboardNew, bestPairNew)
		}
		bestPlugboardNew = RemoveRepByMap(bestPlugboardNew)

	} //the first letter of the pair
	return bestPlugboardNew
}

func main() {

	var maxIoc float64
	var bestRotors [4]string
	var bestPositions [4]string

	//for plugboard
	var bestPlugboard []string
	var newPair string
	var bestPair string

	rotorsStore := [6]string{"Beta", "Gamma", "I", "II", "V", "VI"}
	positionsStore := [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	plugsStore := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	//var count int

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

		//originalPlaintext := strings.Join(ctx.Args(), " ")
		originalPlaintext := string(fd)

		plaintext := enigma.SanitizePlaintext(originalPlaintext)
		//fmt.Println(plaintext)

		if argv.Help || len(plaintext) == 0 {
			com := ctx.Command()
			com.Text = DescriptionTemplate
			ctx.String(com.Usage(ctx))
			return nil
		}

		//settings of the enigma
		rotor_array := strings.Split(argv.Rotors, " ")
		var ring_array []int = make([]int, len(strings.Split(argv.Rings, " ")))
		pos_array := strings.Split(argv.Position, " ")
		for idx, val := range strings.Split(argv.Rings, " ") {
			ring_array[idx], _ = strconv.Atoi(val)
		}

		fmt.Println("rotors:       ")
		fmt.Println(rotor_array)
		//plugboards := strings.Split(argv.Plugboard, " ")

		//find the optimal setting
		//first round, determine plugboard roughly, find the setting of rotors and positions
		for firstRotor := 0; firstRotor < 2; firstRotor++ { //try "Beta, Gamma" for first rotor
			rotor_array[0] = rotorsStore[firstRotor]

			for secondRotor := 2; secondRotor < 6; secondRotor++ { //try "I", "II", "V", "VI" for second rotor
				rotor_array[1] = rotorsStore[secondRotor]

				for firstPosition := 0; firstPosition < 26; firstPosition++ { // try 26 letters for the position of the first rotor
					pos_array[0] = positionsStore[firstPosition]

					for secondPosition := 0; secondPosition < 26; secondPosition++ { // try 26 letters for the position of the second rotor
						pos_array[1] = positionsStore[secondPosition]

						for firstLetter := 0; firstLetter < len(plugsStore); firstLetter++ { //the first letter of the pair
							for secondLetter := 0; secondLetter < len(plugsStore); secondLetter++ { //the second letter of the pair

								if firstLetter == secondLetter { // we don't switch the same letter
									continue
								}

								//record the the number of attempts
								// count++
								// fmt.Print("\nThe ")
								// fmt.Print(count)
								// fmt.Println(" times attempt")

								//set rotor, ring, start position as config
								config := make([]enigma.RotorConfig, len(rotor_array))
								for index, rotor := range rotor_array {
									ring := ring_array[index]
									value := pos_array[index][0]
									config[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
								}

								//hillclimb of the plugboard

								newPair = plugsStore[firstLetter] + plugsStore[secondLetter] //new pair of two letter
								plugboards := []string{newPair}

								// fmt.Print("tempPlugboards:  ")
								// fmt.Println(tempPlugboards)

								e := enigma.NewEnigma(config, argv.Reflector, plugboards)

								//after setting, encode the text
								encoded := e.EncodeString(plaintext)

								newIoc := IocSum(encoded) //get the current Ioc of the encoded text
								if newIoc > maxIoc {
									maxIoc = newIoc // store the max Ioc
									// store the setting of the max Ioc
									argv.Rotors = strings.Join(rotor_array, " ")
									copy(bestRotors[:], rotor_array)
									argv.Position = strings.Join(pos_array, " ")
									copy(bestPositions[:], pos_array)
									argv.Plugboard = strings.Join(bestPlugboard, " ")
									bestPair = newPair
									//copy(bestPlugboard[:], plugboards)
									//bestPlugboard = tempPlugboards
								}
								//fmt.Print("MaxIoc:  ")
								//fmt.Println(maxIoc)
								fmt.Print("bestPlugboard:  ")
								fmt.Println(bestPlugboard)

								if argv.Condensed {
									fmt.Print(encoded)
									return nil
								}
								/*
									tmpl, _ := template.New("cli").Parse(OutputTemplate)
									err := tmpl.Execute(os.Stdout, struct {
										Original, Plain, Encoded string
										Args                     *CLIOpts
										Ctx                      *cli.Context
									}{originalPlaintext, plaintext, encoded, argv, ctx})
									if err != nil {
										return err
									}
								*/

							} //the second letter of the pair
							//resolve conflicts
							if len(plugsStore) > 6 {
								for p := 0; p < len(plugsStore); p++ {
									if strings.Contains(bestPair, plugsStore[p]) {
										plugsStore = append(plugsStore[:p], plugsStore[p+1:]...)
										p--
									}
								}
							}

							//add a new good pair
							if len(bestPlugboard) < 10 {
								bestPlugboard = append(bestPlugboard, bestPair)
							}
							bestPlugboard = RemoveRepByMap(bestPlugboard)

						} //the first letter of the pair

					} //second position of the rotor
				} //first position of the rotor
			} //secondRotor
		} //firstRotor

		//second round, use the setting of rotors and positions, determine plugboard again
		//find the optimal setting
		var lastPlugsBoard []string
		var bestPlugboardNew []string
		var plugsCount int
		for plugsCount <= 10 {
			bestPlugboardNew = findPlugs(maxIoc, bestRotors, ring_array, bestPositions, lastPlugsBoard, plaintext)
			lastPlugsBoard = bestPlugboardNew
			plugsCount = len(lastPlugsBoard)
		}

		//print optimal encoded text
		fmt.Println("final result:  ")
		optimalConfig := make([]enigma.RotorConfig, len(bestRotors))
		for index, rotor := range bestRotors {
			ring := ring_array[index]
			value := bestPositions[index][0]
			optimalConfig[index] = enigma.RotorConfig{ID: rotor, Start: value, Ring: ring}
		}
		optimalE := enigma.NewEnigma(optimalConfig, argv.Reflector, bestPlugboardNew)
		//after setting, encode the text
		optimalEncoded := optimalE.EncodeString(plaintext)
		fmt.Println(optimalEncoded)

		return nil
	})

}
