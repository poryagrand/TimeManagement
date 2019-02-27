// ========================================
// ======= Developer : Porya Elahi ========
// ======= Web Site  : PoryaGrand.ir ======
// ======= Github    : PoryaGrand =========
// ========================================
// ============= 2019/02/27 ===============
// ========================================

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func main() {
	//initial console screen
	SelectedURL := Init()

	// create input ready object
	var reader = bufio.NewReader(os.Stdin)

	// ready file data to parse as json
	data, err := ioutil.ReadFile(SelectedURL)
	check(err)

	// parse file as object into <Response>
	var Response []TimeBox
	json.Unmarshal(data, &Response)

	var command = ""
	// create chanel to make timer and listen to keyboard
	ch := make(chan int)
	for {
		fmt.Printf("Work Time Duration Application is Started up [%s]\n------------------------\n", SelectedURL)
		fmt.Printf("\nChoose One Of Commands:")
		fmt.Printf("\n    1-New Duration")
		fmt.Printf("\n    2-Statistics")
		fmt.Printf("\n    3-Save")

		fmt.Printf("\nInput: ")
		command, _ = reader.ReadString('\n')
		clc()

		// if the comand be showing statistics
		if strings.TrimSpace(command) == "2" {
			rlen := len(Response)
			var tm time.Time
			var allsum int64
			for n := 0; n < rlen; n++ {
				tm = time.Unix(Response[n].Date, 0)
				dlen := len(Response[n].Durations)
				var sum int64
				for k := 0; k < dlen; k++ {
					sum = sum + Response[n].Durations[k]
				}
				allsum = allsum + sum
				fmt.Printf("\nState %d On %s : %s ", n+1, tm.Format("2006/01/02 15:04:05"), fmtDurationSec(sum))
			}
			fmt.Printf("\n________________________________________________\n\n         Sum Of All Work Hours : %s", fmtDurationSec(allsum))
			fmt.Printf("\n\nTo return Main Menu , Press Enter...  ")
			_, _ = reader.ReadString('\n')
			clc()
		} else if strings.TrimSpace(command) == "3" { // if the command be save the file

			obj, _ := json.Marshal(&Response)

			d1 := []byte(string(obj))
			err := ioutil.WriteFile(SelectedURL, d1, 0644)
			check(err)
			clc()
			fmt.Println("\t\t\t░░░░ Data Has Been Saved ░░░░\n")
		} else if strings.TrimSpace(command) == "1" { // if the command be start new duration
			box := TimeBox{Date: time.Now().UnixNano() / int64(time.Second)}

			step := 1
			// loop on all duration
			for {
				// check input and listen to [enter] on each second
				go func() {
					command, _ = reader.ReadString('\n')
					ch <- 1
				}()
				// start time of a duration step
				start := time.Now().Round(0)
			FL:
				for {
					// start a step
					clc()
					fmt.Printf("Step %d Started [Note: Press Enter To Pause]\n", step)
					dur := time.Since(start)
					fmt.Printf("\n    " + fmtDuration(dur))
					select {
					case <-ch:
						break FL
					case <-time.After(time.Second):
					}

				}
				// calculate duration in seconds(int)
				box.Durations = append(box.Durations, int64(time.Since(start)/(time.Second)))

				fmt.Printf("\nStep %d Ended. \ndo you want to create new step [y/n] ?  ", step)

				command, _ = reader.ReadString('\n')
				clc()

				// if the answeer be yes then
				// new duration will be created
				// otherwise returns to main menu
				if strings.TrimSpace(command) == "y" {
					step++
				} else {
					break
				}
			}

			Response = append(Response, box)
		}
	}

}

// main struct of storage json file
type TimeBox struct {
	Date      int64   `json:"date"`
	Durations []int64 `json:"durations"`
}

// initialize start screen to show developer details
func Init() string {
	// clear screan
	clc()

	str1 := "\n\n\n\t\t\t\tTime Management System"
	str2 := "\n\t\t\t\t     PoryaGrand.ir"
	str3 := "\t\t\t\tgithub.com/poryagrand"
	str4 := "\n\t\t\t\t   Please Wait..."

	fmt.Printf("%s \n", str1)
	fmt.Printf("%s \n", str2)
	fmt.Printf("%s \n", str3)
	fmt.Printf("%s \n", str4)

	time.Sleep(3 * time.Second)

	// return a storage file url
	return GetStorage()
}

// generate new or choose from existance file json storage
func GetStorage() string {

	clc()

	var reader = bufio.NewReader(os.Stdin)

	fmt.Printf("Choose Storage Option:")
	fmt.Printf("\n    1-Create Storage")
	fmt.Printf("\n    2-Open Storage")
	fmt.Printf("\nInput: ")
	command, _ := reader.ReadString('\n')
	command = strings.TrimSpace(command)
	clc()
	// if wants to create new storage
	if command == "1" {
		fmt.Printf("\nEnter Storage name and then press [Enter] Button: ")
		command, _ = reader.ReadString('\n')
		command = strings.TrimSpace(command)
		clc()

		if fileExists("storage/" + string(command) + ".json") {
			return "storage/" + string(command) + ".json"
		}

		content := []byte("[]")
		err := ioutil.WriteFile("storage/"+string(command)+".json", content, 0644)

		// after creating storage, it will go to main menu for created storage
		if err == nil {
			return "storage/" + string(command) + ".json"
		}
		return "storage/default.json"
	} else {
		root := "storage"
		var files []string
		// find all json files in storage directory
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			fname := strings.Split(info.Name(), ".")

			if info.Name() != root || (len(fname) >= 2 && fname[len(fname)-1] == "json") {
				files = append(files, strings.Split(info.Name(), ".")[0])
			}
			return nil
		})
		if err == nil {
			fmt.Printf("\nSelect the desired storage file: ")

			flen := len(files)
			for n := 0; n < flen; n++ {
				fmt.Printf("\n    %d- %s", n+1, files[n])
			}

			fmt.Printf("\nInput: ")
			command, _ = reader.ReadString('\n')
			command = strings.TrimSpace(command)
			index, err := strconv.Atoi(command)

			clc()
			if err == nil && (index-1) < flen {
				return root + "/" + files[index-1] + ".json"
			}
			return "storage/default.json"

		} else {
			clc()
			return "storage/default.json"
		}
	}

}

// return time duration in string format
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%d:%d:%d", h, m, s)
}

// return int duration in string format
func fmtDurationSec(s int64) string {
	var m int64
	var h int64
	for s >= 60 {
		s = s - 60
		m++
	}
	for m >= 60 {
		m = m - 60
		h++
	}
	return fmt.Sprintf("%d:%d:%d", h, m, s)
}

// clear screen
func clc() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// check error and panic on need
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// chech file existance
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
