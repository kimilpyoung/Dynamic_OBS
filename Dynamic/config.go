package Dynamic
type dmConfig struct {
	rtcconfig       rtcConfig
	mediaServer     string
	appServer       string
	sdkconfig       sdkConfig
	credconfig      credConfig
	viewconfig      viewConfig
	mediaconfig     mediaConfig
  cast_start chan bool

	channel  int
	observer Observer

	videoCodec string
	audioCodec string
}

type rtcConfig struct {
//	iceServers   []msgIceServer
	simulcast    bool
	sdpSemantics string
}

type sdkConfig struct {
	loglevel string
	contry   string
	version  string
}

type credConfig struct {
  account string
  password string
}

type viewConfig struct {
	local  bool
	remote bool
}

type mediaConfig struct {
	video    bool
	audio    bool
	record   bool
	recvOnly bool
}

func defaultConfig() dmConfig {
	cfg := dmConfig{
		rtcconfig: rtcConfig{
			/*iceServers: []msgIceServer{
				msgIceServer{
					Urls: "stun:stun.l.google.com:19302",
				},
			},*/
			simulcast:    false,
			sdpSemantics: "unified-plan",
		},
		sdkconfig: sdkConfig{
			version: "1.0.0",
		},
    channel:1234,
		mediaServer: "wss://cms-dev.hdmania.com/janus",
		appServer:   "https://cms1.hdmania.com",
    cast_start: make(chan bool, 5),
		mediaconfig: mediaConfig{
			video:  true,
			audio:  true,
			record: false,
		},
	}
	return cfg
}
