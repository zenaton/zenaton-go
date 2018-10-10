package integration_test

import (
	"fmt"
	"github.com/onsi/gomega/types"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/zenaton/zenaton-go/integration/go"

	"testing"

	. "github.com/onsi/gomega"
)

type entry struct {
	exampleDir    string
	directory     string
	context       string
	it            string
	fileRunOutput string
	fileRunErr    bool
	outFile       []string
	errFile       string
	sleep         int //time to wait for the file to be written to
}

var table = []entry{
	//{
	//	exampleDir: "testExamples",
	//	directory:  "testFindWorkflow",
	//	context:    "with a workflow that has an ID method",
	//	it:         "should be able to find a workflow and dispatch it again",
	//	outFile:    []string{"out:  testTaskReturn\nerr:  <nil>\nout:  testTaskReturn\nerr:  <nil>\n"},
	//	errFile:    "",
	//	sleep:      16,
	//},
	//{
	//	exampleDir: "testExamples",
	//	directory:  "testReturnFromTaskInsideTask",
	//	context:    "with a task that launches another task",
	//	it:         "the outer task should be able to get the return value of the inner task",
	//	outFile:    []string{"out:  test return\nerr:  <*>\"test error\"\n"},
	//	errFile:    "",
	//	sleep:      10,
	//},
	//{
	//	exampleDir: "testExamples",
	//	directory:  "testReturnError",
	//	context:    "with a task that returns an error",
	//	it:         "the error should be returned to the workflow",
	//	outFile:    []string{"out:  <nil>\nerr:  <*>\"testTaskError\"\n"},
	//	errFile:    "",
	//	sleep:      8,
	//},
	{
		exampleDir: "examples-go",
		directory:  "waitevent",
		context:    "with a waitevent workflow",
		it:         "should wait for either an event or a given time",
		outFile:    []string{"Task A starts\nTask A ends\n"},
		errFile:    "",
		sleep:      15,
	},
	//{
	//	exampleDir: "examples-go",
	//	directory:  "event",
	//	context:    "with an event based workflow",
	//	it:         "should handle events",
	//	outFile:    []string{"Task A starts\nTask A ends\nTask B starts\nTask B ends\n"},
	//	errFile:    "",
	//	sleep:      15,
	//},
	//{
	//	exampleDir: "examples-go",
	//	directory:  "wait",
	//	context:    "with a wait workflow",
	//	it:         "should wait before a task",
	//	outFile:    []string{"Task A starts\nTask A ends\nTask B starts\nTask B ends\n"},
	//	errFile:    "",
	//	sleep:      22,
	//},
	//{
	//	exampleDir: "examples-go",
	//	directory:  "asynchronous",
	//	context:    "with an asychronous workflow",
	//	it:         "should perform tasks asynchronously",
	//	outFile:    []string{"Task A starts\nTask B starts\nTask C starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
	//							"Task A starts\nTask C starts\nTask B starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
	//							"Task B starts\nTask A starts\nTask C starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
	//							"Task B starts\nTask C starts\nTask A starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
	//							"Task C starts\nTask A starts\nTask B starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",
	//							"Task C starts\nTask B starts\nTask A starts\nTask A ends\nTask B ends\nTask C ends\nTask D starts\nTask D ends\n",},
	//	errFile:    "",
	//	sleep:      21,
	//},
	//{
	//	exampleDir: "examples-go",
	//	directory:  "parallel",
	//	context:    "with a parallel workflow",
	//	it:         "should handle tasks in parallel",
	//	outFile:    []string{"Task A starts\nTask B starts\nTask A ends\nTask B ends\nTask D starts\nTask D ends\n",
	//						"Task B starts\nTask A starts\nTask A ends\nTask B ends\nTask D starts\nTask D ends\n"},
	//	errFile:    "",
	//	sleep:      19,
	//},
	//{
	//	exampleDir: "examples-go",
	//	directory:  "recursive",
	//	context:    "with a recursive workflow",
	//	it:         "should handle tasks recursively",
	//	outFile:    []string{"01\nIteration: 1\n01\nIteration: 2\n01"},
	//	errFile:    "",
	//	sleep:      16,
	//},
	//{
	//	exampleDir: "examples-go",
	//	directory:  "sequential",
	//	context:    "with a sequential workflow",
	//	it:         "should handle tasks sequentially",
	//	outFile:    []string{"Task A starts\nTask A ends\nTask C starts\nTask C ends\nTask D starts\nTask D ends\n"},
	//	errFile:    "",
	//	sleep:      28,
	//},
	//{ // this test is too fickle for now
	//	exampleDir: "examples-go",
	//	directory:  "version",
	//	context:    "with a version workflow",
	//	it:         "should handle versioned workflows",
	//	outFile:    []string{"Task A starts\nTask B starts\nTask C starts\nTask D starts\nTask A ends\nTask B ends\nTask C ends\nTask D ends\n",
	//						"Task B starts\nTask A starts\nTask C starts\nTask D starts\nTask A ends\nTask B ends\nTask C ends\nTask D ends\n"},
	//	errFile:    "",
	//	sleep:      14,
	//},
}

func TestSetup(t *testing.T) {
	g := NewGomegaWithT(t)

	err := CopyExamples()
	g.Expect(err).NotTo(HaveOccurred())

	for _, entry := range table {
		entry := entry //gotcha!

		envFile := entry.directory + ".env"
		err := ChangeClient(entry.exampleDir, entry.directory+"/main.go", envFile)
		g.Expect(err).NotTo(HaveOccurred())

		created, err := AddEnv(envFile)
		g.Expect(err).NotTo(HaveOccurred())

		if created {
			err = WriteAppEnv(envFile, entry.directory)
			g.Expect(err).NotTo(HaveOccurred())
		}

		bootFileName := entry.directory + "boot.go"

		err = AddBoot(bootFileName, envFile)
		g.Expect(err).NotTo(HaveOccurred())

		err = Listen(envFile, bootFileName, entry.exampleDir, entry.directory)
		g.Expect(err).NotTo(HaveOccurred())
	}
}

func TestExamples(t *testing.T) {
	for _, entry := range table {
		entry := entry //gotcha!

		t.Run("", func(st *testing.T) {
			g := NewGomegaWithT(st)
			st.Parallel()
			entry := entry
			st.Log(entry.context)
			{
				st.Log(entry.it)
				{
					errFile, err := os.OpenFile(entry.exampleDir+"/"+entry.directory+"/zenaton.err", os.O_RDWR, 0660)
					switch err.(type) {
					case *os.PathError:
						//this is ok
						errFile, err = os.Create(entry.exampleDir + "/" + entry.directory + "/zenaton.err")
						g.Expect(err).ToNot(HaveOccurred())

					default:
						g.Expect(err).ToNot(HaveOccurred())

						//clear the files
						err = errFile.Truncate(0)
						g.Expect(err).ToNot(HaveOccurred())
						_, err = errFile.Seek(0, 0)
						g.Expect(err).ToNot(HaveOccurred())
					}
					defer errFile.Close()

					outFile, err := os.OpenFile(entry.exampleDir+"/"+entry.directory+"/zenaton.out", os.O_RDWR, 0660)
					switch err.(type) {
					case *os.PathError:
						//this is ok
						outFile, err = os.Create(entry.exampleDir + "/" + entry.directory + "/zenaton.out")
						g.Expect(err).ToNot(HaveOccurred())
					default:
						g.Expect(err).ToNot(HaveOccurred())

						err = outFile.Truncate(0)
						g.Expect(err).ToNot(HaveOccurred())
						_, err = outFile.Seek(0, 0)
						g.Expect(err).ToNot(HaveOccurred())
					}
					defer outFile.Close()

					logFile, err := os.OpenFile(entry.exampleDir+"/"+entry.directory+"/zenaton.log", os.O_RDWR, 0660)
					switch err.(type) {
					case *os.PathError:
						//this is ok
					default:
						g.Expect(err).ToNot(HaveOccurred())

						err = logFile.Truncate(0)
						g.Expect(err).ToNot(HaveOccurred())
						_, err = logFile.Seek(0, 0)
						g.Expect(err).ToNot(HaveOccurred())
					}
					defer logFile.Close()

					cmd := exec.Command("go", "run", "-race", "main.go")
					cmd.Dir = entry.exampleDir + "/" + entry.directory

					out, err := cmd.CombinedOutput()

					g.Expect(string(out)).To(ContainSubstring(entry.fileRunOutput))

					if entry.fileRunErr {
						g.Expect(err).To(HaveOccurred())
					} else {
						if err != nil {
							fmt.Println("out: ", string(out))
						}
						g.Expect(err).ToNot(HaveOccurred())
					}

					time.Sleep(time.Duration(entry.sleep) * time.Second)

					errLog, err := ioutil.ReadAll(errFile)
					g.Expect(err).ToNot(HaveOccurred())
					outLog, err := ioutil.ReadAll(outFile)
					g.Expect(err).ToNot(HaveOccurred())

					g.Expect(string(errLog)).To(Equal(entry.errFile))

					var matchers []types.GomegaMatcher
					for _, possibleOut := range entry.outFile {
						matchers = append(matchers, Equal(possibleOut))
					}

					g.Expect(string(outLog)).To(SatisfyAny(matchers...))

				}
			}
		})
	}
}
