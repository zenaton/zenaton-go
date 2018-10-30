package task_test

import (
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
)

var _ = Describe("WaitTask", func() {

	Context("When executing a waitTask", func(){
		It("should wait", func(){
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

	Context("When waiting for an event", func(){
		It("should be able to get the name of the event", func(){
			w := task.Wait().ForEvent("eventName")
			Expect(w.Event()).To(Equal("eventName"))
		})
	})

	Context("when applying duration methods", func() {
		It("should Wait for a second", func() {
			w := task.Wait().Seconds(1)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))
			Expect(duration).To(Equal(int64(1)))
		})
		It("should Wait for a minute", func() {
			w := task.Wait().Minutes(1)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))
			Expect(duration).To(Equal(int64(60)))
		})
		It("should Wait for an hour", func() {
			w := task.Wait().Hours(1)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))
			Expect(duration).To(Equal(int64(3600)))
		})
		It("should Wait for a day", func() {
			w := task.Wait().Days(1)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))
			Expect(duration).To(Equal(int64(secondsInADay)))
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
			now := time.Now()
			then := now.AddDate(0, 1, 0)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))
			Expect(duration).To(Equal(then.Unix() - now.Unix()))
		})
		It("should Wait for a year", func() {
			w := task.Wait().Years(1)
			now := time.Now()
			then := now.AddDate(1, 0, 0)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(0)))
			Expect(duration).To(Equal(then.Unix() - now.Unix()))
		})
	})

	Context("when applying timestamp methods", func() {

		It("should be able to set timezone which effects the timestamp", func() {
			w := task.Wait().Monday(2).At("5")

			localTimestamp, _, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())

			// change timezone
			err = w.Timezone("Indian/Maldives")
			Expect(err).NotTo(HaveOccurred())

			maldivesTimestamp, _, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())

			if localTimestamp == maldivesTimestamp { // maybe our local timezone is actually "Indian/Maldives"
				err := w.Timezone("Indian/Mauritius")
				Expect(err).NotTo(HaveOccurred())

				mauritiusTimestamp, _, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(localTimestamp).NotTo(Equal(mauritiusTimestamp))
			}
		})

		Context("When using a non-existent timezone", func() {
			It("should return an error", func() {
				err := task.Wait().Timezone("non-existent time zone")
				Expect(err).To(HaveOccurred())
			})
		})

		It("should Wait for the specified timestamp", func() {
			now := time.Now().Unix()
			w := task.Wait().Timestamp(now + 3)
			timestamp, duration, err := w.GetTimestampOrDuration()
			Expect(err).NotTo(HaveOccurred())
			Expect(timestamp).To(Equal(int64(now + 3)))
			Expect(duration).To(Equal(int64(0)))
		})

		Context("when using the timestamp method twice", func() {
			It("should error", func() {
				now := time.Now().Unix()

				w := task.Wait().Timestamp(now + 1).Timestamp(now + 3)
				_, _, err := w.GetTimestampOrDuration()
				Expect(err.Error()).To(Equal("incompatible definition in WaitTask methods"))
			})
		})

		Context("At", func() {

			Context("when using At with just an hour", func() {
				Context("when the given hour is the same as the current hour", func() {
					It("should wait until the next day", func() {
						now := time.Now()
						hour := strconv.Itoa(now.Hour())

						w := task.Wait().At(hour)

						timestamp, duration, err := w.GetTimestampOrDuration()
						Expect(err).NotTo(HaveOccurred())

						Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
						Expect(time.Unix(timestamp, 0).Minute()).To(Equal(0))
						Expect(time.Unix(timestamp, 0).Second()).To(Equal(0))
						Expect(time.Unix(timestamp, 0).Day()).To(Equal(now.AddDate(0, 0, 1).Day()))

						Expect(duration).To(Equal(int64(0)))

					})
				})

				It("should wait until that hour", func() {
					now := time.Now()

					w := task.Wait().At("4")

					timestamp, duration, err := w.GetTimestampOrDuration()
					Expect(err).NotTo(HaveOccurred())

					Expect(time.Unix(timestamp, 0).Hour()).To(Equal(4))
					Expect(time.Unix(timestamp, 0).Minute()).To(Equal(0))
					Expect(time.Unix(timestamp, 0).Second()).To(Equal(0))

					if now.Hour() < 4 {
						Expect(time.Unix(timestamp, 0).Day()).To(Equal(now.Day()))
					} else {
						Expect(time.Unix(timestamp, 0).Day()).To(Equal(now.AddDate(0, 0, 1).Day()))
					}

					Expect(duration).To(Equal(int64(0)))
				})

			})
			Context("when using At with just an hour and minute", func() {
				It("should wait until the the hour:minute", func() {
					now := time.Now().Unix()
					w := task.Wait().At("8:30")
					timestamp, duration, err := w.GetTimestampOrDuration()

					Expect(err).NotTo(HaveOccurred())
					Expect(time.Unix(timestamp, 0).Hour()).To(Equal(8))
					Expect(time.Unix(timestamp, 0).Minute()).To(Equal(30))
					Expect(time.Unix(timestamp, 0).Second()).To(Equal(0))

					// that is, it should happen some time in the next 24 hours
					Expect(now + secondsInADay).To(BeNumerically(">", timestamp))
					Expect(now).To(BeNumerically("<", timestamp))

					Expect(duration).To(Equal(int64(0)))
				})
			})
			Context("when using At with an hour, minute, and second", func() {
				It("should wait until the next hour:minute:second", func() {
					now := time.Now().Unix()
					w := task.Wait().At("8:30:12")
					timestamp, duration, err := w.GetTimestampOrDuration()

					Expect(err).NotTo(HaveOccurred())
					Expect(time.Unix(timestamp, 0).Hour()).To(Equal(8))
					Expect(time.Unix(timestamp, 0).Minute()).To(Equal(30))
					Expect(time.Unix(timestamp, 0).Second()).To(Equal(12))

					// that is, it should happen some time in the next 24 hours
					Expect(now + secondsInADay).To(BeNumerically(">", timestamp))
					Expect(now).To(BeNumerically("<", timestamp))

					Expect(duration).To(Equal(int64(0)))
				})
			})

			Context("when using At in conjunction with a weekday method", func() {
				It("should wait until the specified time and weekday", func() {
					w := task.Wait().Monday(1).At("8:30:12")

					timestamp, duration, err := w.GetTimestampOrDuration()
					Expect(err).NotTo(HaveOccurred())

					Expect(time.Unix(timestamp, 0).Hour()).To(Equal(8))
					Expect(time.Unix(timestamp, 0).Minute()).To(Equal(30))
					Expect(time.Unix(timestamp, 0).Second()).To(Equal(12))
					Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Monday))

					Expect(duration).To(Equal(int64(0)))
				})
			})

			Context("when using At in conjunction with a dayOfMonth method", func() {
				It("should wait until the specified time and day", func() {
					w := task.Wait().DayOfMonth(23).At("8:30:12")

					timestamp, duration, err := w.GetTimestampOrDuration()
					Expect(err).NotTo(HaveOccurred())

					Expect(time.Unix(timestamp, 0).Hour()).To(Equal(8))
					Expect(time.Unix(timestamp, 0).Minute()).To(Equal(30))
					Expect(time.Unix(timestamp, 0).Second()).To(Equal(12))
					Expect(time.Unix(timestamp, 0).Day()).To(Equal(23))

					Expect(duration).To(Equal(int64(0)))
				})
			})

			Context("when calling At twice", func() {
				It("should return an error", func() {
					w := task.Wait().At("5").At("4")
					_, _, err := w.GetTimestampOrDuration()
					Expect(err.Error()).To(Equal("incompatible definition in WaitTask methods"))
				})
			})

			Context("when calling At with incorrectly formatted time", func() {
				It("should return an error", func() {
					w := task.Wait().At("")
					_, _, err := w.GetTimestampOrDuration()
					Expect(err.Error()).To(Equal("time formatted incorrectly"))

					w = task.Wait().At("x")
					_, _, err = w.GetTimestampOrDuration()
					Expect(err.Error()).To(Equal("time formatted incorrectly"))

					w = task.Wait().At("5:x")
					_, _, err = w.GetTimestampOrDuration()
					Expect(err.Error()).To(Equal("time formatted incorrectly"))

					w = task.Wait().At("5:03:x")
					_, _, err = w.GetTimestampOrDuration()
					Expect(err.Error()).To(Equal("time formatted incorrectly"))
				})
			})

			Context("when calling At after a weekday", func() {
				It("should wait until the given time on the given weekday", func() {
					w := task.Wait().Monday(1).At("12:30")
					timestamp1, duration, err := w.GetTimestampOrDuration()
					Expect(err).NotTo(HaveOccurred())

					Expect(duration).To(Equal(int64(0)))
					Expect(time.Unix(timestamp1, 0).Weekday()).To(Equal(time.Monday))
					Expect(time.Unix(timestamp1, 0).Hour()).To(Equal(12))
					Expect(time.Unix(timestamp1, 0).Minute()).To(Equal(30))
					Expect(timestamp1).To(BeNumerically("<", time.Now().AddDate(0, 0, 8).Unix()))
				})

			})

			Context("When calling At after a DayOfMonth", func() {
				Context("when the given day of the month is the same as today", func() {
					It("should wait until the given time of the next month", func() {

						now := time.Now()
						day := now.Day()
						hour := now.Hour() - 1
						strHour := strconv.Itoa(hour)

						w := task.Wait().DayOfMonth(day).At(strHour)
						timestamp1, duration, err := w.GetTimestampOrDuration()
						Expect(err).NotTo(HaveOccurred())

						Expect(duration).To(Equal(int64(0)))

						Expect(time.Unix(timestamp1, 0).Day()).To(Equal(day))
						Expect(time.Unix(timestamp1, 0).Hour()).To(Equal(hour))
						Expect(time.Unix(timestamp1, 0).Minute()).To(Equal(0))
						Expect(time.Unix(timestamp1, 0).Month()).To(Equal(now.AddDate(0, 1, 0).Month()))
					})
				})

				It("should wait until the given time and day of the month", func() {

					now := time.Now()

					w := task.Wait().DayOfMonth(23).At("12:30")
					timestamp1, duration, err := w.GetTimestampOrDuration()
					Expect(err).NotTo(HaveOccurred())

					Expect(duration).To(Equal(int64(0)))

					Expect(time.Unix(timestamp1, 0).Day()).To(Equal(23))
					Expect(time.Unix(timestamp1, 0).Hour()).To(Equal(12))
					Expect(time.Unix(timestamp1, 0).Minute()).To(Equal(30))
					//be within the next month
					Expect(timestamp1).To(BeNumerically("<", now.AddDate(0, 1, 1).Unix()))
				})
			})
		})

		Context("when using DayOfMonth", func() {
			It("should wait until the next specified day of the month (same time)", func() {
				now := time.Now()
				w := task.Wait().DayOfMonth(now.Day())

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))

				// that is, it should happen some time in the next month
				Expect(now.AddDate(0, 1, 1).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})

			Context("twice", func() {
				It("should return an error", func() {
					w := task.Wait().DayOfMonth(4).DayOfMonth(6)
					_, _, err := w.GetTimestampOrDuration()
					Expect(err.Error()).To(Equal("incompatible definition in WaitTask methods"))
				})
			})
		})

		Context("when using Monday", func() {
			It("should wait until the next Monday (same time)", func() {
				now := time.Now()
				w := task.Wait().Monday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Monday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})

		Context("when using Tuesday", func() {
			It("should wait until the next Tuesday (same time)", func() {
				now := time.Now()
				w := task.Wait().Tuesday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Tuesday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})

		Context("when using Wednesday", func() {
			It("should wait until the next Wednesday (same time)", func() {
				now := time.Now()
				w := task.Wait().Wednesday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Wednesday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})

		Context("when using Thursday", func() {
			It("should wait until the next Thursday (same time)", func() {
				now := time.Now()
				w := task.Wait().Thursday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Thursday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})

		Context("when using Friday", func() {
			It("should wait until the next Friday (same time)", func() {
				now := time.Now()
				w := task.Wait().Friday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Friday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})

		Context("when using Sunday", func() {
			It("should wait until the next Saturday (same time)", func() {
				now := time.Now()
				w := task.Wait().Saturday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Saturday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})

		Context("when using Sunday", func() {
			It("should wait until the next Sunday (same time)", func() {
				now := time.Now()
				w := task.Wait().Sunday(2)

				timestamp, duration, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())

				Expect(time.Unix(timestamp, 0).Hour()).To(Equal(now.Hour()))
				Expect(time.Unix(timestamp, 0).Minute()).To(Equal(now.Minute()))
				Expect(time.Unix(timestamp, 0).Second()).To(Equal(now.Second()))
				Expect(time.Unix(timestamp, 0).Weekday()).To(Equal(time.Sunday))

				// that is, it should happen some time in the next two weeks
				Expect(now.AddDate(0, 0, 15).Unix()).To(BeNumerically(">", timestamp))
				Expect(now.AddDate(0, 0, 7).Unix()).To(BeNumerically("<", timestamp))

				Expect(duration).To(Equal(int64(0)))
			})
		})
	})
	//
	//    subject { with_timestamp._get_timestamp_or_duration }
	//
	//    let(:today) { Time.utc(2018, 7, 13, 12, 2, 0) }
	//
	//    before { Timecop.freeze(today) }
	//
	//    after { Timecop.return }
	//
	//    context 'when there is no buffer' do
	//      it { is_expected.to eq [nil, nil] }
	//    end
	//
	//    context 'when applying duration methods' do
	//      before { with_timestamp.seconds }
	//
	//      it { is_expected.to eq [nil, 1] }
	//    end
	//
	//    context 'when specifying a timestamp' do
	//      let(:expected_timestamp) { 1522591200 }
	//
	//      before { with_timestamp.timestamp(1522591200) }
	//
	//      it { is_expected.to eq([expected_timestamp, nil]) }
	//    end
	//
	//    context 'when specifying an full hour' do
	//      let(:expected_time) { Time.utc(2018, 7, 13, 15, 10, 23) }
	//
	//      before { with_timestamp.at('15:10:23') }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying an hour without seconds' do
	//      let(:expected_time) { Time.utc(2018, 7, 13, 15, 10, 0) }
	//
	//      before { with_timestamp.at('15:10') }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying an hour without minutes or seconds' do
	//      let(:expected_time) { Time.utc(2018, 7, 13, 15, 0, 0) }
	//
	//      before { with_timestamp.at('15') }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying a day of the month' do
	//      let(:expected_time) { Time.utc(2018, 8, 12, 12, 2, 0) }
	//
	//      before { with_timestamp.on_day(12) }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next monday' do
	//      let(:expected_time) { Time.utc(2018, 7, 16, 12, 2, 0) }
	//
	//      before { with_timestamp.monday }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next tuesday' do
	//      let(:expected_time) { Time.utc(2018, 7, 17, 12, 2, 0) }
	//
	//      before { with_timestamp.tuesday(1) }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying second wednesday from now' do
	//      let(:expected_time) { Time.utc(2018, 7, 25, 12, 2, 0) }
	//
	//      before { with_timestamp.wednesday(2) }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next thursday' do
	//      let(:expected_time) { Time.utc(2018, 7, 19, 12, 2, 0) }
	//
	//      before { with_timestamp.thursday }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next friday' do
	//      let(:expected_time) { Time.utc(2018, 7, 20, 12, 2, 0) }
	//
	//      before { with_timestamp.friday }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next saturday' do
	//      let(:expected_time) { Time.utc(2018, 7, 14, 12, 2, 0) }
	//
	//      before { with_timestamp.saturday }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next sunday' do
	//      let(:expected_time) { Time.utc(2018, 7, 15, 12, 2, 0) }
	//
	//      before { with_timestamp.sunday }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next monday at 8:00AM' do
	//      let(:expected_time) { Time.utc(2018, 7, 16, 8, 0, 0) }
	//
	//      before { with_timestamp.monday.at('8:00') }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next 12th at 6PM' do
	//      let(:expected_time) { Time.utc(2018, 8, 12, 18, 0, 0) }
	//
	//      before { with_timestamp.on_day(12).at('18') }
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//  end
	//
	//  context 'with timezones' do
	//    subject { with_timestamp._get_timestamp_or_duration }
	//
	//    let(:timezone) { 'America/New_York' }
	//    let(:today) { Time.zone.local(2018, 7, 13, 12, 2, 0) }
	//
	//    before do
	//      klass.timezone = Time.zone = timezone
	//      Timecop.freeze(today)
	//    end
	//
	//    after do
	//      Timecop.return
	//      klass.timezone = Time.zone = nil
	//    end
	//
	//    context 'when there is no buffer' do
	//      it { is_expected.to eq [nil, nil] }
	//    end
	//
	//    context 'when applying duration methods' do
	//      before { with_timestamp.seconds }
	//
	//      it { is_expected.to eq [nil, 1] }
	//    end
	//
	//    context 'when specifying a timestamp' do
	//      let(:expected_timestamp) { 1522591200 }
	//
	//      before { with_timestamp.timestamp(1522591200) }
	//
	//      it { is_expected.to eq([expected_timestamp, nil]) }
	//    end
	//
	//    context 'when specifying an full hour' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 13, 15, 10, 23) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.at('15:10:23')
	//      end
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying an hour without seconds' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 13, 15, 10, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.at('15:10')
	//      end
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying an hour without minutes or seconds' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 13, 15, 0, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.at('15')
	//      end
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying a day of the month' do
	//      let(:expected_time) { Time.zone.local(2018, 8, 12, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.on_day(12)
	//      end
	//
	//      it { is_expected.to eq([expected_time.to_i, nil]) }
	//    end
	//
	//    context 'when specifying next monday' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 16, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.monday
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next tuesday' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 17, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.tuesday(1)
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying second wednesday from now' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 25, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.wednesday(2)
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next thursday' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 19, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.thursday
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next friday' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 20, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.friday
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next saturday' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 14, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.saturday
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next sunday' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 15, 12, 2, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.sunday
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next monday at 8:00AM' do
	//      let(:expected_time) { Time.zone.local(2018, 7, 16, 8, 0, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.monday.at('8:00')
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//
	//    context 'when specifying next 12th at 6PM' do
	//      let(:expected_time) { Time.zone.local(2018, 8, 12, 18, 0, 0) }
	//
	//      before do
	//        expected_time
	//        with_timestamp.on_day(12).at('18')
	//      end
	//
	//      it { is_expected.to eq [expected_time.to_i, nil] }
	//    end
	//  end
	//end

})
