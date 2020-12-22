// Implementation of Position weight matrix on a DnaA 9-mers.

package main

import (
	"fmt"
	"os"
	"bufio"
	"strconv"
	"math"
	"strings"
	"io/ioutil" 
)


// declare functions
func readLines(path string) ([]string, error){
	// Returns array of the string and error
	file, err := os.Open(path)
	
	if err != nil {
		panic(err)
	}
	
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		}


	return lines[4:], scanner.Err()

}

func readRaw(path string) string{
	bs, err := ioutil.ReadFile(path)
	
	if err != nil {
		panic(err)
	}
	str := string(bs)

	return str

}

func writeRaw(lines string, path string)error{
	// Reads from string and write to .txt
	file, err := os.Create(path)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(lines)

	return err
}


func writeLines(lines []string, path string) error {
	// Reads from array and write to .txt
	file, err := os.Create(path)

	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	
	for _, line := range lines{
		
		var new_line []string = strings.Split(line[67:76], ":") // & extract 9 characters from each line from position 67 to 75

		fmt.Fprintln(w, new_line[(len(new_line)) - 1])
	}		
		return w.Flush()
}


func Round(x, unit float64) float64 {
	// math extension
	return math.Round(x * math.Pow(10, unit)) / math.Pow(10, unit)
}

// declare variables

var (
	// name string
	test_path string = "data/temp.txt"
	file_path string = "data/EI_true.seq"
	ei_nine_path string = "outputs/EI_nine.txt"
	score_file string = "outputs/EI_nine_pwm.txt" 
	output_file string = "outputs/EI_nine_output.txt"
)


func main() {
	
	splitted_content, _ := readLines(file_path)
	
	writeLines(splitted_content, ei_nine_path)
	

	// New txt file for position matrix EI_nine
	// initialize the PWM with four rows and nine columns

	var a, c, g, t [9]float64
	input_file, _ := readLines(ei_nine_path)

	var num_lines float64 = float64(len(input_file))
	
	
	// read line by line and update the PWM with the frequencies of each base at
	// the 9 position
	for _, line := range input_file{
		
		for i:=0; i < 9; i++ {

			switch string(line[i]) {
			case "A": a[i] += 1.0
			case "C": c[i] += 1.0
			case "G": g[i] += 1.0
			case "T": t[i] += 1.0
			}
		}
	}


	for i:=0; i < 9; i++ {
		
		a[i] = Round(math.Log2((a[i] + 0.1) / (num_lines + 0.4) / 0.25), 3)
		c[i] = Round(math.Log2((c[i] + 0.1) / (num_lines + 0.4) / 0.25), 3)
		g[i] = Round(math.Log2((g[i] + 0.1) / (num_lines + 0.4) / 0.25), 3)
		t[i] = Round(math.Log2((t[i] + 0.1) / (num_lines + 0.4) / 0.25), 3)
		
	}

	// write the position weight matrix
	genes := map[string][9]float64{
		"A":a,
		"C":c,
		"G": g,
		"T": t,
	}	

	var text string

	for key, gene := range genes {
		text += key + "\t" + "\t"
		
		for i:=0; i < 9; i++ {
			text += strconv.FormatFloat((gene[i]), 'f', 3, 64) + "\t" 
		}
		text += "\n"
	}
	
	writeRaw(text, score_file)

	// Final, write to output
	var score float64
	var temp_str string

	for _, line := range input_file{
		score = 0.0

		for i:=0; i < 9; i++ {

			switch string(line[i]) {
				case "A": score += a[i]
				case "C": score += c[i]
				case "G": score += g[i]
				case "T": score += t[i]		
			}
		
			temp_str += line + "\t" + strconv.FormatFloat(score, 'f', 3, 64) + "\n" 
		
		}

		writeRaw(temp_str, output_file)
		
	}


	fmt.Println()
	fmt.Println("Done...")
	// n := readRaw(output_file)
	// fmt.Println(n)

}

