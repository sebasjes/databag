package databag

import (
  "github.com/gorilla/websocket"
  "encoding/json"
  "sync"
  "time"
)

const bridgeKeepAlive = 6

type BridgeStatus struct {
  status string
}

type Bridge struct {
  bridgeId string
  expires int64
  callerToken string
  calleeToken string
  caller *websocket.Conn
  callee *websocket.Conn
}

type BridgeRelay struct {
  sync sync.Mutex
  bridges []Bridge
}

func (s BridgeRelay) AddBridge(bridgeId string, callerToken string, calleeToken string) {
  s.sync.Lock()
  defer s.sync.Unlock()
  bridge := Bridge{
    bridgeId: bridgeId,
    expires: time.Now().Unix() + (bridgeKeepAlive * 3),
    callerToken: callerToken,
    calleeToken: calleeToken,
  }
  s.bridges = append(s.bridges, bridge)
}

func setStatus(bridge Bridge, status string) {
  msg, _ := json.Marshal(BridgeStatus{ status: status })
  if bridge.caller != nil {
    if err := bridge.caller.WriteMessage(websocket.TextMessage, msg); err != nil {
      LogMsg("failed to notify bridge status");
    }
  }
  if bridge.callee != nil {
    if err := bridge.callee.WriteMessage(websocket.TextMessage, msg); err != nil {
      LogMsg("failed to notify bridge status");
    }
  }
}

func (s BridgeRelay) KeepAlive(bridgeId string) {
  s.sync.Lock()
  defer s.sync.Unlock()
  now := time.Now().Unix()
  var bridges []Bridge
  for _, bridge := range s.bridges {
    if bridge.expires > now {
      bridges = append(bridges, bridge)
      if bridge.bridgeId == bridgeId {
        bridge.expires = now + (bridgeKeepAlive * 3)
        if bridge.caller != nil {
          if err := bridge.caller.WriteMessage(websocket.PingMessage, nil); err != nil {
            LogMsg("failed to ping caller signal");
          }
        }
        if bridge.callee != nil {
          if err := bridge.callee.WriteMessage(websocket.PingMessage, nil); err != nil {
            LogMsg("failed to ping callee signal");
          }
        }
      }
    } else {
      setStatus(bridge, "closed");
    }
  }
  s.bridges = bridges
}

func (s BridgeRelay) RemoveBridge(bridgeId string) {
  s.sync.Lock()
  defer s.sync.Unlock()
  var bridges []Bridge
  for _, bridge := range s.bridges {
    if bridge.bridgeId == bridgeId {
      setStatus(bridge, "closed");
    } else {
      bridges = append(bridges, bridge)
    }
  }
  s.bridges = bridges
}

func (s BridgeRelay) SetConnection(conn *websocket.Conn, token string) {
  s.sync.Lock()
  defer s.sync.Unlock()
  for _, bridge := range s.bridges {
    if bridge.callerToken == token {
      bridge.caller = conn
      if bridge.caller != nil && bridge.callee != nil {
        setStatus(bridge, "connected")
      } else {
        setStatus(bridge, "connecting")
      }
    }
    if bridge.calleeToken == token {
      bridge.callee = conn
      if bridge.caller != nil && bridge.callee != nil {
        setStatus(bridge, "connected")
      } else {
        setStatus(bridge, "connecting")
      }
    }
	}
}

func (s BridgeRelay) ClearConnection(conn *websocket.Conn) {
  s.sync.Lock()
  defer s.sync.Unlock()
  for _, bridge := range s.bridges {
    if bridge.caller == conn {
      bridge.caller = nil
      setStatus(bridge, "connecting")
    }
    if bridge.callee == conn {
      bridge.callee = nil
      setStatus(bridge, "connecting")
    }
	}
}

func (s BridgeRelay) RelayMessage(conn *websocket.Conn, msg []byte) {
  s.sync.Lock()
  defer s.sync.Unlock()
  for _, bridge := range s.bridges {
    if bridge.caller == conn && bridge.callee != nil {
      if err := bridge.callee.WriteMessage(websocket.TextMessage, msg); err != nil {
        LogMsg("failed to relay to callee");
      }
    }
    if bridge.callee == conn && bridge.caller != nil {
      if err := bridge.caller.WriteMessage(websocket.TextMessage, msg); err != nil {
        LogMsg("failed to relay to caller");
      }
    }
  }
}