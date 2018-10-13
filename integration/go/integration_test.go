package integration_test

import (
	"bytes"
	"fmt"
	"github.com/onsi/gomega/types"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"

	. "github.com/onsi/gomega"
)

var errFile *os.File
var outFile *os.File

type entry struct {
	exampleDir    string
	directory     string
	context       string
	it            string
	fileRunOutput string
	fileRunErr    bool
	possibleLogs  []string
	outLog        []byte
	errFile       string
	errLog        []byte
}

var table = []entry{
	{
		exampleDir:   "examples-go",
		directory:    "waitevent",
		context:      "with a waitevent workflow",
		it:           "should wait for either an event or a given time",
		possibleLogs: []string{"Task A starts\nTask A ends\n"},
		errFile:      "",
	},
	{
		exampleDir:   "examples-go",
		directory:    "event",
		context:      "with an event based workflow",
		it:           "should handle events",
		possibleLogs: []string{"Task A starts\nTask A ends\nTask B starts\nTask B ends\n"},
		errFile:      "",
	},
	{
		exampleDir:   "examples-go",
		directory:    "wait",
		context:      "with a wait workflow",
		it:           "should wait before a task",
		possibleLogs: []string{"Task A starts\nTask A ends\nTask B starts\nTask B ends\n"},
		errFile:      "",
	},
	{
		exampleDir: "examples-go",
		directory:  "asynchronous",
		context:    "with an asychronous workflow",
		it:         "should perform tasks asynchronously",
		possibleLogs: []string{"Task A starts\nTask B starts\nTask C starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
			"Task A starts\nTask C starts\nTask B starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
			"Task B starts\nTask A starts\nTask C starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
			"Task B starts\nTask C starts\nTask A starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
			"Task C starts\nTask A starts\nTask B starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
			"Task C starts\nTask B starts\nTask A starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n"},
		errFile: "",
	},
	{
		exampleDir: "examples-go",
		directory:  "parallel",
		context:    "with a parallel workflow",
		it:         "should handle tasks in parallel",
		possibleLogs: []string{"Task A starts\nTask B starts\nTask A ends\nTask B ends\nTask D starts\nTask D ends\n",
			"Task B starts\nTask A starts\nTask A ends\nTask B ends\nTask D starts\nTask D ends\n"},
		errFile: "",
	},
	{
		exampleDir:   "examples-go",
		directory:    "recursive",
		context:      "with a recursive workflow",
		it:           "should handle tasks recursively",
		possibleLogs: []string{"01Iteration:101Iteration:201"},
		errFile:      "",
	},
	{
		exampleDir:   "examples-go",
		directory:    "sequential",
		context:      "with a sequential workflow",
		it:           "should handle tasks sequentially",
		possibleLogs: []string{"Task A starts\nTask A ends\nTask C starts\nTask C ends\nTask D starts\nTask D ends\n"},
		errFile:      "",
	},
	//{ // this test is too fickle for now
	//	exampleDir: "examples-go",
	//	directory:  "version",
	//	context:    "with a version workflow",
	//	it:         "should handle versioned workflows",
	//	possibleLogs:    []string{"Task A starts\nTask B starts\nTask C starts\nTask D starts\nTask A ends\nTask B ends\nTask C ends\nTask D ends\n",
	//						"Task B starts\nTask A starts\nTask C starts\nTask D starts\nTask A ends\nTask B ends\nTask C ends\nTask D ends\n"},
	//	errFile:    "",
	//
	//},
}

func TestSetup(t *testing.T) {

	g := NewGomegaWithT(t)


	t.Log("It Should Copy Examples folder locally")
	{
		err := os.RemoveAll("examples-go")
		if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
			g.Expect(err).NotTo(HaveOccurred())
		}

		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = build.Default.GOPATH
		}

		output, err := exec.Command("cp", "-r", gopath+"/src/github.com/zenaton/examples-go", "examples-go").CombinedOutput()
		if err != nil {
			fmt.Println("error copying examples-go, output: ", string(output), "err: ", err)
			g.Expect(err).NotTo(HaveOccurred())
		}

		output, err = exec.Command("rm", "-rf", "examples-go/vendor", "examples-go/Gopkg.lock", "examples-go/Gopkg.toml").CombinedOutput()

		if err != nil {
			fmt.Println("error removing vendoring, output: ", string(output), "err: ", err)
			g.Expect(err).NotTo(HaveOccurred())
		}
	}

	t.Log("It should inject log prefixes into tasks")
	{
		files, err := ioutil.ReadDir("examples-go/workflows")
		g.Expect(err).NotTo(HaveOccurred())

		for _, f := range files {
			contents, err := ioutil.ReadFile("examples-go/workflows/" + f.Name())
			g.Expect(err).NotTo(HaveOccurred())

			prefix := strings.TrimRight(f.Name(), ".go")
			newContents := bytes.Replace(contents, []byte(".New()"), []byte(".New(`::"+prefix+"`)"), -1)
			newContents = bytes.Replace(newContents, []byte("github.com/zenaton/examples-go/tasks"), []byte("github.com/zenaton/zenaton-go/integration/go/examples-go-new-tasks/tasks"), -1)

			err = ioutil.WriteFile("examples-go/workflows/"+f.Name(), newContents, 0660)
			g.Expect(err).NotTo(HaveOccurred())
		}
	}

	t.Log("It should update the .env")
	{
		// UpdateExamplesEnv adds a .env if there isn't one to the Examples repo. We need this because on Circle CI there
		// won't be a .env file, as it should be in .gitignore.
		ok, err := exists("examples-go/.env")
		g.Expect(err).NotTo(HaveOccurred())

		if !ok {
			env := getEnv()
			err = ioutil.WriteFile("examples-go/.env", []byte(env), 0644)
			g.Expect(err).NotTo(HaveOccurred())
			err = ioutil.WriteFile("client/.env", []byte(env), 0644)
			g.Expect(err).NotTo(HaveOccurred())
		} else {
			godotenv.Load("examples-go/.env")

			env := getEnv()
			err = ioutil.WriteFile("client/.env", []byte(env), 0644)
			g.Expect(err).NotTo(HaveOccurred())
		}
	}

	t.Log("It should listen")
	{
		args := []string{"listen", "--app_id=" + os.Getenv("ZENATON_APP_ID"),
			"--api_token=" + os.Getenv("ZENATON_API_TOKEN"),
			"--app_env=" + os.Getenv("ZENATON_APP_ENV"),
			"--boot=boot/boot.go"}

		cmd := exec.Command("zenaton", args...)

		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("problem listening: ", string(out))
		}
		g.Expect(err).NotTo(HaveOccurred())
	}

	t.Log("It should set up the log files")
	{
		// we sleep for 10 seconds after the listen on CircleCI in case we need to process any tasks that were still in the queue.
		// If so those logs will be added here and we'll have too many logs
		onCircleCI := os.Getenv("CircleCI")
		if onCircleCI == "" {
			time.Sleep(10 * time.Second)
		}

		var err error
		outFile, err = os.OpenFile("zenaton.out", os.O_RDWR, 0660)
		switch err.(type) {
		case *os.PathError:
			//this is ok
			outFile, err = os.Create("zenaton.out")
			g.Expect(err).NotTo(HaveOccurred())
		default:
			err = outFile.Truncate(0)
			g.Expect(err).NotTo(HaveOccurred())
		}

		errFile, err = os.OpenFile("zenaton.err", os.O_RDWR, 0660)
		switch err.(type) {
		case *os.PathError:
			//this is ok
			errFile, err = os.Create("zenaton.err")
			g.Expect(err).NotTo(HaveOccurred())

		default:
			err = errFile.Truncate(0)
			g.Expect(err).NotTo(HaveOccurred())
		}

		logFile, err := os.OpenFile("zenaton.log", os.O_RDWR, 0660)
		if err == nil {
			err = logFile.Truncate(0)
			g.Expect(err).NotTo(HaveOccurred())
		}
	}
}


func TestRunExamples(t *testing.T) {
	for _, entry := range table {
		entry := entry //gotcha!

		t.Run("", func(st *testing.T) {
			st.Parallel()
			g := NewGomegaWithT(st)

			// updates paths in copied examples so that imports work correctly
			if entry.exampleDir == "examples-go" {
				filePath := "examples-go/" + entry.directory + "/main.go"
				contents, err := ioutil.ReadFile(filePath)
				g.Expect(err).ToNot(HaveOccurred())
				newContents := bytes.Replace(contents, []byte("github.com/zenaton/examples-go"), []byte("github.com/zenaton/zenaton-go/integration/go/examples-go"), -1)

				err = ioutil.WriteFile(filePath, newContents, 0660)
			}

			cmd := exec.Command("go", "run", "-race", "main.go")
			cmd.Dir = entry.exampleDir + "/" + entry.directory

			out, err := cmd.CombinedOutput()

			g.Expect(string(out)).To(ContainSubstring(entry.fileRunOutput))

			if entry.fileRunErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				if err != nil {
					fmt.Println("problem running example: ", string(out))
				}
				g.Expect(err).ToNot(HaveOccurred())
			}
		})
	}
}

func TestExamplesOutputs(t *testing.T) {

	g := NewGomegaWithT(t)
	err := waitForLogs()
	g.Expect(err).ToNot(HaveOccurred())

	errLog, err := ioutil.ReadFile("zenaton.err")
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(string(errLog)).To(Equal(""))

	outLog, err := ioutil.ReadFile("zenaton.out")
	g.Expect(err).ToNot(HaveOccurred())

	for _, entry := range table {
		entry := entry
		t.Run(entry.context, func(t *testing.T){
			g := NewGomegaWithT(t)
			t.Log(entry.it)
			{

				var matchers []types.GomegaMatcher
				for _, possibleOut := range entry.possibleLogs {
					matchers = append(matchers, Equal(possibleOut))
				}

				var logs string
				if entry.directory == "recursive" {
					logs = getRecursiveLogs(string(outLog))
				} else {
					logs = getFilteredLogs(string(outLog), "::" + entry.directory+": ")
				}

				g.Expect(logs).To(SatisfyAny(matchers...))
			}
		})
	}
}

//waitForLogs waits in 10 second increments and checks to see if the log file had been updated in that time. If not,
//it returns.
func waitForLogs() error {
	var outFileLogs []byte

	for {
		time.Sleep(15 * time.Second)
		newOutFileLogs, err := ioutil.ReadFile("zenaton.out")
		if err != nil {
			return err
		}
		if string(newOutFileLogs) == string(outFileLogs) {
			break
		}
		outFileLogs = newOutFileLogs
	}
	return nil
}

// getFilteredLogs scans the log file looking for prefixes to lines. If it finds the "filter" prefix, it adds that line
// to the returned lines (removing the filter)
func getFilteredLogs(logs string, filter string) string {

	lines := strings.Split(logs, "\n")

	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, filter) {
			newLine := strings.TrimLeft(line, filter)
			newLines = append(newLines, newLine)
		}
	}

	return strings.Join(newLines, "\n") + "\n"
}

// getFilteredLogs scans the log file looking for prefixes to lines. If it finds the "filter" prefix, it adds that line
// to the returned lines (removing the filter)
func getRecursiveLogs(logs string) string {

	lines := strings.Split(logs, "\n")

	var newLines []string
	for _, line := range lines {
		if strings.HasPrefix(line, "::") {
			continue
		}
		newLines = append(newLines, line)
	}
	out := strings.Join(newLines, "")
	out = strings.Join(strings.Fields(out), "")
	return out
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func getEnv() string {
	zenatonAppID := os.Getenv("ZENATON_APP_ID")
	zenatonAPIToken := os.Getenv("ZENATON_API_TOKEN")
	zenatonAppEnv := os.Getenv("ZENATON_APP_ENV")
	zenatonConcurrentMax := os.Getenv("ZENATON_CONCURRENT_MAX")
	zenatonHandleOnly := os.Getenv("ZENATON_HANDLE_ONLY")
	zenatonHandleExcept := os.Getenv("ZENATON_HANDLE_EXCEPT")

	env := "ZENATON_APP_ID=" + zenatonAppID + "\n" +
		"ZENATON_API_TOKEN=" + zenatonAPIToken + "\n" +
		"ZENATON_APP_ENV=" + zenatonAppEnv + "\n" +
		"ZENATON_CONCURRENT_MAX=" + zenatonConcurrentMax + "\n" +
		"ZENATON_HANDLE_ONLY=" + zenatonHandleOnly + "\n" +
		"ZENATON_HANDLE_EXCEPT=" + zenatonHandleExcept

	return env
}
