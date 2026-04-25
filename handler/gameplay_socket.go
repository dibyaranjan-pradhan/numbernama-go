package handler

import (
	"encoding/json"

	gosocket "github.com/dibyaranjan-pradhan/go-socket"

	"numbernama-go/model"
	"numbernama-go/repo"
	"numbernama-go/service"
	"numbernama-go/utils"
)

type GameplaySocketHandler struct {
	svc    *service.GameplayService
	store  *repo.MemoryGameplay
	logger *utils.Logger
}

func NewGameplaySocketHandler(svc *service.GameplayService, store *repo.MemoryGameplay, logger *utils.Logger) *GameplaySocketHandler {
	return &GameplaySocketHandler{svc: svc, store: store, logger: logger}
}

func (h *GameplaySocketHandler) Register(s *gosocket.Server) {
	s.OnConnect(func(ctx *gosocket.Context) {
		_ = ctx.Emit("go_gameplay_connected", ctx.ClientID())
	})
	s.On("initiateGamePlay", h.onInitiateGamePlay)
	s.On("match", h.onMatch)
	s.On("check", h.onCheck)
	s.On("clear", h.onClear)
	s.On("undo", h.onUndo)
	s.OnDisconnect(func(ctx *gosocket.Context) {
		h.store.Remove(ctx.ClientID())
	})
}

func (h *GameplaySocketHandler) onInitiateGamePlay(ctx *gosocket.Context) {
	var req model.InitiateGamePlayRequest
	if err := ctx.Bind(&req); err != nil {
		_ = ctx.Emit("initiateGamePlay", model.InitiateGamePlayResponse{Error: err.Error()})
		return
	}
	_ = ctx.Emit("initiateGamePlay", h.svc.Initiate(ctx.ClientID(), req))
}

func (h *GameplaySocketHandler) onMatch(ctx *gosocket.Context) {
	var coords []json.RawMessage
	if err := ctx.Bind(&coords); err != nil || len(coords) != 2 {
		h.logger.Warnf("invalid match payload from client=%s", ctx.ClientID())
		_ = ctx.Emit("match", model.MatchResponse{Err: []string{"Provide two positions [[x,y],[x,y]]."}, Matched: false})
		return
	}
	var a, b interface{}
	_ = json.Unmarshal(coords[0], &a)
	_ = json.Unmarshal(coords[1], &b)
	res := h.svc.Match(ctx.ClientID(), a, b)
	res.SelectedElems = positionsForEcho(a, b)
	_ = ctx.Emit("match", res)
}

func (h *GameplaySocketHandler) onCheck(ctx *gosocket.Context) {
	_ = ctx.Emit("check", h.svc.Check(ctx.ClientID()))
}

func (h *GameplaySocketHandler) onClear(ctx *gosocket.Context) {
	cr := h.svc.Clear(ctx.ClientID())
	_ = ctx.Emit("clear", []interface{}{cr.ClearedFlag, cr.ClearedRows})
}

func (h *GameplaySocketHandler) onUndo(ctx *gosocket.Context) {
	u := h.svc.Undo(ctx.ClientID())
	var last interface{}
	if u.LastMatched != nil {
		last = u.LastMatched
	}
	_ = ctx.Emit("undo", []interface{}{u.OK, last})
}

func positionsForEcho(a, b interface{}) [][]int {
	out := make([][]int, 0, 2)
	if p, ok := positionToSlice(a); ok {
		out = append(out, p)
	}
	if p, ok := positionToSlice(b); ok {
		out = append(out, p)
	}
	return out
}

func positionToSlice(v interface{}) ([]int, bool) {
	switch t := v.(type) {
	case []interface{}:
		if len(t) != 2 {
			return nil, false
		}
		x, xok := toInt(t[0])
		y, yok := toInt(t[1])
		if !xok || !yok {
			return nil, false
		}
		return []int{x, y}, true
	case map[string]interface{}:
		x, xok := toInt(t["x"])
		y, yok := toInt(t["y"])
		if !xok || !yok {
			return nil, false
		}
		return []int{x, y}, true
	default:
		return nil, false
	}
}

func toInt(v interface{}) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}
