package common

import "strings"

// WildcardMatch simple wildcard matching, time complexity is O(mn)
func WildcardMatch(str string, ptr string) bool {
	if strings.IndexRune(ptr, '*') == -1 && strings.IndexRune(ptr, '?') == -1 {
		return str == ptr
	}
	m, n := len(str), len(ptr)
	dp := make([][]bool, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]bool, n+1)
	}
	dp[0][0] = true
	for i := 1; i <= n; i++ {
		if ptr[i-1] == '*' {
			dp[0][i] = true
		} else {
			break
		}
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if ptr[j-1] == '*' {
				dp[i][j] = dp[i][j-1] || dp[i-1][j]
			} else if ptr[j-1] == '?' || str[i-1] == ptr[j-1] {
				dp[i][j] = dp[i-1][j-1]
			}
		}
	}
	return dp[m][n]
}
