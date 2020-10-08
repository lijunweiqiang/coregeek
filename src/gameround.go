package main

const (
	//收集命令
	CmprssTake = "COMPRESS_TAKE"
	Take       = "TAKE"
	//购买命令
	BuyLooter = "BUY_LOOTER"
	BuyPower  = "BUY_POWER"
)

type GameRound struct {
	round          int32
	turn           string
	teams          TeamsInfo
	gameMap        GameMap
	resourceZones  []ResourceZone
	collectHistory ResourceUnit
}

type TeamsInfo struct {
	yourTeam  OneTeamInfo
	enemyTeam OneTeamInfo
}

type GameMap struct {
	width       int       //地图宽度
	height      int       //地图长度
	mapEntities []MapInfo //地图详细信息
}

type ResourceZone struct {
	gemType       string         //资源类型
	pos           Pos            //资源位置
	resourceUnits []ResourceUnit //资源单位列表
}

type ResourceUnit struct {
	gemType               string //资源类型
	pos                   Pos    //资源位置
	index                 string //资源唯一索引
	sizeOfUnit            int32  //单位资源大小。当前设置为1024
	sizeOfCompressionUnit int32  //该资源单位被压缩装载后的大小
}

type Pos struct {
	x int32
	y int32
}

type OneTeamInfo struct {
	points  int32  //当前赚取的积分
	golds   int32  //当前的总金币数
	campPos Pos    //己方基地位置信息
	works   Works  //资源机器人信息
	looter  Looter //能量掠夺者
}

type MapInfo struct {
	pos  Pos
	data string
}

type Works struct {
	playerName string         //资源机器人名称
	power      int32          //当前资源机器人能量
	pos        Pos            //当前资源机器人所在位置
	load       int32          //资源机器人已装载资源大小，单位KB
	loadInfo   []ResourceUnit //资源机器人当前装载资源详情
	maxLoad    int32          //资源机器人可装载资源最大值，单位KB
}

type Looter struct {
	playerName string //能量掠夺者名称
	pos        Pos    //当前掠夺者所在位置
}

type GameRoundResp struct {
	path             []Pos
	round            int32
	turn             string
	collectCommand   []CollectCommand
	purchaseCommands []PurchaseCommand
}

type CollectCommand struct {
	zone        Pos
	index       string
	commandType string
}

type PurchaseCommand struct {
	pos          Pos
	purchaseType string
}
