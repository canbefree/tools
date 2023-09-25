package helper

import "time"

func GetLastWeekTime(now time.Time) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	if now.Unix() == 0 {
		now = time.Now()
	}
	offset := int(now.Weekday()) // 今天是周几
	if offset == 0 {
		offset = 6
	} else {
		offset = offset - 1
	}

	weekEndTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc).AddDate(0, 0, 6-offset)
	return weekEndTime.Unix()
}

func GetLastMonthTime(now time.Time) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	if now.Unix() == 0 {
		now = time.Now()
	}
	year := now.Year()
	month := now.Month() + 1
	if month == 13 {
		month = 1
		year = year + 1
	}
	return time.Date(year, month, 1, 0, 0, 0, 0, loc).Unix() - 1
}
