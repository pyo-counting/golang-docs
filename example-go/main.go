package main

import (
	"fmt"
	"iter"
	"time"
)

type s1 struct {
	name string
	age  int
}

type i1 int

type ()

func (r i1) Hello() {
	fmt.Println(r)
}

type m1 func(p1, p2 string)

type i2 interface {
	m1(int)
	m2(a int) (b, c int)
}

func f1(ch <-chan string) {
	fmt.Println(ch)
	for v := range ch {
		fmt.Println(v)
	}
}

func MyIterator1() iter.Seq[int] {
	// 이 반환되는 함수가 실제 이터레이터의 본체입니다.
	return func(yield func(int) bool) {
		// 우리는 이 'yield' 함수를 호출해서 값을 전달합니다.
		fmt.Println("10을 보낼게요")
		if !yield(10) {
			return
		} // 10을 for...range 루프에 전달
		fmt.Println("20을 보낼게요")
		if !yield(20) {
			return
		} // 20을 for...range 루프에 전달
		fmt.Println("30을 보낼게요")
		if !yield(30) {
			return
		} // 30을 for...range 루프에 전달
	}
}

func MyIterator2() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			fmt.Printf("💡 이터레이터: %d 생성... yield 호출!\n", i)
			if !yield(i) {
				return
			} // for 루프가 멈추든 말든 신경쓰지 않고 계속 값을 보냄
			fmt.Printf("🔴 이터레이터: yield 호출 후... (이 메시지는 안 보여야 정상)\n", i)
		}
	}
}

func badIterator1() iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 1; i <= 5; i++ {
			fmt.Printf("💡 이터레이터: %d 생성... yield 호출!\n", i)
			yield(i) // for 루프가 멈추든 말든 신경쓰지 않고 계속 값을 보냄
			fmt.Printf("🔴 이터레이터: yield 호출 후... (이 메시지는 안 보여야 정상)\n", i)
		}
	}
}

func main() {
	v1, v2, v3 := 4, "감사", 1.4
	const c1, c2 = 4, 1.4
	var v4 int = 4

	println(v1, v2, v3, v1+int(v3), v2[0])
	println(c1 + c2)
	println(v4)

	for i := 0; i < 4; i++ {
		println("hello")
	}

	if i := true; i {
		println(i)
	}

	switch true {
	case v1 == 4:
		println(v1)
	case v2 == "감사":
		println(v2)
	}

	v5 := s1{
		name: "psy",
		age:  25,
	}

	println(v5.name, v5.age)

	v6 := struct {
		Name string
		Age  int
	}{
		Name: "psy",
		Age:  25,
	}

	println(v6.Name, v6.Age)

	v7 := [5]int{1, 2, 4, 5, 6}
	fmt.Println(v7)
	var v8 [5]int = [5]int{1, 2, 3, 4, 5}
	fmt.Println(v8)
	v9 := [...]int{1, 2, 3, 4, 5}
	fmt.Println(v9)
	var v10 = [...]int{1, 2, 3, 4, 5}
	fmt.Println(v10)
	var v11 [5]int = [...]int{1, 2, 3, 4, 5}
	fmt.Println(v11)
	fmt.Println(len(v11), cap(v11))

	var v12 []int = []int{1, 2, 3, 4, 5}
	fmt.Println(v12, len(v12), cap(v12))
	fmt.Println(append(v12, 6))
	v13 := make([]int, 5, 10)
	fmt.Println(v13, len(v13), cap(v13))
	v14 := make([]int, 5)
	fmt.Println(v14, len(v14), cap(v14))

	var v15 map[string]string = map[string]string{
		"one": "1",
		"two": "2",
	}
	fmt.Println(v15)
	v15["three"] = "3"

	for k, v := range v15 {
		fmt.Println(k, v)
	}

	if elem, ok := v15["four"]; ok {
		fmt.Println(elem)
	}

	v16 := i1(4)
	v16.Hello()

	var v17 interface{} = v16
	switch v17.(type) {
	case string:
		fmt.Println("string")
	case i1:
		fmt.Println("int", v17)
	default:
		fmt.Println("nothing")
	}

	v18 := v17.(i1)
	fmt.Println(v18)

	var v19 chan string
	fmt.Println(v19)
	v19 = make(chan string)
	fmt.Println(v19)

	go func() {
		fmt.Println(<-v19)
	}()

	v19 <- "hello"
	time.Sleep(2 * time.Second)

	v20 := make(chan string, 10)
	for i := range 10 {
		v20 <- fmt.Sprintf("hello %d", i)
	}
	close(v20) // 채널을 닫아야 range 루프가 종료됨
	f1(v20)

	fmt.Println("hello"[1])
	for i := range "hello" {
		fmt.Print("hello"[i])
	}
	fmt.Printf("\n\n")

	for i := range MyIterator1() {
		fmt.Println(i)
		fmt.Println("i을 받았어요")
	}
	fmt.Printf("\n\n")
	fmt.Printf("\n\n")

	for num := range MyIterator2() {
		fmt.Printf("  ➡️ for 루프: %d 받음\n", num)
		if num == 2 {
			fmt.Println("  💥 for 루프: 2에서 break!")
			break
		}
	}
	fmt.Printf("\n\n")
	fmt.Printf("\n\n")

	for num := range badIterator1() {
		fmt.Printf("  ➡️ for 루프: %d 받음\n", num)
		if num == 2 {
			fmt.Println("  💥 for 루프: 2에서 break!")
			break
		}
	}
}
