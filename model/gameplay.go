package model

// GameState is the in-memory gameplay state for one socket client.
type GameState struct {
	Grid         [][]int
	LastMatched  *LastMatch
	SelectedKind string
}

type InitiateGamePlayRequest struct {
	ResetFlag bool        `json:"resetFlag"`
	GameType  interface{} `json:"gameType"`
}

type InitiateGamePlayResponse struct {
	ResetFlag     bool    `json:"resetFlag"`
	GameplayArray [][]int `json:"gameplayArray,omitempty"`
	Error         string  `json:"error,omitempty"`
}

type MatchResponse struct {
	Err           []string `json:"err,omitempty"`
	Matched       bool     `json:"matched"`
	SelectedElems [][]int  `json:"selectedElems,omitempty"`
}

type CheckResponse struct {
	LastRow          int   `json:"lastRow"`
	LastRowElemCount int   `json:"lastRowElemCount"`
	NewElems         []int `json:"newElems"`
}

type ClearResponse struct {
	ClearedFlag bool  `json:"clearedFlag"`
	ClearedRows []int `json:"clearedRows"`
}

type UndoResponse struct {
	OK          bool       `json:"ok"`
	LastMatched *LastMatch `json:"lastMatched,omitempty"`
}

type LastMatch struct {
	Elem1 Position `json:"elem1"`
	Elem2 Position `json:"elem2"`
	Val1  int      `json:"val1"`
	Val2  int      `json:"val2"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
