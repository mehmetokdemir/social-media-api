package friendship

import (
	"github.com/mehmetokdemir/social-media-api/internal/app/common/httpmodel"
	"github.com/mehmetokdemir/social-media-api/internal/app/entity"
)

type ReadFriendship struct {
	Id        uint                        `json:"id"` // ID of the friendship request
	CreatedAt string                      `json:"created_at"`
	Sender    httpmodel.CommonUser        `json:"sender"`
	Receiver  httpmodel.CommonUser        `json:"receiver"`
	Status    entity.FriendshipStatusEnum `json:"status"`
}
