package Dynamic
import (
	"encoding/json"
	"log"

	"bytes"
	"io/ioutil"
	"net/http"
)

type MainD struct {
  config dmConfig
  signal *mediaSVConnection
//  cms_conn *websocket.Conn
  started bool
}

type Observer interface {
	OnInit(channel int)
	OnCreate(channelId string)
	//OnJoin()
	//OnConnect(channelId string)
	OnComplete()
	OnClose()
	OnError(err error)
	//OnMessage(msg1 string, msg2 string)
	// OnStat(...)
}

type msgInitResponse struct {
	channel     int         `json:"channel"`
}
type msgInitRequest struct {
	Credential msgCredential `json:"credential"`
	Env        msgEnv        `json:"env"`
}

type msgCredential struct {
  Account string `json:"account"`
	Password string `json:"password"`
}
type msgEnv struct {
	Os            string `json:"os"`
	OsVersion     string `json:"osVersion"`
	Device        string `json:"device"`
	DeviceVersion string `json:"deviceVersion"`
	NetworkType   string `json:"networkType"`
	SdkVersion    string `json:"sdkVersion"`
}
type Config struct {
  Account string
  Password string
}

func New(config Config, observer Observer) *MainD {
  dm := MainD{}
  dm.config = defaultConfig()

  dm.config.credconfig.account = config.Account
  dm.config.credconfig.password = config.Password
  dm.config.observer = observer
  return &dm
}

func (dm *MainD) StartCast() error {
	dm.started = true
	err := dm.init()
	if err != nil {
		return err
	}
	err = dm.signal.cast_start()
	return err
}

func (dm *MainD) Close() {
	dm.started = false
	dm.signal.close()
}

func (dm *MainD) WriteVideo(data []byte, timestamp uint64, duration float64) {
	if dm.started {
		dm.signal.pc.chanVideo <- commMedia{
			audio:    false,
			data:     &data,
			ts:       timestamp,
			duration: duration,
		}
	}
}

func (dm *MainD) WriteAudio(data []byte, timestamp uint64, duration float64) {
	if dm.started {
		dm.signal.pc.chanAudio <- commMedia{
			audio:    true,
			data:     &data,
			ts:       timestamp,
			duration: duration,
		}
	}
}

func check_to_start(start bool){
  for {
    if start == true {
      //forward
    }
  }
}


func (dm *MainD) init() error {
	//log.Println("init")
	msg := msgInitRequest{
		Credential: msgCredential{
			Account:  dm.config.credconfig.account,
			Password: dm.config.credconfig.password,
		},
		Env: msgEnv{
			SdkVersion: dm.config.sdkconfig.version,
		},
	}
  //bool처리도 하자 (cast_start bool true 오면 forward 보내자)
  /*dm.conn, _, err = websocket.DefaultDialer.Dial(sc.config.appServer, nil)*/
	jsonBytes, err := json.Marshal(&msg)
	if err != nil {
		//log.Printf("[FATAL]")
		return err
	}
	//log.Println("init: send: " + string(jsonBytes))

	url := dm.config.appServer + "/init"
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonBytes))
	if err != nil {
		//log.Printf("[FATAL]")
		return err
	}
	defer resp.Body.Close()

	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("fail to read response data")
		return err
	}
	//log.Println("init: recv: " + string(respByte))

	var initResp msgInitResponse
	err = json.Unmarshal(respByte, &initResp)
	if err != nil {
		//log.Printf("[FATAL]")
		return err
	}
	//if len(initResp.IceServers) > 0 {
	//	rm.config.rtcconfig.iceServers = initResp.IceServers
	//}

	if dm.config.observer != nil {
		dm.config.observer.OnInit(initResp.channel)//여기서 채널 가져오자
	}
  //go check_to_start(dm.config.cast_start)

	dm.signal = newMediaSVConnection(dm.config)
	err = dm.signal.connect()

	return err
}
