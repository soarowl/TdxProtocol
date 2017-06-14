package util

import "fmt"

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