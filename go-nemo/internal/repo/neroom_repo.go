package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"netease-kit/nemo/internal/db"
	"netease-kit/nemo/internal/dto"
	"time"
)

type NeRoomRepository struct{}

func NewNeRoomRepository() *NeRoomRepository {
	return &NeRoomRepository{}
}

const (
	KeyNeRoomMemberTable   = "nemo:ent:ne_room_table_key:"
	KeyNeRoomSeatUserTable = "nemo:ent:ne_room_seat_user_table_key:"
)

func (r *NeRoomRepository) AddMember(roomArchiveId string, member dto.RoomMember) error {
	key := KeyNeRoomMemberTable + roomArchiveId
	data, err := json.Marshal(member)
	if err != nil {
		return err
	}
	return db.RDB.HSet(context.Background(), key, member.UserUuid, data).Err()
}

func (r *NeRoomRepository) RemoveMember(roomArchiveId string, userUuid string) error {
	key := KeyNeRoomMemberTable + roomArchiveId
	return db.RDB.HDel(context.Background(), key, userUuid).Err()
}

func (r *NeRoomRepository) ExpireMemberTable(roomArchiveId string, duration time.Duration) error {
	key := KeyNeRoomMemberTable + roomArchiveId
	return db.RDB.Expire(context.Background(), key, duration).Err()
}

func (r *NeRoomRepository) AddSeatUser(roomArchiveId string, index int, user dto.SeatUser) error {
	key := KeyNeRoomSeatUserTable + roomArchiveId
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return db.RDB.HSet(context.Background(), key, fmt.Sprintf("%d", index), data).Err()
}

func (r *NeRoomRepository) RemoveSeatUser(roomArchiveId string, index int) error {
	key := KeyNeRoomSeatUserTable + roomArchiveId
	return db.RDB.HDel(context.Background(), key, fmt.Sprintf("%d", index)).Err()
}
