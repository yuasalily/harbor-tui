package ui

type FocusArea int

const (
	FocusNav    FocusArea = iota // sidebar
	FocusPage                    // ページ本体
	FocusDialog                  // 確認/ログモーダル
)


