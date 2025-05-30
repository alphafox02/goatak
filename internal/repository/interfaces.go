package repository

import (
	"time"

	"github.com/kdudkov/goutils/callback"

	"github.com/kdudkov/goatak/pkg/model"
	internal "github.com/kdudkov/goatak/pkg/model"
)

type AuthRepository interface {
	Start() error
	Stop()
	CheckAuth(username, password string) bool
	IsValid(username, sn string) bool
	Get(username string) *internal.Device
}

type DeviceRepository interface {
	Start() error
	Stop()
	CheckAuth(username, password string) bool
	IsValid(username, sn string) bool
	Get(username string) *internal.Device
	SaveSignInfo(username, uid, sn string, till time.Time)
	SaveConnectInfo(username, uid, sn string)
}

type ItemsRepository interface {
	Start() error
	Stop()
	ChangeCallback() *callback.Callback[*model.Item]
	DeleteCallback() *callback.Callback[string]
	Store(i *model.Item)
	Get(uid string) *model.Item
	Remove(uid string)
	ForEach(f func(item *model.Item) bool)
	GetCallsign(uid string) string
}

type FeedsRepository interface {
	Start() error
	Stop()
	Store(i *model.Feed2)
	Get(uid string) *model.Feed2
	Remove(uid string)
	ForEach(f func(item *model.Feed2) bool)
}
