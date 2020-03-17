package runes

// ルーンのスライスsの先頭からcを探し、そのインデックスを返す
//
// @param s --- 検索するスライス
// @param c --- 検索対象のルーン
//
// @return 最初に見つけたポジションのインデックス　見つからない場合は-1を返す
func Index(s []rune, c rune) int {

	for i := range s {
		if s[i] == c {
			return i
		}
	}

	return -1
}
