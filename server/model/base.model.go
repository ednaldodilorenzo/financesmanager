package model

type IUserDependent interface {
	SetUserID(uint64)
}

type UserDependent struct {
	ID     uint32 `gorm:"autoIncrement;primary_key" json:"id"`
	IdUser uint64 `gorm:"primary_key"`
}

func (ud *UserDependent) SetUserID(userId uint64) {
	ud.IdUser = userId
}
