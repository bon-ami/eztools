package eztools

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

const formatComp, formatMinus = "20060102", "2006-01-02"

var (
	// Debugging marks debugging mode
	Debugging bool // whether more debug procedures
	// Verbose marse debugging output level
	Verbose  = 0
	defaults bool // whether no more confirmations asked
	logger   *log.Logger
	logFile  *os.File
)

// InitLogger opens log file
func InitLogger(out *os.File /*io.Writer*/) error {
	if out == nil {
		return ErrInvalidInput
	}
	logger = log.New(out, "", log.LstdFlags)
	if logger == nil {
		return errors.New("no logger created")
	}
	logFile = out
	return nil
}

// logOrPrint logs it
// Parameter print2 = also shown on screen
func logOrPrint(print2 bool, out ...interface{}) {
	if logger != nil {
		logger.Println(out...)
		if print2 {
			fmt.Println(out...)
		}
		_ = logFile.Sync()
	} else {
		log.Println(out...)
	}
}

// Log logs it
func Log(out ...interface{}) {
	logOrPrint(false, out...)
}

// LogPrint logs and prints it
func LogPrint(out ...interface{}) {
	logOrPrint(true, out...)
}

// LogPrintWtTime logs and prints a string with time
func LogPrintWtTime(out string) {
	LogPrint(time.Now().String() + ": " + out)
}

// LogErrPrint logs and prints error
func LogErrPrint(err error) {
	if err != nil {
		LogPrint(err.Error())
	}
}

// LogErrPrintWtInfo logs and prints error with information string
func LogErrPrintWtInfo(info string, err error) {
	if err != nil {
		LogPrint(info + ": " + err.Error())
	}
}

// LogErr logs error
func LogErr(err error) {
	if err != nil {
		Log(err.Error())
	}
}

// LogErrWtInfo logs error with information string
func LogErrWtInfo(info string, err error) {
	if err != nil {
		Log(info + ": " + err.Error())
	}
}

// LogFatal logs and prints it and exits
func LogFatal(out ...interface{}) {
	if logger != nil {
		fmt.Println(out...)
		logger.Fatalln(out...)
	} else {
		log.Fatalln(out...)
	}
}

// LogErrFatal logs and prints error and exits
func LogErrFatal(err error) {
	LogFatal(err.Error())
}

// ShowStrln prints it with a line break
func ShowStrln(ps string) {
	fmt.Println(ps)
}

// ShowStr prints it with no line breaks
func ShowStr(ps string) {
	fmt.Print(ps)
}

// ShowArrln prints a slice in one line
func ShowArrln(arr []string) {
	for _, i := range arr {
		fmt.Print("\"")
		fmt.Print(i)
		fmt.Print("\", ")
	}
	fmt.Print("\n")
}

// ShowByteln prints byte slice
func ShowByteln(ps []byte) {
	os.Stdout.Write(ps)
	fmt.Print("\n")
}

// ShowSthln prints all kinds of stuffs with line feed
func ShowSthln(sth ...interface{}) {
	fmt.Printf("%v\n", sth...)
}

// PromptStr prompts user and get input
func PromptStr(ps string) string {
	fmt.Print(ps + ":")
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	return in.Text()
}

// PromptPwd prompts user and get password
func PromptPwd(ps string) string {
	fmt.Print(ps + ":")
	pwd, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		return ""
	}
	return string(pwd)
}

// PromptIntStr prompts user and gets two inputs
// Return values. zero values are default
func PromptIntStr(pi string, ps string) (i int, s string) {
	fmt.Printf("%s %s:", pi, ps)
	fmt.Scanf("%d %s", &i, &s)
	//to exhaust the buffer
	bufio.NewScanner(os.Stdin).Scan()
	return
}

// PromptInt prompts user and gets input
// Return values. zero values are default
func PromptInt(pi string) (res int, err error) {
	res, err = strconv.Atoi(PromptStr(pi))
	return
}

// ChkCfmNPrompt checks defaults and return false only when user replied exception
// program exits when user replied 'q' or 'e'
// no more confirmations when user replied 'a' or 'c'
// verbose set when user replied a number, in which case the prompt will show again
func ChkCfmNPrompt(noti, exception string) bool {
	if defaults {
		return true
	}
	quitCode := "q"
	confCode := "a"
	switch exception {
	case "q":
		quitCode = "e"
	case "a":
		confCode = "e"
	}
	val := PromptStr(noti + "?(any number=reset verbose level and ask again/" +
		quitCode + "=quit program/" + confCode +
		"=defaults to all confirmations/" + exception + "/...)")
	switch val {
	case quitCode:
		LogFatal("Quiting")
	case confCode:
		defaults = true
	case exception:
		return false
	default:
		if v, err := strconv.Atoi(val); err == nil {
			Verbose = v
			return ChkCfmNPrompt(noti, exception)
		}
	}
	return true
}

// ChooseStringsWtIDs is for general usage to
// ask user to choose from a slice
// parameters.
//	fL=quantity of elements
//	fI=get index to match user's input
//	fV=get message to show for each index
//	notif=notification string for user
func ChooseStringsWtIDs(fL func() int, fI func(int) int,
	fV func(int) string, notif string) (res int) {
	len := fL()
	if len < 1 {
		return InvalidID
	}
	for i := 0; i < len; i++ {
		fmt.Printf("%d: %s\n", fI(i), fV(i))
	}
	res, err := PromptInt(notif)
	if err == nil {
		//check for invalid input
		for i := 0; i < len; i++ {
			if fI(i) == res {
				return
			}
		}
	}
	return InvalidID
}

// ChooseInts asks user to choose from a slice
// Parameters. arr[][0]=id. arr[][1]=string
func ChooseInts(arr [][]string, notif string) (id int) {
	return ChooseStringsWtIDs(
		func() int {
			return len(arr)
		},
		func(i int) int {
			ret, err := strconv.Atoi(arr[i][0])
			if err != nil {
				return InvalidID
			}
			return ret
		},
		func(i int) string {
			return arr[i][1]
		},
		notif)
}

// Return values. zero value is default
func choosePairs(choices []pairs, notif string) (res int) {
	return ChooseStringsWtIDs(
		func() int {
			return len(choices)
		},
		func(i int) int {
			return choices[i].id
		},
		func(i int) string {
			return choices[i].str
		},
		notif)
}

// ChooseStrings asks user to choose from a slice
func ChooseStrings(choices []string) (res int) {
	if len(choices) < 1 {
		return InvalidID
	}
	for i, v := range choices {
		fmt.Printf("%d: %s\n", i, v)
	}
	fmt.Print("Your choice is: ")
	var str string
	fmt.Scanln(&str)
	if len(str) < 1 {
		return InvalidID
	}
	res, err := strconv.Atoi(str)
	if err != nil {
		return InvalidID
	}
	//check for invalid input
	for i := range choices {
		if i == res {
			return
		}
	}
	return InvalidID
}

// GetDate asks user to input a date string
func GetDate(info string) string {
	fmt.Print(info + "date such as " + formatComp + ": ")
	var res string
	fmt.Scanln(&res)
	t, err := time.Parse(formatComp, res)
	if err == nil {
		return t.Format(formatComp)
	}
	return "NULL"
}

// TranDate removes minuses from date string
func TranDate(date string) string {
	t, err := time.Parse(formatMinus, date)
	if err == nil {
		return t.Format(formatComp)
	}
	return ""
}
