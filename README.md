### about docs
- go1.24.0 버전 기준으로 작성

- 특정 타입의 인스턴스(값)를 생성하는 것을 주된 목적으로 하는 함수를 golang에서는 관용적으로 생성자 함수라고 부른다. 보통 `New<타입이름>` 형태의 이름을 갖는다. golang 공식 문서 도구(pkg.go.dev의 기반이 되는 godoc)는 이러한 '생성자' 패턴을 인지해서 생성자 함수를 해당 타입의 하위 목차에 배치한다.

### [A Tour of Go](https://go.dev/tour/list)
- golang은 package로 구성되며 golang으로 개발된 프로그램은 main package을 통해 실행된다.
- `import` 키워드를 사용해 package가 위치한 경로를 명시함으로써 package를 import할 수 있다. module 이름(경로)과 하위 디렉토리 경로로 구성될 수 있으며 보통 편의성을 위해 모듈의 이름 중 마지막 경로 위치를 package 이름과 동일하게 짓는다. 이는 단일 module 내 여러 package를 구성하는 경우에도 동일하게 적용(디렉토리 이름을 package 이름으로 짓는다)된다.
    - go 프로그램은 package로 구성된다. package는 동일 디렉토리에 있는 소스 파일의 집합으로 같이 컴파일된다. repository는 보통 1개 이상의 module을 포함한다. module은 package의 집합으로 같이 릴리즈된다. 일반적으로 repository는 루트 디렉토리 1개의 module만 포함(`go.mod` 파일을 통해 module의 경로를 명시)한다.
- `import` 키워드를 여러 번 사용해 여러 package를 import할 수도 있지만 `import (...)`와 같이 사용하는 것을 권장한다.
- package 내에서 대문자로 시작되는 이름을 갖는 경우 해당 package 밖에서도 참조가 가능하며 이를 exported name이라고 한다. 반대로 소문자로 시작되는 이름을 갖는 경우 package 내부에서만 참조가 가능하다. 내장 타입은 대문자로 시작하지 않아도 접근할 수 있다.
- 함수의 반환 값에 이름을 지정하는 경우 함수의 최상단에서 정의된 변수로 취급된다. 함수에서는 `return` 키워드만 사용해도 반환이된다. `return` 키워드를 생략하는 것은 불가능하다. 이를 naked return이라고 부르며 짧은 길이의 함수에서만 사용하는 것을 권장한다.
- `var` 키워드를 사용해 변수를 선언할 수 있다. 변수 선언 시 초기화를 수행하면 변수의 타입을 생략할 수 있다.
- 함수 내부에서는 `:=` short assignment statement만 이용해 변수를 선언 및 초기화할 수 있다. 함수 외부에서는 항상 키워드로 시작해야 하기 때문에 사용할 수 없다.
- 변수 선언 시 초기화를 하지 않으면 zero value가 할당된다. 숫자 타입일 경우 0, boolean 타입일 경우 false, string 타입일 경우 ""
- golang은 묵시적 형변환이 불가능하며 항상 `T(v)` 표현식을 사용해 타입 변환을 수행해야 한다.
- 변수 선언 시 타입을 지정하지 않는 경우 변수의 타입은 초기화 값으로부터 추론된다. 초기화 값이 명시적인 타입을 갖는 경우 변수는 동일한 타입을 갖는다. 하지만 타입이 없는 숫자 상수의 경우 정밀도에 따라 `int`, `float64`, `complex128` 타입으로 추른된다. 타입이 지정되지 않은 상수는 사용되는 문맥에 따라 필요한 타입을 갖게된다.
- `const` 키워드를 사용해 상수를 선언할 수 있다. 상수는 문자, 문자열, boolean, 숫자 타입일 수 있다(`:=` short assignment statement는 변수 선언시 사용되기 때문에 상수에 대해서는 사용할 수 없다).
- 반복문을 위해 `for` 키워드만 지원한다. ()로 감싸지 않아도 되지만 블럭에는 {}가 항상 필요하다. `for ... range` 문에서 slice, 배열, map, int, string을 사용해 반복할 수 있다.
- `if`문도 `for`문과 동일하게 ()로 감싸지 않아도 되지만 블럭에는 {}가 항상 필요하다. `if`문에 short assignment statement를 사용할 수 있으며 해당 변수는 `else if`, `else` 블럭에서도 사용할 수 있다.
- `switch 테스트 표현식 {...}`문에는 `break`문을 명시할 필요없다. 테스트 표현식을 생략하는 경우 `switch true {...}`와 동일하다.
- `defer`문은 함수의 반환까지 함수의 실행을 연기한다. 함수의 매개변수는 `defer`문에서 결정되지만 호출은 지연된다. 지연된 함수는 stack에 push되며 LIFO로 실행된다.
- `*T`는 포인터 타입으로 zero value은 `nil`이다. `&` 연산자는 피연산자의 포인터를 생성한다. `*`연산자는 포인터가 가리키는 값을 나타낸다. 포인터에 대한 산술 연산자는 없다.
- struct는 필드의 집합으로 각 필드는 .을 이용해 접근한다. struct pointer의 경우 `(*p).X`을 통해 접근 가능하지만 번거롭기 때문에 `p.X`와 같이 접근하는 것을 허용한다. `type 이름 struct {필드목록}`문을 이용해 struct 타입을 정의할 수 있다.
- `type` 키워드를 여러번 사용하는 대신 `type (...)`와 같이 사용할 수 있다.
- struct 리터럴은 순서대로 값을 나열하거나 `Name:`처럼 필드의 이름과 값을 순서 상관없이 나열해 표현할 수 있다. 새로운 필드 추가, 순서 변경에 대한 내구성을 위해 필드를 명시하는 것을 권장한다.
- `[n]T` 타입은 배열이다. 배열의 크기는 배열 타입의 일부이기 때문에 크기를 조정할 수 없다.
- `[]T` 타입은 slice로 동적으로 크기 조절이 가능하며 배열보다 더 일반적으로 사용한다. slice는 내부적으로 배열을 가리키는 포인터, 길이, 용량 정보를 저장하며 zero value은 `nil`이다. 배열 또는 slice 변수 a에 대해 `a[low:high]`의 표현식을 사용해 슬라이싱해 slice 값을 얻을 수 있다.
- slice의 길이는 slice를 통해 접근할 수 있는 요소의 길이를 나타내며 용량은 실제 slice가 참조가 있는 배열 길이를 나타낸다(slice가 가리키는 첫 인덱스부터 배열의 마지막 인덱스까지).
- `func cap(v Type) int` 내장 함수는 배열, slice, 배열을 가리키는 포인터, buffered channel의 크기를 반환한다.
- `func len(v Type) int` 내장 함수는 배열, slice, 배열을 가리키는 포인터, map, 문자열, buffered channel(아직 읽지 않은 메시지)의 길이를 반환한다.
- `func make(t Type, size ...IntegerType) Type` 내장 함수는 slice, map, channel 타입을 생성하는 데 사용할 수 있다. slice일 경우 두 번째, 세 번째 인자는 각각 length, capacity를 나타낸다. map의 경우 첫 번째 인자만 필요하다. channel의 경우 두 번째 인자는 buffered channel을 생성할 때 사용된다.
- `func append(slice []Type, elems ...Type) []Type` 내장 함수를 사용해 slice 마지막 인덱스 뒤에 값을 계속 추가할 수 있다. 기존 배열의 크기가 작으면 크기가 더 큰 배열을 생성 및 할당한다.
- `for i, v := range sli {}`와 같은 표현식을 사용해 slice에 대해 반복문을 사용할 수 있다. 변수 i는 인덱스, v는 해당 인덱스의 복사된 값을 갖는다. `_` 표현식을 사용해 할당을 하지 않을 수 있으며 v는 생략할 수도 있다.
- `map[K]V` 타입은 맵이다. zero value은 nil이며 리터럴은 struct와 다르게 key를 생략할 수 없다. `m[key]` 표현식을 사용해 map에 저장된 요소에 접근할 수 있다. `elem = m[key]` 표현식은 변수 m에 key가 없을 경우 오류가 발생한다. 반면 `elem, ok = m[key]` 표현식은 변수 m이 key가 있으면 ok 변수에 true, 없다면 false 값을 갖는다. key가 없는 경우 elem 변수에 zero value가 할당된다.
- `func delete(m map[Type]Type1, key Type)` 내장 함수를 사용해 맵 변수 m에서 key를 삭제할 수 있다.
- golang은 class 개념이 없지만 type에 method를 정의할 수 있다. method는 receiver라는 인자를 받는 함수다. 함수와 다른 점은 단순히 receiver라는 특별한 인자가 있다는 것이며 기능적으로는 함수와 동일하다. method는 type이 정의된 package 내에서만 선언할 수 있다. 그렇기 때문에 int와 같은 내장 타입에는 method를 사용자가 정의할 수 없다. 물론 `type` 키워드를 사용해 int를 다시 한번 정의한 후 method를 정의할 수 있다. 아래는 int 타입을 다시 정의하고 사용하는 예시다.
    ``` go
    type MyString string

    func (s *MyString) String() string {
    	return string(*s)
    }

    func main() {
    	str := MyString("hello")
    	fmt.Println(str)
    }
    ```
- method 선언 시 pointer receiver를 사용할 수 있다. 이는 함수가 포인터 타입의 인자를 통해 실제 값을 변경하는 것과 같이 동일하게 동작한다. `*T` 타입에 대한 method를 `T` 타입의 변수에서 접근할 때 `(&T).X`와 같이 접근할 수 있지만 번거롭기 때문에 `T.X`와 같이 접근하는 것을 허용한다. 반대인 경우에도 `(*T).X` 대신 `T.X`와 같이 접근 가능하다.
- method를 통해 실제 값을 변경이 필요하거나, 크기가 큰 경우 method 호출 시 항상 값을 복사하는 불필요성을 피하기 위해 pointer receiver를 사용한다. 한 타입에 method를 선언할 때 pointer receiver, value receiver를 섞어 쓰는 것을 권장하지 않는다.
- interface는 struct와 유사하지만 필드의 집합이 아니라 method의 집합을 나타낸다. interface 타입은 method를 모두 구현한 타입을 값으로 가질 수 있다. `type 이름 interface {method 목록}` 표현식을 사용해 선언할 수 있다.
- interface를 구현한다는 것은 어떤 키워드를 통해 명시적으로 수행하는 것은 아니며 단순히 interface 타입에 포함된 모든 method를 선언해 암묵적으로 구현한다. interface의 선언식과 구현이 동일 package일 필요는 없다. interface는 `(value, type)` 값을 갖는다고 생각할 수 있다.
- interface의 zero value은 nil이며 method 호출 시 runtime error가 발생한다.
- method를 명시하지 않은 interface를 empty interfac라고 한다. empty interface는 모든 value, 모든 type을 가질 수 있다. empty interface는 따로 `type` 키워드를 통해 선언할 필요가 없으며 변수 선언 시 `var i interface{}`와 같이 사용할 수 있다. empty interface는 알려지지 않은 타입의 값을 다룰 때 사용된다.
- type assertion은 interface가 nil이 아니며 타입 T임을 확인하는 것을 말한다. `i.(T)` 표현식은 interface i가 T 타입임을 나타내며 i가 가리키는 T 타입의 변수를 반환한다. 만약 nil이나 T 타입이 아닐 경우 오류가 발생한다. `t, ok := i.(T)` 표현식은 두 번째 반환 값을 통해 타입 T가 맞는지에 따른 boolean 값을 반환한다.
- type switch문은 switch의 테스트 표현식에 `i.(type)`을 사용해 interface 변수 i에 대한 내장 타입에 대한 case문을 작성할 수 있다.
- 가장 흔한 interface는 fmt package에 있는 Stringer interface다. fmt package는 값을 출력하기 위해 Stringer interface의 String method를 호출한다.
    ``` go
    type Stringer interface {
        String() string
    }
    ```
- 추가적으로 에러를 나타내는 내장 error interface가 있다. 내장 타입이기 때문에 소문자로 시작할 수 있다.
    ``` go
    type error interface {
        Error() string
    }
    ```
- golang은 generic 함수, generic 타입을 통해 generics를 제공한다.
    ``` go
    func Index[T comparable](s []T, x T) int {}
    ```
    뿐만 아니라 `type` 키워드에 type parameter를 추가함으로써 generic 타입을 정의할 수 있다. 아래는 예시다.
    ``` go
    type List[T any] struct {
	    next *List[T]
	    val  T
    }
    ```
- `func close(c chan<- Type)` 함수를 사용해 사용자는 channel을 닫을 수 있다. 해당 함수는 더 이상 보낼 메시지가 없을 경우 송신 측에서 닫을 수 있도록 사용하는 것을 권장한다. 만약 수신 측에서 닫으면 이를 모르는 송신측에서 메시지를 보내면 panic에 빠질 수 있다. channel은 파일처럼 닫을 필요는 없다. 수신측에서는 `v, ok := <-ch` 문을 사용해 channel이 닫혔는지 여부를 확인할 수 있다. 채널이 닫힌 경우 첫 번째 반환 값은 channel의 underlying type zero value다.
- `for i := range c {...}` 문을 사용해 channel이 닫힐 때까지 반복문을 실행할 수 있다. buffered channel의 경우 channel이 닫힌 후에도 buffer에 남은 메시지 개수만큼 for 문이 실행된다.
    ``` go
    func main() {

        queue := make(chan string, 2)
        queue <- "one"
        queue <- "two"
        close(queue)

        for elem := range queue {
            fmt.Println(elem)
        }
    }
    ```

### [GO 프로그래밍 입문](https://codingnuri.com/golang-book/)
- 문자열은 "(double quote), `(backtick)을 사용해 표현할 수 있다. double quote로 표현하는 문자열은 줄바꿈을 포함할 수 없으며 이스케이프 문자열 사용할 수 있다. backtick으로 표현하는 문자열은 줄바꿈을 포함할 수 있으며 이스케이프 문자열을 지원하지 않는다.
- 문자열은 byte slice(`[]byte`) 표현되기 때문에 `"hello world"[4]`와 같이 인덱스를 통해 접근 가능하다. 문자열에 `len(s)` 함수를 사용해 길이를 확인할 수 있다.
- 함수 정의 시 가변 인자를 사용할 수 있다. 이 때 해당 가변 인자는 함수 내에서 slice로 접근할 수 있다. 함수 호출 시에는 가변 인자에 개별 인자, 배열, slice를 사용할 수 있다. slice를 매개변수로 사용 시 `fmt.Println(sli...)`와 같이 호출해야 한다.
    ``` go
    // fmt.Println() 함수 예시
    func Println(a ...interface{}) (n int, err error){...}
    ```
- `func panic(v any)` 내장 함수는 현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴한다(런타임 오류). 이러한 panic 모드 실행 방식은 다시 상위 함수에도 똑같이 적용되고, 계속 콜스택을 타고 올라가며 적용된다. 그리고 마지막에는 프로그램이 에러를 내고 종료하게 된다. `func recover() any` 내장 함수는 panic() 함수에 의한 패닉 상태를 중단하고 panic() 함수 호출 시 전달했던 인자를 반환한다. panic() 함수 호출 시 런타임 에러가 발생해 즉시 호출 중이던 함수가 종료되기 때문에 recover() 함수를 defer문과 사용해야 한다.
- `func new(Type) *Type` 내장 함수는 매개변수 타입에 대해 메모리를 할당(zero value를 할당)하고 포인터를 반환한다.
- struct에 명시적인 필드 이름 없이 타입만 명시해 embedded field를 사용할 수 있다. embedded field는 타입을 사용해 접근할 수 있다. 만약 struct x의 embedded filed y가 필드 또는 method z를 갖고 있을 때, x.y.z로 접근할 수 있다. 뿐만 아니라 x의 필드인 것처럼 x.z와 같이 접근할 수 있다. 이를 promote 됐다고 말한다. 뿐만 아니라 embedded field의 method도 동일하게 자신의 것인처럼 접근할 수 있으며, 관련 interface를 구현한 것처럼 간주된다.
    ``` go
    // example 1
    // A struct with four embedded fields of types T1, *T2, P.T3 and *P.T4
    type test struct {
    	T1        // field name is T1
    	*T2       // field name is T2
    	P.T3      // field name is T3
    	*P.T4     // field name is T4
    	x, y int  // field names are x and y
    }

    // example 2
    type ClassInfo struct {
      Class int
      No int
    }

    type Student struct {
      ClassInfo
      Name string
    }

    var s1 DupStudent = DupStudent{
      ClassInfo: ClassInfo{Class: 1, No: 1},
      Name:      "John",
      No:        10,
    }

    fmt.Println(s1.Class, s1.No, s1.ClassInfo.No) // 1 10 1
    ```
- 하나 이상의 작업을 동시에 진행하는 것을 동시성(concurrency)라 한다. golang에서는 goroutine, channel을 통해 동시성을 지원한다.
- `go` 키워드를 사용해 goroutine을 생성할 수 있다. `go` 키워드 다음 함수 호출 표현식을 사용하면 된다. main 함수도 goroutine에서 실행되며 main 함수가 종료되면 프로그램의 종료로 이어지기 때문에 다른 goroutine이 모두 종료된 후 main 함수의 goroutine을 종료하도록 해야 한다.
- `chan` 키워드를 사용해 channel을 생성할 수 있다. `chan` 키워드 다음 channel의 타입을 지정할 수 있다. `<-` 연산자를 사용해 channel에 메시지를 전달하거나 channel로부터 메시지를 전달받을 수 있다. 기본적으로 channel은 송신과 수신이 완료되기 전까지 blocking 된다. 이를 통해 channel은 두 goroutine이 서로 통신하고 실행 흐름을 동기화할 수 있다. 아래는 string 타입의 channel 변수를 생성하고 메시지를 송수신하는 예시다.
    ``` go
    func pinger(c chan string) {
        for i := 0; ; i++ {
            c <- "ping" // 송신
        }
    }
    func printer(c chan string) {
        for {
            msg := <- c // 수신
            fmt.Println(msg)
            time.Sleep(time.Second * 1)
        }
    }
    func main() {
        var c chan string = make(chan string) // 송수신 channel 생성

        go pinger(c)
        go printer(c)

        var input string
        fmt.Scanln(&input)
    }
    ```
- channel 타입에 `<-` 연산자를 지정해 수신 또는 송신 전용 channel을 생성할 수 있다. 기본적으로 해당 연산자를 사용하지 않고 생성하면 양방향 channel 타입을 의미한다. 만약 송신 전용 channel에 대해 수신 연산을 수행할 경우 compile 오류가 발생한다(반대도 동일). 아래는 예시다.
    ``` go
    var c1 chan<- string = make(chan<- string) // 송신 전용
    var c2 <-chan string = make(<-chan string) // 수신 전용
    ```
- `switch`문과 유사한 `select`문은 준비된 channel(수신 channel의 경우 받을 메시지가 있는 경우, 송신 channel의 경우 해당 chanenl을 수신대기하고 있거나 buffered channel의 경우 buffer에 메시지가 꽉차지 않은 경우)의 case문을 실행한다. 하나 이상의 channel이 준비되면 어느 channel로부터 메시지를 받을지 무작위로 선택한다. 준비된 channel이 없으면 사용 가능해질 때까지 문장 실행이 차단된다. default case는 준비된 channel이 없을 경우 즉시 실행된다. 아래는 예시다. select 문을 계속 실행하기 위해 for {...}문 내에서 사용할 수 있다.
    ``` go
    for {
        select {
        case msg1 := <- c1:
            fmt.Println("Message 1", msg1)
        case msg2 := <- c2:
            fmt.Println("Message 2", msg2)
        case <- time.After(time.Second):
            fmt.Println("timeout")
        default:
            fmt.Println("nothing ready")
        }
    }
    ```
- `make()` 함수를 사용해 channel을 생성할 때 두 번째 매개변수에 크기를 지정해 buffered channel을 생성할 수 있다. buffered channel은 메시지를 수신, 송신 시 channel이 꽉차있지 않은 이상 기다리지 않는다는 차이점이 있다.
    ``` go
    ch := make(chan int, 100)
    ```

### [예제로 배우는 Go 프로그래밍](http://golang.site/)
- 함수가 결과와 에러를 함께 리턴한다면, 이 에러가 nil 인지를 체크해서 에러가 없는지를 체크할 수 있다. 또 다른 에러 처리로서 error의 타입을 체크(switch문)해서 에러 타입별로 별도의 에러 처리를 하는 방식이 있다.

### [Go by Example](https://gobyexample.com/)
- Arrays
    - 컴파일러 시점에 배열의 개수를 결정하도록 정의할 수 있다.
    - `:`를 사용해 인덱스를 지정할 수 있다.
    ``` go
    b = [...]int{1, 2, 3, 4, 5}
    b = [...]int{100, 3: 400, 500}
    ```
- Slices
    - Slices 내장 package는 slice를 위한 다양한 기능을 제공한다.
- Strings and Runes
    - 문자열은 읽기 전용 byte slice(`[]byte`) 타입이다. go 언어의 문자열은 기본적을 UTF-8 인코딩을 사용한다. 다른 언어에서는 문자열은 character로 구성되지만 go 언어에서는 rune(int32 타입의 alias)이라는 타입을 사용한다.
    - rune literal은 single quote로 표현 가능하다. 예를 들어 'a', '1' '한'
- Structs
    - 일회성의 struct 타입을 정의 및 사용하기 위해 anonymous struct를 사용할 수 있다. 이러한 기법을 table-driven style이라고 부르며 주로 테스트 코드에서 사용한다.
        ``` go
        newCar := struct {
            make    string
            model   string
            mileage int
        }{
            make:    "Ford",
            model:   "Taurus",
            mileage: 200000,
        }
        ```
- Enums
    - enum은 미리 정의된 이름을 갖는 값들의 고정된 집합을 나타내는 특별한 데이터 타입이다. 일반적으로 상수로 원시적인 정수나 문자열을 사용하는 것보다 훨씬 읽기 쉽고 체계적인 방식으로 관련된 상태들을 정의할 때 사용된다. golang에서는 enum을 내장 기능으로 지원하지 않는 대신 const, iota, 사용자 type을 사용해 구현할 수 있다. 주의할 점은 iota가 항상 0 부터 시작하기 때문에 enum의 첫 값을 iota + 1으로 해서 1 부터 시작하도록 하던지 아니면 첫 enum을 Unknown이나 Nil을 나타내는데 사용해야 한다. golang의 zero value 때문에 예상치 못한 부작용이 생길 수 있기 때문이다.
        ``` go
        // 1. enum을 위한 커스텀 type 정의
        type TrafficLightState int

        // 2. const와 iota를 사용하여 상수 선언
        const (
        	Red    TrafficLightState = iota + 1  // Red는 1
        	Yellow                               // Yellow는 2
        	Green                                // Green은 3
        )

        // 3. 커스텀 문자열 표현을 위한 String 메서드
        func (s TrafficLightState) String() string {
        	switch s {
        	case Red:
        		return "빨간불"
        	case Yellow:
        		return "노란불"
        	case Green:
        		return "초록불"
        	default:
        		return fmt.Sprintf("알 수 없는 신호등 상태 (%d)", s)
        	}
        }
        ```
- Generics
    - 아래는 slices standardy library의 Index 함수 예시다.
        ``` go
        func Index[S ~[]E, E comparable](s S, v E) int {...}
        ```
    - 아래는 단일 연결 리스트 예시다.
        ``` go
        type List[T any] struct {
            head, tail *element[T]
        }

        type element[T any] struct {
            next *element[T]
            val  T
        }

        func (lst *List[T]) Push(v T) {
            if lst.tail == nil {
                lst.head = &element[T]{val: v}
                lst.tail = lst.head
            } else {
                lst.tail.next = &element[T]{val: v}
                lst.tail = lst.tail.next
            }
        }
        ```
- Range over Iterators
    - golang은 다른 일부 언어들처럼 Iterator라는 이름의 전용 interface나 내장 타입이 명시적으로 존재하지 않았다. 대신 강력한 동시성 기능과 함수형 프로그래밍 스타일을 활용해 interator와 유사한 동작을 구현할 수 있었따. 하지만 Go 1.23부터는 standard library에 iter package가 공식적으로 추가되면서 iterator 개념이 훨씬 더 명시적이고 표준화된 방식으로 제공되기 시작했다. 이는 golang generics에 이어 언어의 표현력을 한층 더 확장한 중요한 변화다(콜백 함수란, 다른 함수의 인자로 넘겨져서 그 함수 안에서 특정 조건이나 이벤트가 발생했을 때 호출되는 함수를 말한다).
        ``` go
        // Numbers 함수가 iter.Seq 타입의 함수를 반환합니다.
        func Numbers(max int) iter.Seq[int] {
            // 이 반환되는 익명 함수가 iter.Seq의 실제 구현체입니다.
            // 이 함수는 'yield func(int) bool' 이라는 매개변수를 받습니다.
            return func(yield func(int) bool) { // <-- 여기서 'yield'는 콜백 매개변수입니다.
                for i := 0; i < max; i++ {
                    // 개발자는 이 'yield' 콜백 매개변수를 호출하여 'i' 값을 '내보냅니다'.
                    if !yield(i) { // yield가 false를 반환하면 소비자가 더 이상 값을 원하지 않으므로 중단합니다.
                        return
                    }
                }
            }
        }

        func main() {
            for num := range Numbers(5) { // <--- 여기가 소비하는 쪽입니다.
                fmt.Printf("%d", num) // 01234
            }
        }
        ```
        - iter package의 `type Seq[V any] func(yield func(V) bool)`로 정의된 Seq 타입은 `for...range 문`에서 하나의 값을 반환할 때 사용된다. for _, v := range slice와 유사하다. Seq 타입은 함수이며 매개변수로 callback 함수를 매개변수로 전달 받는다. callback 함수는 개발자가 정의하지 않으며 for...range 문 호출 시, go runtime이 내부적으로 yield func(int) bool 시그니처를 갖는 익명 함수(실제 callback 함수)를 만들어서 전달한다.
            ``` go
            type (
            	Seq[V any]     func(yield func(V) bool)
            	Seq2[K, V any] func(yield func(K, V) bool)
            )
            ```
- Errors
    - golang은 에러를 별도의 함수 반환 값을 통해 명시적으로 전달한다. 일반적으로 마지막 반환 값으로 사용한다. 에러는 universe block에 정의된 error 타입(basic interface)으로 표현한다. standard library에 포함된 errors package는 error를 위한 다양한 기능을 제공한다.
        ``` go
	    err1 := errors.New("error1")
	    err2 := fmt.Errorf("error2: [%w]", err1)
	    fmt.Println(err2)
	    fmt.Println(errors.Unwrap(err2))
        ```
        - `func New(text string) error` 함수는 문자열 메시지가 포함된 error을 생성 및 반환한다.
        - error e가 `Unwrap() error` 또는 `Unwrap() []error` method를 갖고 있는 경우 다른 error를 wrapping할 수 있다. wrapping을 사용해 여러 error 간 chain을 만들어 error의 root cause 등의 정보를 포함할 수 있다. error wrapping의 가장 쉬운 방법은 fmt package의 `func Errorf(format string, a ...any) error` 함수를 이용하는 것이다.
        - `func As(err error, target any) bool`(err 또는 해당 cahin에서 target에 할당할 수 있는 error가 있는지 확인), `func Is(err, target error) bool`(err 또는 해당 cahin에서 target과 동일한 error가 있는지 확인) 함수는 재귀적으로 Unwrap method를 실행해 검사한다.
        - `func Unwrap(err error) error` 함수는 매개변수 err의 Unwrap method를 호출한 결과 값을 반환한다.
- Custom Errors
    - custom error를 직접 구현하는 경우 보통 Error postfix를 붙인다.
- Testing and Benchmarking
    - go의 standard library에 포함된 testing package는 go package를 위한 자동화된 테스트 코드를 작성할 수 있도록 도와준다. 이 package는 주로 `go test` 명령어와 함께 사용된다. `go test` 명령어가 테스트로 인식하려면, 함수는 다음과 같은 signature를 따라야 한다.
        ``` go
        func TestXxx(*testing.T)
        ```
        - Xxx 부분은 소문자로 시작해선 안 된다. 보통 TestAddition, TestGetUserByID처럼 테스트의 목적을 설명하는 이름으로 짓는다.
        - 함수는 `*testing.T` 타입의 단일 인자를 받는다. 이 `*testing.T` 타입은 테스트 실패를 알리거나, 정보를 로깅하거나, 테스트 흐름을 제어하는 데 사용되는 메서드를 제공한다.
        - `*testing.T` 타입에는 테스트의 실패를 알리기 위한 `Error(args ...any)`, `Errorf(format string, args ...any)`, `Fail()`, `FailNow()` 메서드를 지원한다.
    - TestXxx 함수들을 포함하는 파일의 이름은 반드시 `_test.go`로 끝나야 한다 (예: my_package_test.go). 이 파일들은 `go build`(일반적인 package 빌드) 명령어 사용 시 제외된다.
        - 테스트 파일이 테스트 대상 package와 동일 package에 속하는 경우 white box 테스트라고 부른다. 동일 package이기 때문에 테스트 대상 package의 unexported name에도 접근이 가능해 내부 구현 상세를 테스트할 수 있다.
        - 테스트 파일이 `_test` 접미사가 붙은 별도의 package에 속하는 경우 black box 테스트라고 부른다. 테스트 대상 package가 abs라면, 테스트 파일은 `package abs_test` package에 포함된다. 동일 package가 아니기 때문에 테스트 대상 package의 exported name에만 접근이 가능하다.
- Channel Synchronization
    - channel을 사용해 goroutine 간 실행을 동기화할 수 있다. 아래는 main() 함수가 종료되기 전에 worker() 함수가 실행되는 goroutine이 먼저 실행이 완료되도록 하는 예시다.
        ``` go
        func worker(done chan bool) {
            fmt.Print("working...")
            time.Sleep(time.Second)
            fmt.Println("done")

            done <- true
        }

        func main() {

            done := make(chan bool, 1)
            go worker(done)

            <-done
        }
        ```
- Timeouts
    - channel, select, standard library의 time package에 포함된 `func After(d Duration) <-chan Time` 함수를 사용해 구현할 수 있다. 실제로는 standard library의 context package를 사용해 timeout을 구현하는 것이 더 일반적이다. Time 타입의 channel은 현재 시간을 메시지로 갖는다. `func Sleep(d Duration)` 함수는 현재 goroutine을 일정 시간 동안 중지한다.
        ``` go
        func main() {

            c1 := make(chan string, 1)
            go func() {
                time.Sleep(2 * time.Second)
                c1 <- "result 1"
            }()

            select {
            case res := <-c1:
                fmt.Println(res)
            case <-time.After(1 * time.Second):
                fmt.Println("timeout 1")
            }

            c2 := make(chan string, 1)
            go func() {
                time.Sleep(2 * time.Second)
                c2 <- "result 2"
            }()
            select {
            case res := <-c2:
                fmt.Println(res)
            case <-time.After(3 * time.Second):
                fmt.Println("timeout 2")
            }
        }
        ```
- Non-Blocking Channel Operations
    - select 문의 default 절을 사용해 non-blocking channel 동작을 구현할 수 있다.
        ``` go
        func main() {
            messages := make(chan string)
            signals := make(chan bool)

            select {
            case msg := <-messages:
                fmt.Println("received message", msg)
            default:
                fmt.Println("no message received")
            }

            msg := "hi"
            select {
            case messages <- msg:
                fmt.Println("sent message", msg)
            default:
                fmt.Println("no message sent")
            }

            select {
            case msg := <-messages:
                fmt.Println("received message", msg)
            case sig := <-signals:
                fmt.Println("received signal", sig)
            default:
                fmt.Println("no activity")
            }
        }
        ```
- Timers
    - standard library의 time package에 포함된 Timer 타입은 미래의 단일 이벤트를 나타낸다. `func NewTimer(d Duration) *Timer` 생성자 함수는 Timer 타입을 반환하고 duration이 지난 후 현재 시간을 해당 Timer 타입의 변수가 갖는 channel의 메시지로 전달한다.
        ``` go
        type Timer struct {
        	C <-chan Time
        	// contains filtered or unexported fields
        }
        ```
    - `func Sleep(d Duration)` 함수를 사용할 수도 있지만 Timer를 타입을 사용하면 `func (t *Timer) Stop() bool` method를 사용해 중간에 취소할 수도 있다.
        ``` go
        func main() {

            timer1 := time.NewTimer(2 * time.Second)

            <-timer1.C
            fmt.Println("Timer 1 fired")

            timer2 := time.NewTimer(time.Second)
            go func() {
                <-timer2.C
                fmt.Println("Timer 2 fired")
            }()
            stop2 := timer2.Stop()
            if stop2 {
                fmt.Println("Timer 2 stopped")
            }

            time.Sleep(2 * time.Second)
        }
        ```
- Tickers
    - standard library의 time package에 포함된 `func NewTicker(d Duration) *Ticker` 생성자 함수는 Ticker 타입을 반환하고 duration 마다 현재 시간을 해당 Ticker 타입의 변수가 갖는 channel의 메시지로 전달한다.
        ``` go
        type Ticker struct {
        	C <-chan Time // The channel on which the ticks are delivered.
        	// contains filtered or unexported fields
        }
        ```
- Worker Pools
    - golang에서는 goroutine과 channel을 사용해 worker pool을 구현할 수 있다.
        ``` go
        func worker(id int, jobs <-chan int, results chan<- int) {
            for j := range jobs {
                fmt.Println("worker", id, "started  job", j)
                time.Sleep(time.Second)
                fmt.Println("worker", id, "finished job", j)
                results <- j * 2
            }
        }

        func main() {

            const numJobs = 5
            jobs := make(chan int, numJobs)
            results := make(chan int, numJobs)

            for w := 1; w <= 3; w++ {
                go worker(w, jobs, results)
            }

            for j := 1; j <= numJobs; j++ {
                jobs <- j
            }
            close(jobs)

            for a := 1; a <= numJobs; a++ {
                <-results
            }
        }
        ```
- 
- Environment Variables
    - standard library의 os package는 환경 변수 관련 기능을 제공한다.
        - `func Environ() []string` 함수는 key=value 형식으로 표현되는 환경 변수 목록 복사본을 slice로 반환한다.
        - `func Getenv(key string) string` 함수는 매개 변수 key에 해당하는 환경 변수 값을 반환한다. 없는 환경 변수라면 빈 값을 반환한다.
        - `func Setenv(key, value string) error` 함수는 매개 변수 key, value 환경 변수를 설정한다.
- Logging
    - standard library의 log package는 자유 형식의 logger를 생성할 수 있다. 뿐만 아니라 미리 정의된 standard logger를 사용할 수 있다. standard logger는 log package의 exported(top-level) 함수를 통해 직접 사용 가능하다. standard logger는 stderr에 날짜/시간/메시지를 출력한다(로그 메시지가 newline으로 종료되지 않을 경우 자동 추가). `func SetFlags(flag int)` 함수를 사용해 standard logger의 flag bit를 설정할 수 있다.
        ``` go
        // standard logger
        log.Println("standard logger")
        log.SetFlags(log.LstdFlags | log.Lmicroseconds)
        log.Println("with micro")

        // custom logger
        	var (
        		buf    bytes.Buffer
        		logger = log.New(&buf, "logger: ", log.Lshortfile)
        	)

        	logger.Print("Hello, log file!")

        	fmt.Print(&buf)
        ```
    - standard library의 log/slog package는 structued logging을 제공한다.
        - package의 exported 함수를 통해 기본 정의된 logger를 사용할 수 있다. 해당 함수는 내부적으로 기본 정의된 logger의 method를 호출한다.
        - slog package는 3가지 중요 타입을 정의한다.
            - Logger: Log, Debug, Warn, Error와 같은 method를 제공하며, 해당 method가 호출되면 매개변수를 사용해 Record를 생성한다.
            - Record: Record는 하나의 로그 항목에 대한 모든 정보를 담고 있는 데이터 구조를 나타낸다.
            - Handler: Logger가 Record를 생성한 후, 연결된 Handler에게 전달한다. Handler는 이 Record를 받아서 자신이 정의된 방식대로 처리한다.
        - slog에는 미리 정의된 2개 handler(TextHandler(logfmt 형식), JSONHandler(json 형식))를 제공한다.
        - 자세한 내용은 https://velog.io/@als2004/go%EB%A1%9C%EA%B7%B8-%ED%8C%A8%ED%82%A4%EC%A7%80log-packageGo-Slog-vs-Zap-%EC%96%B4%EB%96%A4%EA%B2%83%EC%9D%B4-%EC%A2%8B%EC%9D%84%EA%B9%8C 정리 필요
- HTTP Client
- HTTP Server
    - HTTP router는 웹 애플리케이션에서 특정 URL 요청에 대한 처리 로직을 정의하고 연결하는 역할을 한다. 즉, 클라이언트의 요청(예: 웹 페이지 요청)을 받아 어떤 코드를 실행할지 결정하고 해당 코드를 실행하여 결과를 클라이언트에게 반환하는 과정을 관리한다. 간단히 말해, 특정 URL에 접근했을 때 어떤 페이지를 보여줄지, 어떤 데이터를 처리할지 결정하는 일종의 길잡이 역할을 수행한다.
    - package의 exported 함수를 통해 기본 정의된 HTTP 서버를 사용할 수 있다.
    - http.ServerMux struct는 http.Handler interface를 구현한 HTTP 요청 multiplexer다. 각 요청 URL의 패턴과 등록된 handler를 매칭하고 가장 유사한 handler를 통해 요청을 처리하도록 한다.
- Context

### [The Go Programming Language Specification](https://go.dev/ref/spec)
- 변수는 값을 갖는 저장 공간을 의미한다. 허용된 값의 목록은 변수의 타입에 의해 결정된다. static type은 변수 선언 시 알려진 타입이다(컴파일 시점에 결정됨). 반면 dynamic type은 interface 변수에 실제로 저장된 값의 실제 타입이다(runtime에 결정됨).
- identifier(식별자)는 변수나 타입과 같은 프로그램 엔티티의 이름을 지정한다. identifier는 하나 이상의 문자(letter)와 숫자(digit)로 이루어진 연속된 문자열로 이루어진다.
- 아래 keyword(키워드)는 golang에 의해 예약됐기 때문에 identifier로 사용할 수없다.
    ```
    break        default      func         interface    select
    case         defer        go           map          struct
    chan         else         goto         package      switch
    const        fallthrough  if           range        type
    continue     for          import       return       var
    ```
- byte는 uint8 타입, rune은 int32의 alias declaration이다.
- struct 타입은 tagging을 사용해 속성을 나타낼 수 있다. tag는 reflection interface을 통해 확인할 수 있으며 struct의 타입 식별에 영향을 미치며 이외 경우에는 무시된다. 주로 json serialization 등에 사용할 수 있다. 문자열 literal을 사용해 속성을 나타낼 수 있다.
    ``` go
    type test1 struct {
    	x, y float64 ""  // an empty tag string is like an absent tag
    	name string  "any string is permitted as a tag"
    	_    [4]byte "ceci n'est pas un champ de structure"
    }

    // A struct corresponding to a TimeStamp protocol buffer.
    // The tag strings define the protocol buffer field numbers;
    // they follow the convention outlined by the reflect package.
    type test2 struct {
    	microsec  uint64 `protobuf:"1"`
    	serverIP6 uint64 `protobuf:"2"`
    }
    ```
- interface 타입은 type set(타입 집합)을 정의한다. interface 타입의 변수는 interface type set(interface를 구현한 non-interface 타입)의 값을 저장할 수 있으며 type assertion을 통해 interface 변수가 갖고 있는 타입을 확인할 수 있다. interface는 method, 타입 목록을 가질 수 있다.
    - method 목록만 포함하는 경우 basic interface라고 부른다. basic interface의 type set은 모든 method를 구현헤야 한다.
        ``` go
        // A simple File interface.
        interface {
            Read([]byte) (int, error)
            Write([]byte) (int, error)
            Close() error
        }
        ```
        - method 이름은 유일해야 하며 blank(`_`) 이면 안된다.
        - 하나의 타입이 여러 interface를 구현할 수도 있다.
        - empty interface는 모든 non-interface 타입이 구현한다. 편의를 위해 universe block에는 empty interface의 alias declaration인 any 타입을 선언한다.
            ``` go
            // any is an alias for interface{} and is equivalent to interface{} in all ways.
            type any = interface{}
            ```
    - interface는 다른 interface를 포함하는 embedded interface를 지원한다. 이 때 이를 구현하기 위해서는 두 interface의 method를 모두 구현해야 한다. 그리고 두 interface에 중복된 이름의 method가 있을 경우 동일한 형태를 가져야 한다. interface의 type set에는 interface가 포함되지 않음을 유의해야 한다. 즉 interface는 type set을 정의하는 개념으로 사용되기 때문에 interface 자체가 type set의 목록에 포함되지 않는다. embedded interface의 경우에도 마찬가지로 type set은 관련 interface에 속한 모든 method를 구현한 non-interface 타입이다.
        ``` go
        type Reader interface {
        	Read(p []byte) (n int, err error)
        	Close() error
        }

        type Writer interface {
        	Write(p []byte) (n int, err error)
        	Close() error
        }

        // ReadWriter's methods are Read, Write, and Close.
        type ReadWriter interface {
        	Reader  // includes methods of Reader in ReadWriter's method set
        	Writer  // includes methods of Writer in ReadWriter's method set
        }

        type ReadCloser interface {
        	Reader   // includes methods of Reader in ReadCloser's method set
        	Close()  // illegal: signatures of Reader.Close and Close are different
        }
        ```
    - method 목록과 추가적으로 타입을 갖는 경우 general interface라고 부른다. general interface를 통해 단순히 메서드 집합만을 정의하는 것을 넘어 특정 타입들을 명시적으로 포함하거나 제외하는 방식으로 type set을 정의할 수 있다. general interface를 통해 해당 interface를 구현할 수 있는 타입을 제한할 수 있다.
        ``` go
        // An interface representing only the type int.
        interface {
        	int
        }

        // An interface representing all types with underlying type int.
        interface {
        	~int
        }

        // An interface representing all types with underlying type int that implement the String method.
        interface {
        	~int
        	String() string
        }

        // An interface representing an empty type set: there is no type that is both an int and a string.
        interface {
        	int
        	string
        }
        ```
        - 타입은 `T`, `~T`(underlying type이 T인 모든 타입), `T1|T2|T3`(union 연산자는 or) 형태로 명시할 수 있다.
            - `T`에는 interface를 사용할 수 있다.
            - `~T`에는 underying type만 사용할 수 있으며, interface를 사용할 수 없다.
                ``` go
                type MyInt int

                interface {
                	~[]byte  // the underlying type of []byte is itself
                	~MyInt   // illegal: the underlying type of MyInt is not MyInt
                	~error   // illegal: error is an interface
                }
                ```
            - `T`, `~T`에는 type parameter를 사용할 수 없다. non-interface 타입들은 서로 교집합이 없어야 하지만 interface 타입은 다른 타입들과 교집합이 있어도 무방하다.
                ``` go
                interface {
                	P                // illegal: P is a type parameter
                	int | ~P         // illegal: P is a type parameter
                	~int | MyInt     // illegal: the type sets for ~int and MyInt are not disjoint (~int includes MyInt)
                	float32 | Float  // overlapping type sets but Float is an interface
                }
                ```
            - `T1|T2|T3`는 type set의 전체 집합을 나타낸다.
                ``` go
                // The Float interface represents all floating-point types
                // (including any named types whose underlying types are
                // either float32 or float64).
                type Float interface {
                	~float32 | ~float64
                }
                ```
                - `T1|T2|T3`에는 추가 제약 사항이 있다.
                    - predeclared identifier인 comparable을 포함할 수 없다.
                    - method를 갖는 interface를 포함할 수 없다. comparable을 embed할 수 없다.
                    - method을 포함한 interface를 embed할 수 없다.
        - general interface는 type constraint, 다른 interface의 타입 요소로만 사용 가능하다. 반면 변수 선언, non-interface 타입의 구성 요소로 사용할 수 없다. 타입에는 자기 자신을 직접적, 간접적으로 사용할 수 없다. general interface의 유일한 목적은 generic의 type argument를 제약하는 것이다. general interface는 변수를 선언하거나, struct 필드, 컬렉션의 요소로 사용할 수 없다다. 왜냐하면 이들은 단일하고 구체적인 타입을 나타내는 것이 아니라, 타입들의 집합을 나타내기 때문이다.
            ``` go
            // general interface는 non-interface 타입의 구성 요소로 사용할 수 없다.
            type DataHolder struct {
                Value Number // 컴파일 오류!
            }
            ```
            ``` go
            // general interface는 다른 interface의 타입 구성 요소로 사용할 수 있다.
            type BasicComparable interface {
                ~int | ~string // 하나의 제약
                // 메서드 없음
            }

            type MyCombinedConstraint interface {
                BasicComparable // 제약을 임베딩
                fmt.Stringer    // 기본 인터페이스(메서드) 임베딩
                // 다른 메서드나 제약
            }

            func Process[T MyCombinedConstraint](val T) {
                // ...
            }
            ```
            ``` go
            var x Float                     // illegal: Float is not a basic interface

            var x interface{} = Float(nil)  // illegal

            type Floatish struct {
                f Float                 // illegal
            }

            // illegal: Bad may not embed itself
            type Bad interface {
                Bad
            }

                // illegal: Bad1 may not embed itself using Bad2
            type Bad1 interface {
                Bad2
            }
            type Bad2 interface {
                Bad1
            }

            // illegal: Bad3 may not embed a union containing Bad3
            type Bad3 interface {
                ~int | ~string | Bad3
            }

            // illegal: Bad4 may not embed an array containing Bad4 as element type
            type Bad4 interface {
                [10]Bad4
            }
            ```
    - 타입 T가 interfaece I를 구현(implement)하는 조건은 다음과 같다.
	    - T가 interface가 아니면서 I의 type set에 속한다.
	    - T가 interface이면서 T의 type set이 I의 type set의 부분 집합이다.
- map의 key 타입은 비교 가능한 타입( `==`, `!=`)이어야 한다. 만약 key 타입이 interface라면 dynamic type에 대해 비교 연산이 가능해야 한다.
- universe block은 모든 golang 소스 코드를 포함한다. 아래 identifier는 universe block에 선언된 predeclared identifier다. [builtin](https://pkg.go.dev/builtin) package documentation을 통해 확인할 수 있다. builtin package는 단순히 golang documentation을 위해 작성된 코드이다.
    ```
    Types:
    	any bool byte comparable
    	complex64 complex128 error float32 float64
    	int int8 int16 int32 int64 rune string
    	uint uint8 uint16 uint32 uint64 uintptr

    Constants:
    	true false iota

    Zero value:
    	nil

    Functions:
    	append cap clear close complex copy delete imag len
    	make max min new panic print println real recover
    ```
- constant declaration에서 predeclared identifier인 `iota`은 index를 나타낸다. `iota`는 주로 enum을 정의하는데 사용된다.
    ``` go
    const (
    	c0 = iota  // c0 == 0
    	c1 = iota  // c1 == 1
    	c2 = iota  // c2 == 2
    )

    const (
    	a = 1 << iota  // a == 1  (iota == 0)
    	b = 1 << iota  // b == 2  (iota == 1)
    	c = 3          // c == 3  (iota == 2, unused)
    	d = 1 << iota  // d == 8  (iota == 3)
    )

    const (
    	u         = iota * 42  // u == 0     (untyped integer constant)
    	v float64 = iota * 42  // v == 42.0  (float64 constant)
    	w         = iota * 42  // w == 84    (untyped integer constant)
    )

    const x = iota  // x == 0
    const y = iota  // y == 0
    ```
    - const 키워드 뒤에 괄호 ()를 사용해 여러 상수를 한 번에 선언할 때, 첫 번째 상수를 제외하고는 값을 할당하는 부분을 생략할 수 있다. 값을 생략하면 바로 위에서 사용된 값과 타입을 그대로 따라간다. 만약 이전 상수가 a, b = 1, 2 와 같이 두 개의 값을 가졌다면, 다음 줄에서 값을 생략하고 c, d 와 같이 두 개의 식별자를 사용해야 한다. c 만 사용하면 에러가 발생한다.
        ``` go
        const (
            a = iota // a == 0
            b        // b == 1 (iota가 증가하고, 이전 표현식 iota가 반복됨)
            c        // c == 2 (iota가 증가하고, 이전 표현식 iota가 반복됨)
        )

        const (
            x = 10
            y      // y == 10 (이전 표현식 10이 반복됨)
            z      // z == 10 (이전 표현식 10이 반복됨)
        )

        const (
            p, q = iota, iota + 10 // p == 0, q == 10
            r, s                   // r == 1, s == 11 (iota가 증가하고, 이전 표현식 목록이 반복됨)
        )

        const (
        	size int64 = 1024
        	eof        = -1  // untyped integer constant
        )
        ```
- type declaration(타입 선언)은 타입에 identifier(타입 이름)을 바인딩하는 것을 말한다. alias declaration, type definition 두 가지 종류가 있다.
    - alias declaration은 타입에 identifier(별칭)을 바인딩하는 것을 말한다. type parameter를 명시하는 경우 generic alias라고 부른다. 이 때 type parameter를 대상 타입으로 사용할 수 없다.
        ``` go
        type (
            nodeList          = []*Node     // nodeList and []*Node are identical types
            Polar             = polar       // Polar and polar denote identical types
            set[P comparable] = map[P]bool  // generic alias
            A[P any]          = P           // illegal: P is a type parameter
        )
        ```
    - type definition은 기존 타입과 동일한 기능을 제공하지만 독립적인 새로운 타입을 identifier(타입 이름)에 바인딩하는 것을 말한다. type parameter를 명시하는 경우 generic type이라고 부른다. generic type에 대한 method 정의시 receiver에도 동일한 type parameter를 명시해야 한다. 이 때 type parameter를 대상 타입으로 사용할 수 없다. 그리고 새로운 타입을 정의할 때 기존 타입은 컴파일 시점에 알려진 타입이어야 한다.
        ``` go
        type (
        	Point struct{ x, y float64 }  // Point and struct{ x, y float64 } are different types
        	polar Point                   // polar and Point denote different types
        )

        type TreeNode struct {
        	left, right *TreeNode
        	value any
        }

        type Block interface {
        	BlockSize() int
        	Encrypt(src, dst []byte)
        	Decrypt(src, dst []byte)
        }

        type T[P any] P    // illegal: P is a type parameter

        func f[T any]() {
        	type L T   // illegal: T is a type parameter declared by the enclosing function
        }
        ```
- type parameter(타입 파라미터)는 generic 함수, generic 타입의 type parameter 목록을 나타낸다. type parameter는 type constraint(타입 제약)이 있으며 이는 type parameter에 대한 일종의 메타 타입 역할을 수행한다. type parameter는 일반적으로 여러 타입의 집합을 나타내지만 컴파일 시점에는 단일 타입을 나타낸다. 사용자는 generic 함수, generic 타입 사용 시 type argument(타입 매개변수)를 명시해야 하지만, 컴파일러가 타입을 추측할 수 있는 경우 type argument를 생략할 수 있다.
    ``` go
    func PrintKeyValue[K ~string, V int](k K, v V) {
	    fmt.Printf("key: %v, value: %v\n", k, v)
    }

    func main() {
    	PrintKeyValue[string, int]("first", 1) // can call PrintKeyValue("first", 1)
    }
    ```
    - 단일 type parameter를 사용하는 경우 파싱의 모호성이 생길 수 있다. 이러한 특이 케이스에서 type parameter가 아닌 일반적인 표현식으로 해석될 수도 있다. 이를 해결하기 위해 interface로 감싸거나 끝에 ,를 추가할 수 있다.
        ``` go
        type T[P *C] …   // `P *C`가 포인터 타입으로 해석될 수도 있음
        type T[P (C)] …  // `P (C)`가 타입 변환으로 해석될 수도 있음
        type T[P *C|Q] … // `*C | Q`가 비트 연산자로 해석될 수도 있음

        type T[P interface{*C}] …
        type T[P *C,] …
        ```
    - generic 타입 T의 type parameter 목록에 type constraint T를 직접적으로 또는 간접적으로 참조할 수 없다. 이는 순환 참조, 타입의 정의가 완전히 확정되지 않은 상태에서 자신을 의존하는 모순적인 상황을 만들 수 있기 때문에 금지된다.
        ``` go
        type T1[P T1[P]] …                    // illegal: T1 refers to itself
        type T2[P interface{ T2[int] }] …     // illegal: T2 refers to itself
        type T3[P interface{ m(T3[int])}] …   // illegal: T3 refers to itself
        type T4[P T5[P]] …                    // illegal: T4 refers to T5 and
        type T5[P T4[P]] …                    //          T5 refers to T4

        type T6[P int] struct{ f *T6[P] }     // ok: reference to T6 is not in type parameter list
        ```
- type constraint는 interface로 type parameter의 type argument로 사용할 수 있는 타입과 연산를 제한한다. 표현식 `interface{E}`와 같이 표현하며 E가 method가 아닌 경우 단순하게 `E`로 표현할 수 있다.
    ``` go
    [T []P]                      // = [T interface{[]P}]
    [T ~int]                     // = [T interface{~int}]
    [T int|string]               // = [T interface{int|string}]
    ```
    - type parameter가 아닌 interface 변수는 비교가 가능하지만 strictly comparable하지 않기 때문에 comparable interface를 구현하지는 않는다. 하지만 type parameter로서 interface는 comparable interface를 충족할 수 있다(interface의 구현(implementation)과 type constraint에 대한 충족(satisfaction)은 다른 의미로 생각해야 함).
        ``` go
        int                          // implements comparable (int is strictly comparable)
        []byte                       // does not implement comparable (slices cannot be compared)
        interface{}                  // does not implement comparable (see above)
        interface{ ~int | ~string }  // type parameter only: implements comparable (int, string types are strictly comparable)
        interface{ comparable }      // type parameter only: implements comparable (comparable implements itself)
        interface{ ~int | ~[]byte }  // type parameter only: does not implement comparable (slices are not comparable)
        interface{ ~struct{ any } }  // type parameter only: does not implement comparable (field any is not strictly comparable)
        ```
    - gerneral interface인 comparable은 strictly comparable non-inferface 타입들이 구현한다. 즉, 이 interface를 구현한 타입은 비교 연산자 `==`, `!=`를 사용할 수 있는 타입으로 `bool`, 숫자(`int`, `uint`, `float32`, `complex64` 등), string, pointer, channel, 일부 struct(필드가 모두 comparable 타입인 경우), 일부 배열(comparable 타입 배열인 경우)이 있다. inteface 타입은 비교가 가능하지만 strictly comparable하지 않기 때문에 comparable을 구현하지 않는다.
    - comparable은 general interface이기 때문에 type constraint로만 사용 가능하며 변수의 타입으로는 사용할 수 없다.
    - type argument T가 type constraint C를 충족한다는 의미는 T가 C의 type set에 포함된다는 것이다. 즉, T가 C를 구현하는 것을 말한다. 예외적으로 비교 가능한 type argument는 strictly comparable type constraint(comparable interface)을 충족한다. 이러한 예외 규칙으로 인해 비교 연산 수행 시 runtime panic이 발생할 수 있다.
        ``` go
        type argument      type constraint                // constraint satisfaction

        int                interface{ ~int }              // satisfied: int implements interface{ ~int }
        string             comparable                     // satisfied: string implements comparable (string is strictly comparable)
        []byte             comparable                     // not satisfied: slices are not comparable
        any                interface{ comparable; int }   // not satisfied: any does not implement interface{ int }
        any                comparable                     // satisfied: any is comparable and implements the basic interface any
        struct{f any}      comparable                     // satisfied: struct{f any} is comparable and implements the basic interface any
        any                interface{ comparable; m() }   // not satisfied: any does not implement the basic interface interface{ m() }
        interface{ m() }   interface{ comparable; m() }   // satisfied: interface{ m() } is comparable and implements the basic interface interface{ m() }
        ```
- 비교 연산자 `==`, `!=`는 비교 가능한 타입에 사용할 수 있다.
    - type parameter가 아닌 interface 변수는 비교가 가능하다(type parameter는 strictly comparable인 경우에만 비교 가능). 두 interface가 동일하다는 의미는 두 interface 모두 `nil`이거나 dynamic 타입, 값이 동일한 경우다. dynamic 타입이 비교 가능하지 않을 경우 runtime panic이 발생할 수 있다.
    - slice, map, 함수 타입은 `nil` predeclared identifier와만 비교할 수 있다.
    - type parameter는 strictly comparable(비교가 가능한 타입이면서 interface 타입이 아니고 interface 타입으로 구성되지 않는 타입)한 경우에만 비교할 수 있다.
- `*`, `<-` 연산자로 시작하는 타입과 func keyword로 시작하지만 리턴 목록이 없는 타입은 모호성을 피하기 위해 필요한 경우 소괄호로 표현해야 한다.
    ``` go
    *Point(p)        // same as *(Point(p))
    (*Point)(p)      // p is converted to *Point
    <-chan int(c)    // same as <-(chan int(c))
    (<-chan int)(c)  // c is converted to <-chan int
    func()(x)        // function signature func() x
    (func())(x)      // x is converted to func()
    (func() int)(x)  // x is converted to func() int
    func() int(x)    // x is converted to func() int (unambiguous)
    ```
- built-in 함수는 predeclared identifier다. 일반적인 함수와 동일하지만 몇 built-in 함수는 매개변수로 타입을 요구한다. 그리고 함수의 값으로 사용할 수 없다.
- `import` 문을 사용해 다른 package의 exported 필드를 사용할 수 있다.
    ``` go
    Import declaration          Local name of Sin

    import   "lib/math"         math.Sin
    import m "lib/math"         m.Sin
    import . "lib/math"         Sin
    ```
    - import문에는 package의 경로를 명시하며 다음과 같은 규칙이 있다.
        - package 이름과 경로가 일치하면 별칭을 사용하지 않아도 되며 package 이름으로 접근할 수 있다.
        - package 이름과 경로가 일치하지 않으면 별칭이 필요하다.
        - 별칭 `.`은 해당 package의 exported 필드를 자신의 소스 파일 block에 선헌하기 때문에 수식어 없이 접근해야 한다.
        - 별칭 `_`을 사용해 package의 exported 필드를 사용하지 않지만 초기화 목적으로 사용할 수 있다.

### [Developing modules](https://go.dev/doc/)
- Organizing a Go module
    - go 프로젝트는 `import`-ing 코드, `go install`코드로 구성될 수 있다.
        - Basic package or command: 단일 module과 단일 package로 구성된 경우 module의 루트 디렉토리에서 모든 파일을 관리
            ```
            // example 1
            project-root-directory/
              go.mod
              modname.go
              modname_test.go

            // example 2
            project-root-directory/
              go.mod
              modname.go
              modname_test.go
              auth.go
              auth_test.go
              hash.go
              hash_test.go
            ```
        - Package or command with supporting packages: 구성이 복잡해지는 경우 일부 기능은 별도의 package로 구성한다. 처음에는 `internal` package를 사용해 module 외부로 해당 기능을 노출하지 않는다.
            ```
            project-root-directory/
              internal/
                auth/
                  auth.go
                  auth_test.go
                hash/
                  hash.go
                  hash_test.go
              go.mod
              modname.go
              modname_test.go
            ```
        - Multiple packages: module이 여러 importable package를 제공하는 경우 각 package를 별도 디렉토리로 구성한다.
            ```
            project-root-directory/
              go.mod
              modname.go
              modname_test.go
              auth/
                auth.go
                auth_test.go
                token/
                  token.go
                  token_test.go
              hash/
                hash.go
              internal/
                trace/
                  trace.go
            ```
        - Multiple commands: 단일 module에 여러 cli(또는 executable program) 프로그램을 구성하는 경우 디렉토리를 구분한다.
            ```
            project-root-directory/
              go.mod
              internal/
                ... shared internal packages
              prog1/
                main.go
              prog2/
                main.go
            ```
        - Packages and commands in the same repository: importable package와 cli(또는 executable program) 프로그램을 포함하는 경우 아래와 같이 구분한다.
            ```
            project-root-directory/
              go.mod
              modname.go
              modname_test.go
              auth/
                auth.go
                auth_test.go
              internal/
                ... internal packages
              cmd/
                prog1/
                  main.go
                prog2/
                  main.go
            ```
        - Server project: 서버 프로젝트는 대부분 자체 사용 목적을 갖기 때문에 internal package를 사용한다. 따라서 서버 관련 로직을 구현하는 package는 `internal` 디렉토리에 두는 것을 권장한다. 그리고 서버 프로젝트의 경우 대부분 go 관련 파일이 아닌 다른 많은 디렉토리가 있을 가능성이 높으므로 cli(또는 executable program) 프로그램을 `cmd` 디렉토리에 두는 것을 권장한다.