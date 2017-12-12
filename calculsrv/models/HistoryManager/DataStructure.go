package HistoryManager

import "gopkg.in/mgo.v2/bson"

type DbsLotteryTrend struct {
	Id_      bson.ObjectId `bson:"_id"`
	Game_tag string        `bson:"game_tag"`
	Expect   int `bson:"expect"`
	Position int  `bson:""` // 百:0,十:1,个:2
	Trend    []int // 对应索引表示对应号码,对应数值表示对应遗漏期次
}
