package util

import (
	"fmt"
	"time"
)

func FormatDayDate(dateValue uint32) string {
	return fmt.Sprintf("%d", dateValue)
}

func FormatMinuteDate(dateValue uint32) string {
	dayValue := uint16(dateValue & 0xFFFF)
	minuteValue := uint16((dateValue >> 16) & 0xFFFF)

	year := (dayValue / 2048) + 2004
	month := (dayValue % 2048) / 100
	day := (dayValue % 2048) % 100

	hour := minuteValue / 60
	minute := minuteValue % 60

	return fmt.Sprintf("%04d%02d%02d %02d:%02d:00", year, month, day, hour, minute)
}

func ToWindMinuteDate(dateValue uint32) uint32 {
	dayValue := uint16(dateValue & 0xFFFF)
	minuteValue := uint16((dateValue >> 16) & 0xFFFF)

	if minuteValue == 0x30c {
		minuteValue = 0x2b2
	}

	minuteValue--

	return (uint32(minuteValue) << 16) | uint32(dayValue)
}

func GetTodayString() string {
	now := time.Now()
	return fmt.Sprintf("%04d%02d%02d", now.Year(), now.Month(), now.Day())
}

func FormatLongDate(date time.Time) string {
	return fmt.Sprintf("%04d%02d%02d %02d:%02d:%02d", date.Year(), date.Month(), date.Day(),
		date.Hour(), date.Minute(), date.Second())
}

func GetNowString() string {
	now := time.Now()
	return FormatLongDate(now)
}

func GetTimeString() string {
	now := time.Now()
	return fmt.Sprintf("%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
}
