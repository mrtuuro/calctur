package main

import (
	"time"

	"calctur/backend/types"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{next: next}
}

func (l *LogMiddleware) ProduceData(data types.Coordinate) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"_id":  primitive.NewObjectID(),
			"lat":  data.Lat,
			"long": data.Lng,
			"took": time.Since(start),
		})
	}(time.Now())
	return l.next.ProduceData(data)
}
