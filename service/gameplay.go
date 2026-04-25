package service

import (
	"fmt"
	"sort"

	"numbernama-go/model"
	"numbernama-go/repo"
)

const elementsAddTotal = 10
const rowLength = 9

var defaultElem1to18 = [][]int{
	{1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 1, 1, 2, 1, 3, 1, 4, 1},
	{5, 1, 6, 1, 7, 1, 8},
}

var defaultElem1to19 = [][]int{
	{1, 2, 3, 4, 5, 6, 7, 8, 9},
	{1, 1, 1, 2, 1, 3, 1, 4, 1},
	{5, 1, 6, 1, 7, 1, 8, 1, 9},
}

type GameplayService struct {
	store *repo.MemoryGameplay
}

func NewGameplayService(store *repo.MemoryGameplay) *GameplayService {
	return &GameplayService{store: store}
}

func (s *GameplayService) Initiate(clientID string, req model.InitiateGamePlayRequest) model.InitiateGamePlayResponse {
	state := s.store.State(clientID)
	if req.ResetFlag {
		if state.SelectedKind == "" {
			return model.InitiateGamePlayResponse{ResetFlag: true, Error: "nothing to reset; start a game first"}
		}
		s.start(state, state.SelectedKind)
		return model.InitiateGamePlayResponse{ResetFlag: true, GameplayArray: cloneGrid(state.Grid)}
	}

	kind, err := normalizeGameType(req.GameType)
	if err != nil {
		return model.InitiateGamePlayResponse{ResetFlag: false, Error: err.Error()}
	}
	s.start(state, kind)
	return model.InitiateGamePlayResponse{ResetFlag: false, GameplayArray: cloneGrid(state.Grid)}
}

func (s *GameplayService) Match(clientID string, a, b interface{}) model.MatchResponse {
	state := s.store.State(clientID)
	errs, ok := s.matchState(state, a, b)
	return model.MatchResponse{Err: errs, Matched: ok}
}

func (s *GameplayService) Check(clientID string) model.CheckResponse {
	state := s.store.State(clientID)
	if len(state.Grid) == 0 {
		return model.CheckResponse{}
	}
	var left []int
	for _, row := range state.Grid {
		for _, v := range row {
			if v != 0 {
				left = append(left, v)
			}
		}
	}
	original := append([]int(nil), left...)
	lastRow := len(state.Grid) - 1
	lastRowElemCount := rowLength - len(state.Grid[lastRow])
	if lastRowElemCount > 0 && len(left) > 0 {
		take := min(lastRowElemCount, len(left))
		state.Grid[lastRow] = append(state.Grid[lastRow], left[:take]...)
		left = left[take:]
	}
	for len(left) > 0 {
		take := min(rowLength, len(left))
		state.Grid = append(state.Grid, append([]int(nil), left[:take]...))
		left = left[take:]
	}
	return model.CheckResponse{LastRow: lastRow, LastRowElemCount: lastRowElemCount, NewElems: original}
}

func (s *GameplayService) Clear(clientID string) model.ClearResponse {
	state := s.store.State(clientID)
	rows := make([]int, 0)
	for i := len(state.Grid) - 1; i >= 0; i-- {
		if len(state.Grid[i]) == rowLength && nonZeroCount(state.Grid[i]) == 0 {
			state.Grid = append(state.Grid[:i], state.Grid[i+1:]...)
			rows = append(rows, i)
			state.LastMatched = nil
		}
	}
	sort.Ints(rows)
	return model.ClearResponse{ClearedFlag: len(rows) > 0, ClearedRows: rows}
}

func (s *GameplayService) Undo(clientID string) model.UndoResponse {
	state := s.store.State(clientID)
	if state.LastMatched == nil {
		return model.UndoResponse{OK: false}
	}
	last := *state.LastMatched
	state.Grid[last.Elem1.X][last.Elem1.Y] = last.Val1
	state.Grid[last.Elem2.X][last.Elem2.Y] = last.Val2
	return model.UndoResponse{OK: true, LastMatched: &last}
}

func (s *GameplayService) start(state *model.GameState, kind string) {
	state.SelectedKind = kind
	state.LastMatched = nil
	if kind == "elem1to18" {
		state.Grid = cloneGrid(defaultElem1to18)
		return
	}
	state.Grid = cloneGrid(defaultElem1to19)
}

func normalizeGameType(t interface{}) (string, error) {
	switch v := t.(type) {
	case float64:
		if v == 1 {
			return "elem1to18", nil
		}
		if v == 2 {
			return "elem1to19", nil
		}
	case string:
		if v == "elem1to18" || v == "1" {
			return "elem1to18", nil
		}
		if v == "elem1to19" || v == "2" {
			return "elem1to19", nil
		}
	}
	return "", fmt.Errorf("Select a proper game method.")
}

func (s *GameplayService) matchState(state *model.GameState, pos1, pos2 interface{}) ([]string, bool) {
	e1, e2, err := verifyAndConvert(state.Grid, pos1, pos2)
	if err != nil {
		return []string{err.Error()}, false
	}
	a, b := swapOrder(e1, e2)
	valid := false
	if a.Y == b.Y {
		c := 0
		for i := a.X + 1; i < b.X; i++ {
			if state.Grid[i][a.Y] != 0 {
				break
			}
			c++
		}
		valid = a.X+c+1 == b.X
	} else if a.X == b.X {
		c := 0
		for i := a.Y + 1; i < b.Y; i++ {
			if state.Grid[a.X][i] != 0 {
				break
			}
			c++
		}
		valid = a.Y+c+1 == b.Y
	} else if countNonZeroFrom(state.Grid[a.X], a.Y) == 1 && countNonZeroUntil(state.Grid[b.X], b.Y) == 0 {
		c := 0
		for i := a.X + 1; i < b.X; i++ {
			if nonZeroCount(state.Grid[i]) != 0 {
				break
			}
			c++
		}
		valid = a.X+c+1 == b.X
	}
	if !valid {
		return nil, false
	}

	p1, _ := asPosition(pos1)
	p2, _ := asPosition(pos2)
	state.LastMatched = &model.LastMatch{
		Elem1: p1, Elem2: p2,
		Val1: state.Grid[p1.X][p1.Y],
		Val2: state.Grid[p2.X][p2.Y],
	}
	state.Grid[p1.X][p1.Y] = 0
	state.Grid[p2.X][p2.Y] = 0
	return nil, true
}

func verifyAndConvert(grid [][]int, a, b interface{}) (model.Position, model.Position, error) {
	p1, err := asPosition(a)
	if err != nil {
		return model.Position{}, model.Position{}, fmt.Errorf("Provide inputs in proper format.")
	}
	p2, err := asPosition(b)
	if err != nil {
		return model.Position{}, model.Position{}, fmt.Errorf("Provide inputs in proper format.")
	}
	if p1.X < 0 || p1.Y < 0 || p2.X < 0 || p2.Y < 0 ||
		p1.X >= len(grid) || p2.X >= len(grid) ||
		p1.Y >= len(grid[p1.X]) || p2.Y >= len(grid[p2.X]) {
		return model.Position{}, model.Position{}, fmt.Errorf("Provide inputs in proper format.")
	}
	if grid[p1.X][p1.Y] != grid[p2.X][p2.Y] && grid[p1.X][p1.Y]+grid[p2.X][p2.Y] != elementsAddTotal {
		return model.Position{}, model.Position{}, fmt.Errorf("Inputs are not matching.")
	}
	if grid[p1.X][p1.Y] == 0 || grid[p2.X][p2.Y] == 0 {
		return model.Position{}, model.Position{}, fmt.Errorf("Invalid selections.")
	}
	return p1, p2, nil
}

func asPosition(v interface{}) (model.Position, error) {
	switch t := v.(type) {
	case []interface{}:
		if len(t) != 2 {
			return model.Position{}, fmt.Errorf("bad")
		}
		x, ok1 := toInt(t[0])
		y, ok2 := toInt(t[1])
		if !ok1 || !ok2 {
			return model.Position{}, fmt.Errorf("bad")
		}
		return model.Position{X: x, Y: y}, nil
	case map[string]interface{}:
		x, ok1 := toInt(t["x"])
		y, ok2 := toInt(t["y"])
		if !ok1 || !ok2 {
			return model.Position{}, fmt.Errorf("bad")
		}
		return model.Position{X: x, Y: y}, nil
	default:
		return model.Position{}, fmt.Errorf("bad")
	}
}

func swapOrder(a, b model.Position) (model.Position, model.Position) {
	if b.X < a.X || (b.X == a.X && b.Y < a.Y) {
		return b, a
	}
	return a, b
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

func cloneGrid(src [][]int) [][]int {
	out := make([][]int, len(src))
	for i := range src {
		out[i] = append([]int(nil), src[i]...)
	}
	return out
}

func nonZeroCount(row []int) int {
	n := 0
	for _, v := range row {
		if v != 0 {
			n++
		}
	}
	return n
}

func countNonZeroFrom(row []int, start int) int {
	n := 0
	for i := start; i < len(row); i++ {
		if row[i] != 0 {
			n++
		}
	}
	return n
}

func countNonZeroUntil(row []int, end int) int {
	n := 0
	for i := 0; i < end && i < len(row); i++ {
		if row[i] != 0 {
			n++
		}
	}
	return n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
