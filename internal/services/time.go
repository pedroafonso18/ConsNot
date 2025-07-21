package services

import "time"

func IsAllowedTime() bool {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 7, 30, 0, 0, currentTime.Location())
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 2, 30, 0, 0, currentTime.Location())
	weekday := currentTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	if startTime.Before(endTime) || startTime.Equal(endTime) {
		return (currentTime.Equal(startTime) || currentTime.After(startTime)) && (currentTime.Equal(endTime) || currentTime.Before(endTime))
	} else {
		return currentTime.After(startTime) || currentTime.Before(endTime) || currentTime.Equal(startTime) || currentTime.Equal(endTime)
	}
}
