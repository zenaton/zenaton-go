package task_test

import (
	"fmt"
	"strconv"
	"time"

	"github.com/onsi/gomega/types"
	"github.com/zenaton/zenaton-go/v1/zenaton/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	secondsInADay   = 86400
	secondsInAnHour = 3600
	secondsInAWeek  = 604800
	fakeYear        = 2018
	fakeMonth       = 10
	fakeDay         = 30
	fakeHour        = 17
	fakeMin         = 4
	fakeSec         = 5
	fakeNsec        = 0
	fakeLocation    = "America/New_York"
)

// fakeDate is a Tuesday
var fakeDate = time.Date(fakeYear, fakeMonth, fakeDay, fakeHour, fakeMin, fakeSec, fakeNsec, getFakeLocation())

func getFakeLocation() *time.Location {
	loc, err := time.LoadLocation(fakeLocation)
	if err != nil {
		panic(err)
	}
	return loc
}

var _ = BeforeSuite(func() {
	task.Now = func() time.Time {
		return fakeDate
	}
})

var _ = Describe("WaitTask", func() {

	Context("When executing a waitTask", func() {
		It("should wait", func() {
			execution := task.Wait().Seconds(1).Execute()

			//didn't receive an event
			Expect(execution.EventReceived()).To(BeFalse())
		})
	})

	//Context("When waiting for an event", func(){
	//	It("should wait until receiving the event", func(){
	//		now1 := time.Now()
	//		execution := task.Wait().Seconds(1).Execute()
	//		Expect(time.Now().Unix() - 1).To(Equal(now1.Unix()))
	//
	//		//didn't receive an event
	//		Expect(execution.EventReceived()).To(BeFalse())
	//	})
	//})

	Context("When waiting for an event", func() {
		It("should be able to get the name of the event", func() {
			w := task.Wait().ForEvent("eventName")
			Expect(w.Event()).To(Equal("eventName"))
		})
	})

	Context("Duration", func() {

		It("should Wait for a second", func() {
			w := task.Wait().Seconds(1)
			expectDurationInSeconds(w, 1)
		})
		It("should Wait for a minute", func() {
			w := task.Wait().Minutes(1)
			expectDurationInSeconds(w, 60)
		})
		It("should Wait for an hour", func() {
			w := task.Wait().Hours(1)
			expectDurationInSeconds(w, 3600)
		})
		It("should Wait for a day", func() {
			w := task.Wait().Days(1)
			expectDurationInSeconds(w, secondsInADay)
		})
		It("should Wait for a week", func() {
			w := task.Wait().Weeks(1)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))

			matchers := []types.GomegaMatcher{
				Equal(int64(secondsInAWeek)),
				Equal(int64(secondsInAWeek + secondsInAnHour)), // fall back from daylight savings
				Equal(int64(secondsInAWeek - secondsInAnHour)), // spring forward from daylight savings
			}

			Expect(duration).To(SatisfyAny(matchers...))
		})
		It("should Wait for a month", func() {
			w := task.Wait().Months(1)
			now := task.Now()
			then := now.AddDate(0, 1, 0)
			fmt.Println("test then: ", then, "now: ", now)
			expectDurationInSeconds(w, int(then.Unix()-now.Unix()))
		})
		It("should Wait for a year", func() {
			w := task.Wait().Years(1)
			now := task.Now()
			then := now.AddDate(1, 0, 0)
			expectDurationInSeconds(w, int(then.Unix()-now.Unix()))
		})
	})

	Context("Timezone", func() {
		It("should be able to set timezone which effects the timestamp", func() {
			w := task.Wait().Monday(2).At("5")

			localTimestamp, _, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())

			// change timezone
			err = w.Timezone("Indian/Maldives")
			Expect(err).NotTo(HaveOccurred())

			maldivesTimestamp, _, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())

			Expect(localTimestamp).NotTo(Equal(maldivesTimestamp))
		})

		Context("When using a non-existent timezone", func() {
			It("should return an error", func() {
				err := task.Wait().Timezone("non-existent time zone")
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Context("Timestamp", func() {
		It("should Wait for the specified timestamp", func() {

			w := task.Wait().Timestamp(fakeDate.Unix() + 3)
			expectTimestampEqualTime(w, fakeDate.Add(3*time.Second))
		})

		Context("when using the timestamp method twice", func() {
			It("should error", func() {
				now := time.Now().Unix()

				w := task.Wait().Timestamp(now + 1).Timestamp(now + 3)
				expectError(w, "incompatible definition in WaitTask methods")
			})
		})

		Context("when using the timestamp method with any other method", func() {
			It("should error", func() {
				w := task.Wait().Timestamp(time.Now().Unix() + 3).Tuesday(1)
				expectError(w, "timestamp can only be used alone")
			})
		})
	})

	Context("At", func() {

		Context("when using At with just an hour", func() {
			Context("when the given hour is the same as the current hour", func() {
				It("should wait until that hour the next day", func() {

					w := task.Wait().At(strconv.Itoa(fakeHour))
					expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 1).Add(-fakeMin*time.Minute-fakeSec*time.Second))
				})
			})
			Context("when the given hour is after the current hour", func() {
				It("should wait until the hour on the same day", func() {

					w := task.Wait().At(fakeTimePlus1Hour())
					expectTimestampEqualTime(w, fakeDate.Add(time.Hour))
				})
			})
		})

		Context("when using At with just an hour and minute", func() {
			It("should wait until the the hour:minute", func() {

				hour := fakeHour + 1
				min := fakeMin + 1
				w := task.Wait().At(strconv.Itoa(hour) + ":" + strconv.Itoa(min))
				expectTimestampEqualTime(w, fakeDate.Add(time.Hour+time.Minute-(fakeSec*time.Second)))
			})
		})

		Context("when using At with an hour, minute, and second", func() {
			It("should wait until the next hour:minute:second", func() {

				hour := fakeHour + 1
				min := fakeMin + 1
				sec := fakeSec + 1
				w := task.Wait().At(strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec))
				expectTimestampEqualTime(w, fakeDate.Add(time.Hour+time.Minute+time.Second))
			})
		})

		Context("when using At in conjunction with a weekday method", func() {
			Context("when the specified weekday is today and the specified time is later than now", func() {
				It("should wait until the specified time today", func() {

					//fakeDate is a Tuesday
					w := task.Wait().Tuesday(1).At(fakeTimePlus1Hour())
					expectTimestampLater(w, time.Hour)
				})
			})

			Context("when the specified weekday is NOT today", func() {
				It("should wait until the specified time and weekday", func() {

					//fakeDate is a Tuesday
					w := task.Wait().Wednesday(1).At(fakeTimePlus1Hour())

					expectTimestampLater(w, 25*time.Hour)
				})
			})
		})

		Context("when using At in conjunction with a dayOfMonth method", func() {
			Context("when the specified day is today and the specified time is later than now", func() {
				It("should wait until the specified time today", func() {

					w := task.Wait().DayOfMonth(fakeDay).At(fakeTimePlus1Hour())

					expectTimestampLater(w, time.Hour)
				})
			})

			Context("when the specified day is NOT today", func() {
				It("should wait until the specified time and day", func() {

					//fakeDate is a Tuesday
					w := task.Wait().DayOfMonth(fakeDay + 1).At(fakeTimePlus1Hour())

					expectTimestampLater(w, 25*time.Hour)
				})
			})
		})

		Context("when calling At twice", func() {
			It("should return an error", func() {
				w := task.Wait().At("5").At("4")
				expectError(w, "incompatible definition in WaitTask methods")
			})
		})

		Context("when calling At with incorrectly formatted time", func() {
			It("should return an error", func() {
				w := task.Wait().At("")
				expectError(w, "time formatted incorrectly")

				w = task.Wait().At("x")
				expectError(w, "time formatted incorrectly")

				w = task.Wait().At("5:x")
				expectError(w, "time formatted incorrectly")

				w = task.Wait().At("5:03:x")
				expectError(w, "time formatted incorrectly")
			})
		})
	})

	Context("DayOfMonth", func() {
		It("should wait until the next specified day of the month (same time)", func() {

			w := task.Wait().DayOfMonth(fakeDay + 1)
			expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 1))
		})

		Context("when the specified day is the same as today", func() {
			It("should wait for one month", func() {

				w := task.Wait().DayOfMonth(fakeDay)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 1, 0))
			})
		})

		Context("twice", func() {
			It("should return an error", func() {
				w := task.Wait().DayOfMonth(4).DayOfMonth(6)
				expectError(w, "incompatible definition in WaitTask methods")
			})
		})
	})

	Context("Weekdays", func() {

		Context("when using Monday", func() {
			It("should wait until the nth Monday (same time)", func() {
				w := task.Wait().Monday(2)
				//fakeDate is a tuesday
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 13))
			})
		})

		Context("when using Tuesday", func() {
			It("should wait until the nth Tuesday (same time)", func() {
				w := task.Wait().Tuesday(2)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 14))
			})
		})

		Context("when using Wednesday", func() {
			It("should wait until the nth Wednesday (same time)", func() {
				w := task.Wait().Wednesday(2)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 8))
			})
		})

		Context("when using Thursday", func() {
			It("should wait until the nth Thursday (same time)", func() {
				w := task.Wait().Thursday(2)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 9))
			})
		})

		Context("when using Friday", func() {
			It("should wait until the nth Friday (same time)", func() {
				w := task.Wait().Friday(2)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 10))
			})
		})

		Context("when using Saturday", func() {
			It("should wait until the nth Saturday (same time)", func() {
				w := task.Wait().Saturday(2)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 11))
			})
		})

		Context("when using Sunday", func() {
			It("should wait until the nth Sunday (same time)", func() {
				w := task.Wait().Sunday(2)
				expectTimestampEqualTime(w, fakeDate.AddDate(0, 0, 12))
			})
		})

	})
})

func fakeTimePlus1Hour() string {
	hour := fakeHour + 1
	min := fakeMin
	sec := fakeSec
	return strconv.Itoa(hour) + ":" + strconv.Itoa(min) + ":" + strconv.Itoa(sec)
}

func expectTimestampLater(w *task.WaitTask, dif time.Duration) {
	err := w.Timezone(fakeLocation)
	Expect(err).NotTo(HaveOccurred())

	timestamp, duration, err := w.GetTimestampOrDuration()
	Expect(err).NotTo(HaveOccurred())

	expectedDate := fakeDate.Add(dif)
	Expect(timestamp, 0).To(Equal(expectedDate.Unix()))
	Expect(duration).To(Equal(int64(0)))
}

func expectDurationInSeconds(w *task.WaitTask, seconds int) {
	timestamp, duration, err := w.GetTimestampOrDuration()
	fmt.Println("duration:", duration)
	Expect(err).NotTo(HaveOccurred())
	Expect(timestamp).To(Equal(int64(0)))
	Expect(duration).To(Equal(int64(seconds)))
}

func expectTimestampEqualTime(w *task.WaitTask, date time.Time) {

	err := w.Timezone(fakeLocation)
	Expect(err).NotTo(HaveOccurred())

	timestamp, duration, err := w.GetTimestampOrDuration()
	Expect(err).NotTo(HaveOccurred())

	Expect(duration).To(Equal(int64(0)))
	Expect(timestamp).To(Equal(date.Unix()))
}

func expectError(w *task.WaitTask, errMessage string) {
	_, _, err := w.GetTimestampOrDuration()
	Expect(err.Error()).To(Equal(errMessage))
}
