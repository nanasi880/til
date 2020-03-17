package runes

import "unicode"

// スライスの前後から空白文字を除去する
// この関数はスライスをコピーせず、サブスライスを返す
//
// @param s --- スライス
//
// @return 空白文字を除去したスライス
func TrimSpace(s []rune) []rune {
	return TrimFunc(s, unicode.IsSpace)
}

// スライスの前後から、条件に合致するルーンを除去する
// この関数はスライスをコピーせず、サブスライスを返す
//
// @param s --- スライス
// @param f --- 条件関数
//
// @return 条件に合致する文字を除去したスライス
func TrimFunc(s []rune, f func(rune) bool) []rune {
	return TrimRightFunc(TrimLeftFunc(s, f), f)
}

// スライスの前方から、条件に合致するルーンを除去する
// この関数はスライスをコピーせず、サブスライスを返す
//
// @param s --- スライス
// @param f --- 条件関数
//
// @return 条件に合致する文字を除去したスライス
func TrimLeftFunc(s []rune, f func(rune) bool) []rune {

	index := indexFunc(s, f, false)
	if index == -1 {
		return nil
	}

	return s[index:]
}

// スライスの後方から、条件に合致するルーンを除去する
// この関数はスライスをコピーせず、サブスライスを返す
//
// @param s --- スライス
// @param f --- 条件関数
//
// @return 条件に合致する文字を除去したスライス
func TrimRightFunc(s []rune, f func(rune) bool) []rune {

	index := lastIndexFunc(s, f, false)
	if index == -1 {
		return nil
	}

	return s[:index+1]
}

// スライスを前方からイテレートし、条件式が期待する値を返した最初のインデックスを返す
// 最後まで条件式が期待する値を返さない場合は-1を返す
//
// @param s     --- スライス
// @param f     --- 条件関数
// @param truth --- 条件関数が返す値の期待値
//
// @return インデックス
func indexFunc(s []rune, f func(rune) bool, truth bool) int {

	for i, c := range s {
		if f(c) == truth {
			return i
		}
	}

	return -1
}

// スライスを後方からイテレートし、条件式が期待する値を返した最初のインデックスを返す
// 最後まで条件式が期待する値を返さない場合は-1を返す
//
// @param s     --- スライス
// @param f     --- 条件関数
// @param truth --- 条件関数が返す値の期待値
//
// @return インデックス
func lastIndexFunc(s []rune, f func(rune) bool, truth bool) int {

	for i := len(s) - 1; i >= 0; i-- {
		if f(s[i]) == truth {
			return i
		}
	}

	return -1
}
