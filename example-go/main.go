package main

import (
	"fmt"
	"iter" // Go 1.23부터 표준 라이브러리에 포함
)

// Numbers 함수가 iter.Seq 타입의 함수를 반환합니다.
func Numbers(max int) iter.Seq[int] {
	// 이 반환되는 익명 함수가 iter.Seq의 실제 구현체입니다.
	// 이 함수는 'yield func(int) bool' 이라는 매개변수를 받습니다.
	return func(yield func(int) bool) { // <-- 여기서 'yield'는 콜백 매개변수입니다.
		for i := 0; i < max; i++ {
			fmt.Println("hello", i)
			// 개발자는 이 'yield' 콜백 매개변수를 호출하여 'i' 값을 '내보냅니다'.
			if !yield(i) { // yield가 false를 반환하면 소비자가 더 이상 값을 원하지 않으므로 중단합니다.
				return
			}
		}
	}
}

func main() {
	for num := range Numbers(5) { // <--- 여기가 소비하는 쪽입니다.
		fmt.Println("bye", num) // 1
		if num == 3 {
			break
		}
	}
}
