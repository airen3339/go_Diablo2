package status

var Config *StatusManage

type StatusManage struct {
	DropItemName                              string //掉落在地面上的物品名字
	ChangeScenceFlg, DoorCountFlg, LoadingFlg bool
	MusicIsPlay                               bool
	OpenBag, OpenMiniPanel                    bool
	IsWalk, IsRun, IsAttack                   bool
	CalculateEnd                              bool
	UIOFFSETX                                 int
	ShadowOffsetX, ShadowOffsetY              int
	PLAYERCENTERX, PLAYERCENTERY              int64
	IsTakeItem, IsTakeJoyStick                bool //是否拿起物品/是否使用摇杆
	Mouseoffset                               int
	CamerOffsetX, CamerOffsetY                float64
	ReadMapSizeWidth, ReadMapSizeHeight       int
	MapTitleX, MapTitleY                      int
	MapZoom                                   int
	CurrentGameScence                         int
	DisPlayDebugInfo                          bool        //是否显示Debug信息
	IsPlayDropAnmi                            bool        //是否播放掉落物品动画
	IsDropDeal                                bool        //是否掉落物品处理中
	DisplaySort                               bool        //人物和物体渲染顺序
	Queue                                     chan []byte //消息队列
	IsNetPlay                                 bool        //是否联网游玩
	IsMobile                                  bool        //是否手机端
}

func NewStatusManage() *StatusManage {
	n := &StatusManage{
		IsWalk:        true,
		UIOFFSETX:     0,
		ShadowOffsetX: -60,
		ShadowOffsetY: -10,
		PLAYERCENTERX: 790 / 2, //LAYOUTX/2
		PLAYERCENTERY: 480 / 2, //LAYOUTY/2
		Mouseoffset:   -1800,
		//玩家初始位置偏移设置
		CamerOffsetX: -5280 + 790/2,
		CamerOffsetY: -1880 + 480/2,
		//读取地图的尺寸
		ReadMapSizeWidth:  0,
		ReadMapSizeHeight: 0,
		//玩家初始逻辑地图坐标
		MapTitleX: 36,
		MapTitleY: 11,
		//
		MapZoom:           8,
		CurrentGameScence: 1,
		Queue:             make(chan []byte),
		//IsNetPlay:         true, //是否联网
		IsMobile: false,
	}
	return n
}

// 初始化
func init() {
	Config = NewStatusManage()
}
