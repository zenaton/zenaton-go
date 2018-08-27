package zenaton

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Wait", func() {

	Describe("WithTimestamp", func() {
		Context("Timezone", func() {
			It("sets the timezone", func() {
				w := NewWait()
				err := w.Timezone("invalidTimezone")
				Expect(err.Error()).To(ContainSubstring("cannot find invalidTimezone"))

				err = w.Timezone("America/Sao_Paulo")
				Expect(err).NotTo(HaveOccurred())
				Expect(w.timezone.String()).To(Equal("America/Sao_Paulo"))
			})
		})
	})

	Context("without timezones", func() {
		Context("when applying duration methods", func() {
			It("should Wait for the specified time", func() {
				w := NewWait().Seconds(1)
				a, b, err := w.GetTimestampOrDuration()
				Expect(err).NotTo(HaveOccurred())
				fmt.Println("a, b", a, b)

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
