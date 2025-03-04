### about docs
- go1.24.0 버전 기준으로 작성

### [A Tour of Go](https://go.dev/tour/list)
- golang은 package로 구성되며 golang으로 개발된 프로그램은 main package을 통해 실행된다.
- `import` 키워드를 사용해 package가 속한 경로를 import할 수 있다. 편의를 위해 import 경로의 마지막 구성 요소를 package 이름으로 사용하는 것이 일반적이다.
- `import` 키워드를 여러 번 사용해 여러 package를 import할 수도 있지만 import (...)와 같이 사용하는 것을 권장한다.
- package 내에서 대문자로 시작되는 이름을 갖는 경우 해당 package 밖에서도 참조가 가능하며 이를 exported name이라고 한다. 반대로 소문자로 시작되는 이름을 갖는 경우 package 내부에서만 참조가 가능하다. 내장 타입은 대문자로 시작하지 않아도 접근할 수 있다.
- 함수의 반환 값에 이름을 지정하는 경우 함수의 최상단에서 정의된 변수로 취급된다. 함수에서는 `return` 키워드만 사용해도 반환이된다. `return` 키워드를 생략하는 것은 불가능하다. 이를 naked return이라고 부르며 짧은 길이의 함수에서만 사용하는 것을 권장한다.
- `var` 키워드를 사용해 변수를 선언할 수 있다. 변수 선언 시 초기화도 수행하면 변수의 타입을 생략할 수 있다. 
- 함수 내부에서는 `:=` short assignment statement만 이용해 변수를 선언 및 초기화할 수 있다. 함수 외부에서는 항상 키워드로 시작해야 하기 때문에 사용할 수 없다.
- 변수 선언 시 초기화를 하지 않으면 zero value가 할당된다. 숫자 타입일 경우 0, boolean 타입일 경우 false, string 타입일 경우 ""
- 묵시적 형변환이 불가능하며 항상 `T(v)` 표현식을 사용해 타입 변환을 수행해야 한다.
- 변수 선언 시 타입을 지정하지 않는 경우 변수의 타입은 초기화 값으로부터 추론된다. 초기화 값이 명시적인 타입을 갖는 경우 변수는 동일한 타입을 갖는다. 하지만 타입이 없는 숫자 상수의 경우 정밀도에 따라 `int`, `float64`, `complex128` 타입으로 추른된다.
- `const` 키워드를 사용해 상수를 선언할 수 있다. 상수는 문자, 문자열, boolean, 숫자 타입일 수 있다(`:=` short assignment statement는 변수 선언시 사용되기 때문에 상수에 대해서는 사용할 수 없다).
- 타입이 지정되지 않은 상수는 사용되는 문맥에 따라 필요한 타입을 갖게된다.
- 반복문을 위해 `for` 키워드만 지원한다. ()로 감싸지 않아도 되지만 블럭에는 {}가 항상 필요하다.
- `if`문도 `for`문과 동일하게 ()로 감싸지 않아도 되지만 블럭에는 {}가 항상 필요하다. `if`문에 short assignment statement를 사용할 수 있으며 해당 변수는 `else if`, `else if` 블럭에서도 사용할 수 있다.
- `switch (테스트 표현식)`문에는 `break` 문을 명시할 필요없다. 테스트 표현식을 생략하는 경우 `switch true`와 동일하다.
- `defer`문은 함수의 반환까지 함수의 실행을 연기한다. 함수의 매개변수는 `defer` 문에서 결정되지만 호출은 지연된다. 지연된 함수는 stack에 push되며 LIFO로 실행된다.
- `*T`는 포인터 타입으로 zero value은 `nil`이다. `&` 연산자는 피연산자의 포인터를 생성한다. `*` 연산자는 포인터가 가리키는 값을 나타낸다. 포인터에 대한 산술 연산자는 없다.
- struct는 필드의 집합으로 각 필드는 .을 이용해 접근한다. struct pointer의 경우 `(*p).X`을 통해 접근 가능하지만 번거롭기 때문에 `p.X`와 같이 접근하는 것을 허용한다. `type 이름 struct{필드목록}` 문을 이용해 struct 타입을 정의할 수 있다.
- struct 리터럴은 순서대로 값을 나열하거나 `Name:`처럼 필드의 이름과 값을 순서 상관없이 나열해 표현할 수 있다. 
- `[n]T` 타입은 배열이다. 배열의 크기는 배열 타입의 일부기 때문에 크기를 조정할 수 없다.
- `[]T` 타입은 슬라이스로 동적으로 크기 조절이 가능하며 배열보다 더 일반적으로 사용한다. 슬라이스는 내부적으로 배열을 가리키는 포인터, 길이, 용량 정보를 저장하며 zero value은 `nil`이다. 배열 또는 슬라이스 변수 a에 대해 `a[low:high]`의 표현식을 사용해 슬라이싱해 슬라이스 값을 얻을 수 있다.
- 슬라이스의 길이는 슬라이스를 통해 접근할 수 있는 요소의 길이를 나타내며 용량은 실제 슬라이스가 참조가 있는 배열 길이를 나타낸다(슬라이스가 가리키는 첫 인덱스부터 배열의 마지막 인덱스까지). 슬라이스의 길이는 `len(s)`, 용량은 `cap(s)` 내장 함수를 통해 얻을 수 있다.
- `make([]T, len, (cap))` 함수를 사용해 zero value를 갖는 슬라이스 값을 생성할 수 있다.
- `func append(s []T, vs ...T) []T` 함수를 사용해 슬라이스 마지막 인덱스 뒤에 값을 계속 추가할 수 있다. 기존 배열의 크기가 작으면 크기가 더 큰 배열을 생성 및 할당한다.
- `for i, v := range sli {}`와 같은 표현식을 사용해 slice에 대해 반복문을 사용할 수 있다. 변수 i는 인덱스, v는 해당 인덱스의 복사된 값을 갖는다. `_` 표현식을 사용해 할당을 하지 않을 수 있으며 v는 생략할 수도 있다.
- `map[K]V` 타입은 맵이다. zero value은 nil이며 리터럴은 struct와 다르게 key를 생략할 수 없다. `m[key]` 표현식을 사용해 map에 저장된 요소에 접근할 수 있다. `elem = m[key]` 표현식은 변수 m에 key가 없을 경우 오류가 발생한다. 반면 `elem, ok = m[key]` 표현식은 변수 m이 key가 있으면 ok 변수에 true, 없다면 false 값을 갖는다. key가 없는 경우 elem 변수에 zero value가 할당된다. `delete(m, key)` 함수를 사용해 맵 변수 m에서 key를 삭제할 수 있다.
- golang은 class 개념이 없지만 타입에 method를 정의할 수 있다. method는 receiver라는 인자를 받는 함수다. 함수와 다른 점은 단순히 reciver라는 특별한 인자가 있다는 것이며 기능적으로는 함수와 동일하다. method는 type이 정의된 package 내에서만 선언할 수 있다. 그렇기 때문에 int와 같은 내장 타입에는 method를 사용자가 정의할 수 없다. 물론 `type` 키워드를 사용해 int를 다시 한번 정의한 후 method를 정의할 수 있다.
- method 선언 시 pointer receiver를 사용할 수 있다. 이는 함수가 포인터 타입의 인자를 통해 실제 값을 변경하는 것과 같이 동일하게 동작한다. `*T` 타입에 대한 method를 `T` 타입의 변수에서 접근할 때 `(&T).X`와 같이 접근할 수 있지만 번거롭기 때문에 `T.X`와 같이 접근하는 것을 허용한다. 반대인 경우에도 `(*T).X` 대신 `T.X`와 같이 접근 가능하다.
- method를 통해 실제 값을 변경이 필요하거나, 크기가 큰 경우 method 호출 시 항상 값을 복사하는 불필요성을 피하기 위해 pointer receiver를 사용한다. 한 타입에 method를 선언할 때 pointer receiver, value receiver를 섞어 쓰는 것을 권장하지 않는다.
- interface는 struct와 유사하지만 필드의 집합이 아니라 method의 집합을 나타낸다. interface 타입은 method를 모두 구현한 타입을 값으로 가질 수 있다. `type 이름 interface{필드목록}` 표현식을 사용해 선언할 수 있다.
- interface를 구현한다는 것은 어떤 키워드를 통해 명시적으로 수행하는 것은 아니며 단순히 interface 타입에 포함된 모든 method를 선언해 암묵적으로 구현한다. interface의 선언식과 구현이 동일 package일 필요는 없다. interface는 `(value, type)` 값을 갖는다고 생각할 수 있다.
- interface의 zero value은 nil이며 method 호출 시 runtime error가 발생한다.
- methods를 명시하지 않은 interface를 empty interfac라고 한다. empty interface는 모든 타입, 모든 값을 가질 수 있다. empty interface는 따로 `type` 키워드를 통해 선언할 필요가 없으며 변수 선언 시 `var i interface{}`와 같이 사용할 수 있다. empty interface는 알려지지 않은 타입의 값을 다룰 때 사용된다.
- type assertion은 interface가 nil이 아니며 타입 T임을 확인하는 것을 말한다. `i.(T)` 표현식은 interface i가 T 타입임을 나타내며 i가 가리키는 T 타입의 변수를 반환한다. 만약 nil이나 T 타입이 아닐 경우 오류가 발생한다. `t, ok := i.(T)` 표현식은 두 번째 반환 값을 통해 타입 T가 맞는지에 따른 boolean 값을 반환한다.
- type switch 문은 switch의 테스트 표현식에 `i.(type)`을 사용해 interface 변수 i에 대한 내장 타입에 대한 case 문을 작성할 수 있다.
- 가장 흔한 interface는 fmt package에 있는 Striner interface다. fmt package는 값을 출력하기 위해 Stringer interface의 String method를 호출한다.
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
- golang은 제네릭 함수, 제네릭 타입을 통해 제네릭스를 제공한다. 함수에 타입 파라미터를 추가함으로써 제네릭 함수를 정의할 수 있다. 아래는 예시다.
    ``` go
    func Index[T comparable](s []T, x T) int
    ```
    뿐만 아니라 `type` 키워드에 타입 파라미터를 추가함으로써 제네릭 타입을 정의할 수 있다.
    ``` go
    type List[T any] struct {
	    next *List[T]
	    val  T
    }
    ```

### [GO 프로그래밍 입문](https://codingnuri.com/golang-book/)
- 문자열은 "(double quote), `(backtick)을 사용해 표현할 수 있다. double quote로 표현하는 문자열은 줄바꿈을 포함할 수 없으며 이스케이프 문자열 사용할 수 있다. backtick으로 표현하는 문자열은 줄바꿈을 포함할 수 있으며 이스케이프 문자열을 지원하지 않는다.
- 문자열은 바이트로 표현되기 때문에 `"hello world"[4]`와 같이 인덱스를 통해 접근 가능하다. 문자열에 `len(s)` 함수를 사용해 길이를 확인할 수 있다.
- 함수 정의 시 가변 인자를 사용할 수 있다. 이 때 해당 가변 인자는 함수 내에서 배열로 접근할 수 있다. 함수 호출 시에는 가변 인자에 개별 인자, 배열, 슬라이스를 사용할 수 있다.
    ``` go
    // fmt.Println() 함수 예시
    func Println(a ...interface{}) (n int, err error){...}
    ```
- `panic(v any)` 내장 함수는 현재 함수를 즉시 멈추고 현재 함수에 defer 함수들을 모두 실행한 후 즉시 리턴한다(런타임 오류). 이러한 panic 모드 실행 방식은 다시 상위함수에도 똑같이 적용되고, 계속 콜스택을 타고 올라가며 적용된다. 그리고 마지막에는 프로그램이 에러를 내고 종료하게 된다. `recover()` 내장 함수는 panic() 함수에 의한 패닉 상태를 중단하고 panic() 함수 호출 시 전달했던 인자를 반환한다. panic() 함수 호출 시 런타임 에러가 발생해 즉시 호출 중이던 함수가 종료되기 때문에 recover() 함수를 defer문과 사용해야 한다.
- `new(Type)` 내장 함수를 사용해 매개변수 타입에 대해 메모리를 할당(zero value를 할당)하고 포인터를 반환한다. 
- struct의 필드로 struct를 사용할 수 있다. embedded type은 필드로 struct를 사용할 때 필드의 이름을 지정하지 않으면 된다. embedded type은 해당 struct를 통해 직접 접근 가능하다. 구조체 생성 시 필드 명은 해당 필드 타입을 그대로 사용하면 된다.
    ``` go
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
    
    fmt.Println(s1.No, s1.ClassInfo.No) // 10 1
    ```
-

### [예제로 배우는 Go 프로그래밍](http://golang.site/)
- 함수가 결과와 에러를 함께 리턴한다면, 이 에러가 nil 인지를 체크해서 에러가 없는지를 체크할 수 있다. 또 다른 에러 처리로서 error의 타입을 체크(switch 문)해서 에러 타입별로 별도의 에러 처리를 하는 방식이 있다.
- 하나 이상의 작업을 동시에 진행하는 것을 동시성(concurrency)라 한다. golang에서는 goroutine, channel을 통해 동시성을 지원한다.
- `go` 키워드를 사용해 goroutine을 생성할 수 있다. main 함수도 goroutine에서 실행되며 main 함수가 실행되면 프로그램의 종료로 이어지기 때문에 다른 gorutine이 모두 종료된 후 main 함수의 goroutine을 종료하도록 해야 한다. channel은 두 goroutine이 서로 통신하고 실행 흐름을 동기화하는 것을 지원한다.