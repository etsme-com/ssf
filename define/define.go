package define

const EtsmeUserID = "etsme"

// //////////////// 产品平台 //////////////////
const (
	//////////////////// 云端 ////////////////////
	PlatformCloud = 0

	//////////////////// 1000 以内为自有产品 ////////////////////
	// MEOS2.0 版本产品
	// 老一代产品定义
	PlatformETSME1  = 1 // 标准版
	PlatformETSME1S = 2 // 希捷 buildroot MEOS2.0 版本
	PlatformMobile  = 3 // marvsmart硬件的演示版本

	// MEOS3.0 版本产品
	// ARM基础版
	PlatformArm     = 100 // 基础标准版
	PlatformArmBase = 100 // 基础标准版
	PlatformArmPro  = 101 // MEOS 3.0 企业版 2合一，家庭版，包括华数基础版

	// X86基础版
	PlatformX86       = 200
	PlatformX86Minipc = 200

	// ARM移动版
	PlatformArmMobile           = 300
	PlatformArmMobileTechvision = 300 // 慧为硬件版本

	// X86企业版
	PlatformX86Pro = 500

	// AI 版
	PlatformAI      = 600
	PlatformAI1684X = 600

	//////////////////// 1000 以外为定制，每个定制产品的固件包必定不一样 ////////////////////
	// 中移动定制版
	PlatformCMCC       = 1000
	PlatformCMCCTCN200 = 1000 // TCN200
	PlatformCMCCTCN100 = 1001 // TCN100

	// 希捷定制版
	PlatformSeagate    = 1100
	PlatformSeagatePro = 1101

	// 广明光电定制版
	PlatformQSI    = 1300
	PlatformQSIPro = 1301

	// 华数定制版
	PlatformWasu    = 1500
	PlatformWasuPro = 1501

	// 华三定制版
	PlatformH3C    = 1600
	PlatformH3CPro = 1601

	// 理想汽车定制版
	PlatformLiAuto    = 1700
	PlatformLiAutoPro = 1701

	// 运维调试版
	PlatformMaintain = 10000
)

// //////////////// 运行场景 //////////////////
var (
	ProductionRunEnv  = 0
	TestingRunEnv     = 1
	Testing2RunEnv    = 2
	DevelopmentRunEnv = 3
)

// //////////////// 升级包 //////////////////
const (
	UpdaterPackagePlatform = "platform"
)

// //////////////// HeartbeatCast 服务心跳Cast routingkey //////////////////
const HeartbeatCast = "HeartbeatCast"

var (
	ServiceInit         = 0
	ServiceHealth       = 1
	ServiceWarning      = 2
	ServiceFault        = 3
	ServicePanicRecover = 4
)

// ServiceHBInfo 服务心跳信息
type ServiceHBInfo struct {
	Name        string
	Version     string
	Sn          string
	Status      int
	Pid         int
	Description string
}

// //////////////// 操作日志 //////////////////
const PHDOperateCast = "PHDOperateCast"

// OperateInfo 操作信息
type OperateInfo struct {
	UniqueID    string // 唯一ID
	User        string // 用户
	Timestamp   string // 时间
	Action      string // 行为
	Description string // 描述
}

// //////////////// EventNotify 系统事件Notify routingkey //////////////////
const PHDEventNotify = "PHDEventNotify"
const PHDEventCast = "PHDEventCast"
const CloudEventCast = "CloudEventCast"
const CloudEventNotify = "CloudEventNotify"

// 事件来源
var (
	PhdEvent   uint32 = 0
	CloudEvent uint32 = 1
	MoonEvent  uint32 = 2
)

// 最高位定义
// // 1: 状态型事件，修复后才能再发，或只能发一次，+1 为修复事件
// // 0: 触发型事件，不可修复
// ///////////////////////////////////////////////////////// 设备端事件 0x000xxxxx ///////////////////////////////////////////////////////////
// 硬件事件 0x0000xxxx
var (
	TestEvt           uint32 = 0x00000000
	PowerInputLossEvt uint32 = 0x80000001
	// ResetRebootEvt    uint32 = 0x80000003
	HardwareInitEvt uint32 = 0x80000005

	CPUTempHighCriticalEvt uint32 = 0x80000011
	CPUTempHighWarningEvt  uint32 = 0x80000013
	CPUTempHighRebootEvt   uint32 = 0x80000015
	CPUPowersaveEvt        uint32 = 0x80000017
	CPUPerformanceEvt      uint32 = 0x80000019

	SSDNotExistEvt             uint32 = 0x80000021
	SSDFsReadOnlyEvt           uint32 = 0x80000023
	SSDLifeEvt                 uint32 = 0x00000025
	SSDUuidChangeEvt           uint32 = 0x80000027
	SSDMountFailEvt            uint32 = 0x80000029
	SSDInitSelfCheckHealthEvt  uint32 = 0x8000002b
	SSDInitSelfCheckWarningEvt uint32 = 0x8000002d
	SSDOverallHealthEvt        uint32 = 0x80000031

	UsbNotExistEvt        uint32 = 0x80000041
	BtHardwareAbnormalEvt uint32 = 0x80000043
	BtServiceAbnormalEvt  uint32 = 0x00000045

	MDNotExistEvt   uint32 = 0x80000051
	MDDegradeEvt    uint32 = 0x80000053
	MDRecoveringEvt uint32 = 0x80000055
	MDMountFailEvt  uint32 = 0x80000057

	PCieNotExistEvt uint32 = 0x80000061

	BatteryChargingEvt uint32 = 0x80000071
)

// 系统事件 0x0001xxxx
var (
	SystemStartEvt           uint32 = 0x80010001
	SystemRebootEvt          uint32 = 0x00010003
	SystemShutdownEvt        uint32 = 0x00010005
	SystemExecRebootCmdEvt   uint32 = 0x80010007
	SystemFactoryCfgResetEvt uint32 = 0x00010009
	SystemDisableSvcCheckEvt uint32 = 0x8001000b

	SystemBindEvt         uint32 = 0x00010011
	SystemGPANingBoundEvt uint32 = 0x00010013
	SystemGPANBoundEvt    uint32 = 0x00010015
	SystemPANBoundEvt     uint32 = 0x00010017
	SystemUnbindEvt       uint32 = 0x00010019
	SystemUnboundEvt      uint32 = 0x0001001b

	SystemUpgradEvt          uint32 = 0x00010031
	SystemUpgradeSuccessEvt  uint32 = 0x00010033
	SystemUpgradeFailEvt     uint32 = 0x00010035
	SystemUpgradeForceEvt    uint32 = 0x00010037
	SystemDownloadSuccessEvt uint32 = 0x00010039
	SystemDownloadFailEvt    uint32 = 0x0001003b

	SystemNetDownEvt uint32 = 0x80010041
	SystemNetUpEvt   uint32 = 0x80010042

	SystemLoggerBreakOutEvt uint32 = 0x00010051
	SystemLoggerEnDebugEvt  uint32 = 0x80010053
	SystemLoggerTooManyEvt  uint32 = 0x80010055

	SystemKernelPanicEvt uint32 = 0x00010061
	SystemKernelDieEvt   uint32 = 0x00010063

	// factory start
	SystemFactoryBTestingEvt uint32 = 0x00010071
	SystemFactoryBSuccessEvt uint32 = 0x00010073
	SystemFactoryBFailEvt    uint32 = 0x00010075

	SystemFactoryDTestingEvt uint32 = 0x00010077
	SystemFactoryDSuccessEvt uint32 = 0x00010079
	SystemFactoryDFailEvt    uint32 = 0x0001007b

	SystemRegisteringEvt     uint32 = 0x00010081
	SystemRegisterSuccessEvt uint32 = 0x00010083
	SystemRegisterFailEvt    uint32 = 0x00010085

	SystemAgeingTestingEvt uint32 = 0x00010091
	SystemAgeingSuccessEvt uint32 = 0x00010093
	SystemAgeingFailEvt    uint32 = 0x00010095
	SystemAgeingNetFailEvt uint32 = 0x00010097
	// factory end

	SystemRunEnvChangeEvt uint32 = 0x000100a9

	SystemPerformanceLoadEvt   uint32 = 0x000100b1
	SystemPerformanceMemoryEvt uint32 = 0x800100b3

	// 数据备份
	SystemMetadataBackuSendEvt  uint32 = 0x000100c1
	SystemMetadataBackupRecvEvt uint32 = 0x000100c3
	SystemDataBackupSendEvt     uint32 = 0x000100c5
	SystemDataBackupRecvEEvt    uint32 = 0x000100c7
	SystemDataBackupFailEEvt    uint32 = 0x800100c9

	// 时间同步失败
	SystemTimeSyncFailEvt uint32 = 0x800100d1

	// 修改设备名
	SystemModifyDevnameEvt uint32 = 0x000100e1

	// 移动版专属
	SystemPasserConnectedEvt uint32 = 0x80011001
	SystemListeningPasserEvt uint32 = 0x80011003
)

// 服务事件 0x0002xxxx
var (
	ServiceExceptionEvt     uint32 = 0x80020001
	ServicePanicRecoverEvt  uint32 = 0x00020003
	ServiceRestServerEvt    uint32 = 0x00020005
	ServiceYamlFailedEvt    uint32 = 0x00020007
	ServiceTooManyOpenFiles uint32 = 0x80020009

	ServiceMemoryCriticalEvt uint32 = 0x00020011
	ServiceMemoryWarningEvt  uint32 = 0x00020013

	// Etsnetd
	ServiceEtsnetdNodeOfflineEvt  uint32 = 0x80021003
	ServiceEtsnetdJoinNetworkEvt  uint32 = 0x00021005
	ServiceEtsnetdLeaveNetworkEvt uint32 = 0x00021007
	ServiceEtsnetdPeerUpdateEvt   uint32 = 0x00021101

	// Estor
	ServiceEstorSSDResidualSpaceEvt        uint32 = 0x80022001
	ServiceEstorFileStarLevelChangeFailEvt uint32 = 0x00022003
	ServiceEstorOsEvt                      uint32 = 0x00022005
	ServiceEstorSqliteEvt                  uint32 = 0x00022007
	ServiceEstorLeveldbEvt                 uint32 = 0x00022009
	ServiceEstorAsyncTaskEvt               uint32 = 0x0002200b
	ServiceEstorPodFullEvt                 uint32 = 0x0002200d
	ServiceEstorNodeIdChangeEvt            uint32 = 0x0002200f
	ServiceEstorSearchCmpEvt               uint32 = 0x00022011
	ServiceEstorThumbUsedRatioEvt          uint32 = 0x00022013

	// Etsecd
	ServiceEtsecdInterceptEvt      uint32 = 0x00023001
	ServiceEtsecdUserCANotExistEvt uint32 = 0x00023003

	// AI
	ServiceSmartalbumModelsNotReadyEvt uint32 = 0x00024001

	// FFMPEG
	ServiceFFMPEGH264EncodeFailEvt uint32 = 0x00025001

	// Updater
	ServiceUpdaterUpgradeEvt          uint32 = 0x00026001
	ServiceUpdaterUpgradeSuccessEvt   uint32 = 0x00026003
	ServiceUpdaterUpgradeFailEvt      uint32 = 0x00026005
	ServiceUpdaterUpgradeForceEvt     uint32 = 0x00026007
	SServiceUpdaterDownloadSuccessEvt uint32 = 0x00026009
	ServiceUpdaterDownloadFailEvt     uint32 = 0x0002600b

	// Ebox
	ServiceEboxAccountCreatedEvt uint32 = 0x00027001

	// Devmgmt
	ServiceDevmgmtMaintainDisableEvt uint32 = 0x00028001

	// Etsfs
	ServiceEtsfsFSRebuildEvt uint32 = 0x00029001

	// tasksched
	ServiceTaskNoAliveJobEvt uint32 = 0x0002a001

	// usermgmt
	ServiceUserContactCreatedEvt uint32 = 0x0002b001
	ServiceUserUpdatedEvt        uint32 = 0x0002b003
	ServiceUserDeletedEvt        uint32 = 0x0002b005
	ServiceUserCreatedEvt        uint32 = 0x0002b007
	ServiceUserContactUpdatedEvt uint32 = 0x0002b009
	ServiceUserContactDeletedEvt uint32 = 0x0002b00b
	ServiceUserProfileCreatedEvt uint32 = 0x0002b00d
	ServiceUserProfileUpdatedEvt uint32 = 0x0002b00f
)

// ///////////////////////////////////////////////////////// 云端事件 0x001xxxxxx ///////////////////////////////////////////////////////////
// 集群和基础组件事件 0x0010xxxx
var (
	// K8S 0x00100xxx
	CloudClusterNodeCpuEvt uint32 = 0x80100001
	CloudClusterNodeMemEvt uint32 = 0x80100003

	// CEPH 0x00101xxx
	CloudClusterCephIOEvt  uint32 = 0x80101001
	CloudClusterCephCapEvt uint32 = 0x80101003
)

// 服务事件 0x0011xxxx
var (
	// Upgrade
	CloudUpgradeFailEvt uint32 = 0x00110001

	// PHD 0x00110xxx
	CloudPhdOfflineEvt  uint32 = 0x00110003
	CloudPhdOnlineEvt   uint32 = 0x00110004
	CloudEboxOfflineEvt uint32 = 0x00110005
	CloudEboxOnlineEvt  uint32 = 0x00110006

	// NC 0x00111xxx
	//CloudNcMoonOfflineEvt    uint32 = 0x80120001
	//CloudNcOfflineEvt        uint32 = 0x80120003
	//CloudNcChannelErrorEvt   uint32 = 0x80120005
	//CloudNcNodeLinkChangeEvt uint32 = 0x00120007
	CloudNcMoonOfflineEvt    uint32 = 0x80111001
	CloudNcOfflineEvt        uint32 = 0x80111003
	CloudNcChannelErrorEvt   uint32 = 0x80111005
	CloudNcNodeLinkChangeEvt uint32 = 0x00111007

	// UC 0x00112xxx
	CloudUCSMSTimeoutEvt uint32 = 0x00112001
	CloudUCSMSFailEvt    uint32 = 0x00112003
)

// ///////////////////////////////////////////////////////// 业务事件 0x002xxxxxx ///////////////////////////////////////////////////////////
// 用户事件 0x0020xxxxx
var (
	UserLoginEvt uint32 = 0x00200001
)

// APP事件 0x0021xxxxx
var (
	AppConnectPhdEvt    uint32 = 0x00210001
	AppDisconnectPhdEvt uint32 = 0x00210003
)

// 事件实体
const (
	// 设备端
	SystemEntity    = "SYSTEM"
	FactoryEntity   = "FACTORY"
	PowerEntity     = "POWER"
	CpuEntity       = "CPU"
	DiskEntity      = "DISK"
	PcieEntity      = "PCIE"
	UsbEntity       = "USB"
	BatteryEntity   = "BATTERY"
	ServiceEntity   = "SERVICE"
	MdEntity        = "MD"
	SSDEntity       = "SSD"
	CAEntity        = "CA"
	FFMPEGEntity    = "FFMPEG"
	AIEntity        = "AI"
	MaintainEntity  = "MAINTAIN"
	DevmonEntity    = "DEVMON"
	EtsecdEntity    = "ETSECD"
	EtsnetdEntity   = "ETSNETD"
	EstorEntity     = "ESTOR"
	EboxEntity      = "EBOX"
	UpgradeEntity   = "UPGRADE"
	TaskschedEntity = "TASKSCHED"
	UsermgmtEntity  = "USERMGMT"

	// 云端
	ClusterEntity  = "CLUSTER"
	MysqlEntity    = "MYSQL"
	RedisEntity    = "REDIS"
	RabbitmqEntity = "RABBITMQ"
	EmqxEntity     = "EMQX"
	CephEntity     = "CEPH"
	SqliteEntity   = "SQLITE"
	LeveldbEntiry  = "LEVELDB"

	CloudmonEntity      = "CLOUDMON"
	UpgrademgmtEntity   = "UPGRADEMGMT"
	PhdmgmtEntity       = "PHDMGMT"
	ClustermonEntity    = "CLUSTERMON"
	ControllermonEntity = "CONTROLLERMON"
	NCEntity            = "NC"
	UCEntity            = "UC"

	// 基础设施
	ECSEntity = "ECS"
	SLBEntity = "SLB"
)

// 事件级别
var (
	InfoLevelEvt     = 0
	MinorLevelEvt    = 1
	MajorLevelEvt    = 2
	CriticalLevelEvt = 3
)

// EventInfo 事件信息
type EventInfo struct {
	UniqueID    string        // 每个事件分配一个唯一ID
	Sn          string        // PHD端填SN，云端填服务名
	Source      uint32        // 0 为设备端，1 为K8S， 2 ECS， 3 SLB
	Service     string        // 来源于哪个服务
	EvtID       uint32        // +1 为恢复事件
	Entity      string        // 实体名，例如 CPU0温度高，实体填CPU，Name 填CPU0
	Name        string        // 事件对象名
	Timestamp   string        // 事件产生时间
	Level       int           // 事件级别
	Sequence    int           // 事件产生序列号
	Description string        // 事件描述
	Parameter   []interface{} // 事件携带参数
}

// //////////////// PHDStateChangeNotify PHD状态改变 routingkey //////////////////
const PHDStateChangeNotify = "PHDStateChangeNotify"
const PHDTransientStateName = "StateMachine:Transient:"

// 100以下为暂态，重启丢失
var (
	StateRebootOrShutdown = 1

	StateHWInit = 2 // 硬件初始化，场景：硬件在线修复，售后维修

	StateUnbanding = 3 // 解绑过程中，屏蔽一切其他故障状态

	StateSSDNotExist   = 10
	StateSSDMountFail  = 11
	StateSSDFsReadOnly = 12
	StateNodeOffline   = 13

	StateFactoryFail   = 20
	StateRegisterFail  = 21
	StateAgeingFail    = 22
	StateAgeingNetFail = 23

	StateMatedataBackupSend = 30
	StateMatedataBackupRecv = 31
	StateDataBackupSend     = 32
	StateDataBackupRecv     = 33

	StateFactoryTesting = 50
	StateRegistering    = 51
	StateAgeingTesting  = 52
	StateUpgrading      = 53
	StateBanding        = 54

	// 移动版专属
	StatePasserConnected = 60 // 有 passer 连接上来
	StateListeningPasser = 61 // 可被 passer 发现并连接

	StateBatteryCharging = 70 // 充电状态，暂时没用
)

// 100以上为稳态，持久化保存
var (
	StateFactory      = 100
	StateRegister     = 101
	StateAgeing       = 102
	StateUnbound      = 110
	StateGPANingBound = 111
	StateGPANBound    = 112
	StatePANBound     = 113 // 断开外网模式
	StateRecycle      = 120

	StateUnknown = 0 // 未知状态
)

// PHDStateChangeInfo PHD状态
type PHDStateChangeInfo struct {
	LastState int
	CurState  int
}

// //////////////// 用户告警 //////////////////
const UserAlarmCast = "UserAlarmCast"

// UserAlarmInfo 用户告警信息
type UserAlarmInfo struct {
	UserID  string     `json:"user_id"`
	Evtinfo *EventInfo `json:"evtinfo"`
}

// //////////////// ETSRPC  //////////////////
// C2DRequest RPC请求 一级头描述符
type C2DRequest struct {
	P     []byte `json:"p"`
	UniId string `json:"uni_id"`
}

// D2CResponse RPC响应 一级头描述符
type D2CResponse struct {
	C     int         `json:"c"`
	P     interface{} `json:"p"`
	UniId string      `json:"uni_id"`
}

type CommonRpcRequestRest struct {
	Method string `json:"method"`
	Url    string `json:"url"`
	Body   []byte `json:"body"`
}

// //////////////// check phd ready  //////////////////
type ElementReadyState struct {
	State       bool
	Description string
}

var (
	AllReadyState   = 0
	SSDReadyState   = 1
	PowerReadyState = 2
	NetReadyState   = 3
	UpgradeState    = 4
	RedisState      = 5
)

type PHDReadyStateInfo struct {
	Strong struct {
		BackupPath   ElementReadyState
		MetadataPath ElementReadyState
		RamblockPath ElementReadyState
	}

	Power struct {
		Powerlost bool
	}

	Net struct {
		Etsnetd ElementReadyState
	}

	Upgrade struct {
		Running bool
	}

	Redis struct {
		State bool
	}
}

// //////////////// PHD Runtime Info //////////////////
var (
	AllRunningState   = 0
	CurCPUTemperature = 1
	CurCPULoad        = 2
)

type PHDRuntimeInfo struct {
	Hardware struct {
		CPU struct {
			Temperature float64
		}
	}

	System struct {
		Load struct {
			Load1  float64
			Load5  float64
			Load15 float64
		}
	}
}

// //////////////// media encode //////////////////
const MediaVideoEncodeRespNotify = "MediaVideoEncodeRespNotify"
const MediaVideoEncodeReqCast = "MediaVideoEncodeReqCast"

type MediaVideoEncodeReq struct {
	Id              int
	SrcWeight       int
	SrcHeight       int
	SrcSize         int64
	DownloadFileUrl string
	UploadPodID     string
	DesWeight       int
	DesHeight       int
	DesFps          int
	DesBps          int
}

var (
	MediaVideoEncodeRespECsuccess    = 0 // 成功
	MediaVideoEncodeRespECnonsupport = 1 // 不符合转码条件
	MediaVideoEncodeRespECffmpegFail = 2 // 转码失败
	MediaVideoEncodeRespECsytemFail  = 3 // 系统环境异常
)

type MediaVideoEncodeResp struct {
	Req *MediaVideoEncodeReq

	UploadFileUrl  string
	UploadFileSize uint64
	Result         int
	Description    string
}

// //////////////// AI //////////////////
const AITaskCast = "AITaskCast"

const HumanPicLabel = "HumanPicLabel"

// ////////////////////////////////// REST 相关 ////////////////////////////////////
type RestCommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Payload []byte `json:"payload"`
}

type RestComResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// //////////////// MQ 选项 //////////////////
type MQOption struct {
	MaxLength int64 //队列深度
}
