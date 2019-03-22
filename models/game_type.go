package models

type tGameType struct {
	TbName string
}

var GameType = &tGameType{TbName: "game_types"}

const (
	GameTypeClosed = 1
)
