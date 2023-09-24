package util

import "time"

func ConvertToJapanTime(dbTime time.Time) time.Time {
	japanLocation, _ := time.LoadLocation("Asia/Tokyo")
	return dbTime.In(japanLocation)
}
