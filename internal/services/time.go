package services

import "time"

func IsAllowedTime() bool {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 11, 27, 0, 0, currentTime.Location())
	endTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 2, 30, 0, 0, currentTime.Location())

	if startTime.Before(endTime) || startTime.Equal(endTime) {
		return (currentTime.Equal(startTime) || currentTime.After(startTime)) && (currentTime.Equal(endTime) || currentTime.Before(endTime))
	} else {
		return currentTime.After(startTime) || currentTime.Before(endTime) || currentTime.Equal(startTime) || currentTime.Equal(endTime)
	}
}
