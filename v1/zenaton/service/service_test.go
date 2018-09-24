package service_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zenaton/zenaton-go/v1/zenaton/service"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

)

var _ = BeforeSuite(func(){
	os.Setenv("ZENATON_LOG_LEVEL", "0")
	createServer(8080, 100*time.Microsecond)
})

var _ = Describe("http", func() {
	// a bit of context for this test. Before making the change where connections are not reused, we would sometimes
	// get an EOF or "read: connection reset by peer" error.
	url := "http://127.0.0.1:8080"
	It("should be able to POST without an error", func() {

		for i := 0; i < 60; i++ {
			time.Sleep(100 * time.Microsecond)
			resp, err := service.Post(url, `{"key":"value"}`)
			Expect(err).NotTo(HaveOccurred())

			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("should be able to PUT without an error", func() {

		for i := 0; i < 60; i++ {
			time.Sleep(100 * time.Microsecond)
			resp, err := service.Put(url, `{"key":"value"}`)
			Expect(err).NotTo(HaveOccurred())

			_, err = ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
		}
	})
})


func createServer(port int, idleTime time.Duration) error {
	server := &http.Server{
		Handler: http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if _, err := w.Write([]byte(`{"key":"value"}`)); err != nil {
					panic(fmt.Sprint("error writing to client: %s", err))
				}
			},
		),
		IdleTimeout: idleTime,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return fmt.Errorf("error listening: %s", err)
	}

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(fmt.Sprint("error serving: %s", err))
		}
		if err := listener.Close(); err != nil {
			panic(fmt.Sprint("error closing listener: %s", err))
		}
	}()

	return nil
}