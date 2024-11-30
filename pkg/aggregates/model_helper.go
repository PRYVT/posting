package aggregates

import "github.com/PRYVT/posting/pkg/models/query"

func GetPostModelFromAggregate(userAggregate *PostAggregate) *query.Post {
	return &query.Post{
		Id:          userAggregate.AggregateId,
		Text:        userAggregate.Text,
		ImageBase64: userAggregate.ImageBase64,
		ChangeDate:  userAggregate.ChangeDate,
		UserId:      userAggregate.UserId,
	}
}
