package friendship

import "github.com/mehmetokdemir/social-media-api/internal/app/entity"

type ReadFriendship struct {
	Id        uint                        `json:"id"` // ID of the friendship request
	CreatedAt string                      `json:"created_at"`
	Sender    ReadFriendshipUser          `json:"sender"`
	Receiver  ReadFriendshipUser          `json:"receiver"`
	Status    entity.FriendshipStatusEnum `json:"status"`
}

type ReadFriendshipUser struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
}
