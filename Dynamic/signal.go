package Dynamic

import (
	"log"
	"time"

	janus "github.com/notedit/janus-go"
	//"github.com/pion/webrtc/v2"
)
type mediaSVConnection struct {
	comm   chan commChan
	config dmConfig
	gateway *janus.Gateway
  session *janus.Session
  handle *janus.Handle
	pc     *peerconnection
	closed bool
}

func watchHandle(handle *janus.Handle) {
	// wait for event
	for {
		msg := <-handle.Events
		switch msg := msg.(type) {
		case *janus.SlowLinkMsg:
			log.Println("SlowLinkMsg type ", handle.Id)
		case *janus.MediaMsg:
			log.Println("MediaEvent type", msg.Type, " receiving ", msg.Receiving)
		case *janus.WebRTCUpMsg:
			log.Println("WebRTCUp type ", handle.Id)
		case *janus.HangupMsg:
			log.Println("HangupEvent type ", handle.Id)
		case *janus.EventMsg:
			log.Printf("EventMsg %+v", msg.Plugindata.Data)
		}
	}
}

func newMediaSVConnection(config dmConfig) *mediaSVConnection {
	return &mediaSVConnection{
		comm:   make(chan commChan, 5),
		config: config,
	}
}

func (msc *mediaSVConnection) connect() error{
	// Prepare the configuration

  msc.pc = newPeerConnection(msc.config, msc.comm)
	gateway, err := janus.Connect(msc.config.mediaServer)
	if err != nil {
		return err
	}
  msc.gateway = gateway
	session, err := gateway.Create()
	if err != nil {
		return err
	}
  msc.session = session
	handle, err := session.Attach("janus.plugin.videoroom")
	if err != nil {
		return err
	}
  msc.handle = handle
	go func() {
		for {
			if _, keepAliveErr := session.KeepAlive(); err != nil {
    		 panic(keepAliveErr)
			}

			time.Sleep(25 * time.Second)
		}
	}()

	go watchHandle(handle)

	msg, err := msc.handle.Message(map[string]interface{}{
		"request": "join",
		"ptype":   "publisher",
		"room":    msc.config.channel,
		"id":      1825,
	}, nil)

	if err != nil {
		return err
	}
	if msg.Plugindata.Data != nil {

	}

	select {}
}

func (msc *mediaSVConnection) cast_start() error {
    offer, err := msc.pc.peerConnection.CreateOffer(nil)
    if err != nil {
      panic(err)
    }
    err = msc.pc.peerConnection.SetLocalDescription(offer)
    if err != nil {
      panic(err)
    }
  	msg, err := msc.handle.Message(map[string]interface{}{
  		"request": "publish",
  		"audio":   true,
  		"video":   true,
  		"data":    false,
  	}, map[string]interface{}{
  		"type":    "offer",
  		"sdp":     offer.SDP,
  		"trickle": false,
  	})
  	if err != nil {
  		return err
  	}

  	if msg.Jsep != nil {
      msc.pc.comm <- commChan{
        cmd: "onSdp",
        sdp_type: "answer",
        sdp_des:msg.Jsep["sdp"].(string),
      }
  		// Start pushing buffers on these tracks
			//여기서 포워드 시작을 하면 됨.
      //msg.config.cast_start <- true
  	}
		return nil
}


func (msc *mediaSVConnection) close() error {
	msc.pc.comm <- commChan{
		cmd: "mediastop",
	}
  msc.gateway.Close()
	msc.quit()
	return nil
}

func (msc *mediaSVConnection) quit() {
	// make pcReadLoop stop
	msc.comm <- commChan{
		cmd: "quit",
	}
	msc.closed = true
	if msc.config.observer != nil {
		msc.config.observer.OnClose()
	}
  close(msc.config.cast_start)
	close(msc.pc.comm)
}
