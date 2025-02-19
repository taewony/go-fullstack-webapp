### 1) `text/template` vs `templ` 패키지를 사용하는 HTTP 핸들러의 차이  
Go에서 `text/template`과 `templ`을 사용할 때 주요 차이점은 **템플릿 정의 방식**, **타입 안정성**, **컴파일 시점 검증** 등이 있습니다.

#### `text/template`을 사용하는 경우  
`text/template`은 기본적으로 HTML 파일을 파싱하고 데이터를 삽입하는 방식으로 작동합니다.

```go
package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title string
	Body  string
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title: "Hello, Go!",
		Body:  "This is a Go template example.",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

**특징**  
- `.html` 파일을 별도로 관리하고 `template.ParseFiles()`로 로드해야 함.
- `Execute()` 호출 시 런타임에 데이터 바인딩.
- 문법 오류가 있어도 런타임에야 알 수 있음.

---

#### `templ`을 사용하는 경우  
`templ`은 `.templ` 파일에서 Go 코드와 함께 타입 안정성을 유지하면서 HTML을 정의합니다.

1. 먼저 `hello.templ` 파일을 작성합니다.

```templ
@use templ

Hello(title string, body string) {
	<!DOCTYPE html>
	<html>
	<head>
		<title>{title}</title>
	</head>
	<body>
		<h1>{title}</h1>
		<p>{body}</p>
	</body>
	</html>
}
```

2. `go generate`를 실행하여 `hello.templ.go`를 생성합니다.

```sh
go generate ./...
```

3. HTTP 핸들러에서 호출:

```go
package main

import (
	"net/http"

	"example.com/templates" // 생성된 .templ.go 파일을 포함하는 패키지
)

func handler(w http.ResponseWriter, r *http.Request) {
	templates.hello("Hello, Go!", "This is a Go templ example.").Render(r.Context(), w)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

**특징**  
- `templ`은 `.templ` 파일을 Go 코드로 변환하여 빌드 시 검증 가능.
- `Hello(title, body).Render()`와 같은 함수 호출 방식으로 템플릿을 사용.
- 런타임 오류가 아니라 **컴파일 타임**에 검증 가능.
- 타입 안정성이 보장됨.

---

### 2) `text/template` 기반 HTML 파일을 `templ` 파일로 변환할 때의 문법적 차이

1. **변수 표기법 차이**  
   - `text/template`: `{{ .Variable }}`
   - `templ`: `{Variable}`

```html
<!-- text/template -->
<h1>{{ .Title }}</h1>
<p>{{ .Body }}</p>
```

```templ
<!-- templ -->
Hello(title string, body string) {
	<h1>{title}</h1>
	<p>{body}</p>
}
```

2. **반복문 (`range`) 변환**  
   - `text/template`: `{{ range .Items }}`
   - `templ`: `for item in items {}`

```html
<!-- text/template -->
<ul>
{{ range .Items }}
	<li>{{ . }}</li>
{{ end }}
</ul>
```

```templ
<!-- templ -->
List(items []string) {
	<ul>
		for item in items {
			<li>{item}</li>
		}
	</ul>
}
```

3. **조건문 (`if`) 변환**  
   - `text/template`: `{{ if .Condition }} ... {{ end }}`
   - `templ`: `if condition {}`

```html
<!-- text/template -->
{{ if .ShowMessage }}
	<p>Hello, {{ .User }}!</p>
{{ end }}
```

```templ
<!-- templ -->
Greeting(user string, showMessage bool) {
	if showMessage {
		<p>Hello, {user}!</p>
	}
}
```

4. **함수 호출 방식**  
   - `text/template`은 실행 시 `Execute()` 호출.
   - `templ`은 일반 Go 함수처럼 호출.

```go
// text/template 방식
tmpl.Execute(w, data)
```

```go
// templ 방식
templates.Hello(w, "Hello", "Welcome to templ!")
```

---

### 요약
| 비교 항목           | `text/template` | `templ` |
|------------------|------------------|---------|
| 템플릿 파일 | `.html` | `.templ` |
| 데이터 삽입 | `{{ .Variable }}` | `{Variable}` |
| 반복문 | `{{ range .Items }}` | `for item in items {}` |
| 조건문 | `{{ if .Condition }}` | `if condition {}` |
| 함수 호출 | `tmpl.Execute(w, data)` | `templates.Component(w, args...)` |
| 타입 안정성 | X (런타임 오류) | O (컴파일 시 검증 가능) |

즉, `templ`로 변환하면 더 타입 안정적이고, 함수형 접근 방식으로 `Execute()` 없이 직접 렌더링할 수 있습니다. 

템플릿을 변환하려면 위 문법 차이를 참고하여 기존 `text/template`을 `.templ` 형식으로 다시 작성하면 됩니다! 🚀

### `text/template` vs `templ`에서 템플릿 내부에서 다른 템플릿(컴포넌트)을 호출하는 방법 비교  

템플릿을 재사용하거나, 부모-자식 관계로 구성할 때 사용하는 방식이 `text/template`과 `templ`에서 다릅니다.

---

## 1. `text/template`에서 다른 템플릿을 호출하는 방법  

`text/template`에서는 `define`과 `template` 키워드를 사용하여 다른 템플릿을 포함합니다.  

### 예제: `header.html`, `body.html`, `layout.html`  

**(1) `header.html` - 헤더 템플릿**  
```html
{{ define "header" }}
<header>
    <h1>{{ .Title }}</h1>
</header>
{{ end }}
```

**(2) `body.html` - 본문 템플릿**  
```html
{{ define "body" }}
<main>
    <p>{{ .Content }}</p>
</main>
{{ end }}
```

**(3) `layout.html` - 부모 템플릿에서 자식 템플릿 포함**  
```html
{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    {{ template "header" . }}
    {{ template "body" . }}
</body>
</html>
{{ end }}
```

### 사용 방법: `template.ParseFiles()`로 여러 템플릿을 불러오고 실행  
```go
package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Title   string
	Content string
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("layout.html", "header.html", "body.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := PageData{
		Title:   "My Page",
		Content: "Welcome to my website!",
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

**특징**  
- `{{ define "template_name" }}` 으로 템플릿을 정의하고 `{{ template "name" . }}`로 포함.  
- `ParseFiles()`에서 여러 템플릿을 로드해야 함.  
- 런타임에서 템플릿 이름을 찾아서 실행.  
- 문법 오류를 컴파일 타임이 아닌 실행 중에 발견할 수 있음.  

---

## 2. `templ`에서 다른 템플릿을 호출하는 방법  

`templ`에서는 Go 함수처럼 다른 템플릿(컴포넌트)을 호출할 수 있습니다.  

### 예제: `Header.templ`, `Body.templ`, `Layout.templ`  

**(1) `Header.templ` - 헤더 컴포넌트**  
```templ
@use templ

Header(title string) {
    <header>
        <h1>{title}</h1>
    </header>
}
```

**(2) `Body.templ` - 본문 컴포넌트**  
```templ
@use templ

Body(content string) {
    <main>
        <p>{content}</p>
    </main>
}
```

**(3) `Layout.templ` - 부모 템플릿에서 자식 컴포넌트 호출**  
```templ
@use templ
@use "./Header.templ"
@use "./Body.templ"

Layout(title string, content string) {
    <!DOCTYPE html>
    <html>
    <head>
        <title>{title}</title>
    </head>
    <body>
        Header(title)
        Body(content)
    </body>
    </html>
}
```

### 사용 방법: `templ`은 컴파일된 Go 함수를 직접 호출  
```go
package main

import (
	"net/http"
	"example.com/templates" // go generate 후 생성된 코드가 포함된 패키지
)

func handler(w http.ResponseWriter, r *http.Request) {
	templates.Layout(w, "My Page", "Welcome to my website!")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

**특징**  
- `templ`에서는 `Header(title)`처럼 함수 호출 방식으로 재사용 가능.  
- `@use "./Header.templ"`을 통해 다른 `.templ` 파일을 import.  
- 모든 템플릿이 **Go 코드로 변환되므로 컴파일 타임에 오류를 감지**.  
- `Execute()` 없이 컴포넌트 함수를 바로 실행 가능.  

---

## 3. `text/template` vs `templ` 비교  

| 비교 항목 | `text/template` | `templ` |
|-----------|----------------|---------|
| 템플릿 호출 방식 | `{{ template "name" . }}` | `ComponentName(args...)` |
| 실행 방식 | `ExecuteTemplate(w, "layout", data)` | `Layout(w, title, content)` |
| 데이터 전달 방식 | `.`(context object) | 명시적 함수 인자 전달 |
| 템플릿 정의 방식 | `{{ define "name" }}` | GO 함수 `ComponentName(args...) {}` |
| 템플릿 포함 방식 | `{{ template "name" . }}` | `ComponentName(args...)` |
| 템플릿 로딩 | `ParseFiles("layout.html", ...)` | Template composition `templ showAll() { @left() }` |
| 기타1 |      | Components as parameters `templ layout(contents templ.Component) {` |
| 기타2 |      | Children can be passed to a component `{ children... }` |
---

### **Go 템플릿 시스템 비교: `text/template` vs `templ`**  
이번 비교에서는 `text/template`과 `templ`의 주요 차이점뿐만 아니라,  
- **컴포넌트를 인자로 전달하는 방식 (`templ layout(contents templ.Component) {}`)**
- **`{ children... }` 표현식을 활용하는 방식**  
까지 포함하여 설명하겠습니다.  

---

## 1. **기본 개념 및 주요 차이점**  

| 비교 항목 | `text/template` | `templ` |
|-----------|----------------|---------|
| 템플릿 정의 | `{{ define "name" }}` | `ComponentName(args...) {}` |
| 템플릿 호출 방식 | `{{ template "name" . }}` | `ComponentName(args...)` |
| 데이터 전달 방식 | `.`(context object) | 명시적 함수 인자 전달 |
| 템플릿 로딩 | `ParseFiles("layout.html", ...)` | `@use "Component.templ"`로 import |
| 실행 방식 | `ExecuteTemplate(w, "layout", data)` | `Layout(w, title, content)` |
| **컴포넌트 지원** | 없음 (partial 개념) | 명시적 컴포넌트 지원 |
| **컴포넌트 인자 전달** | `{{ template "name" . }}` | `Layout(contents)` |
| **Children slot 지원** | 없음 | `{ children... }` 사용 가능 |
| 타입 검증 | 런타임에서 오류 탐지 | 컴파일 타임에서 검증 가능 |

---

## 2. **컴포넌트 사용 방식 비교**  

### **2.1 `text/template`에서의 컴포넌트 개념**  

`text/template`에서는 컴포넌트를 직접 정의하는 기능이 없기 때문에 **템플릿을 다른 템플릿에서 호출하는 방식**으로 처리해야 합니다.

#### **📌 예제: `text/template`에서의 레이아웃 + 본문 구조**
##### **`layout.html`** (부모 템플릿)
```html
{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <title>{{ .Title }}</title>
</head>
<body>
    {{ template "header" . }}
    {{ template "body" . }}
</body>
</html>
{{ end }}
```

##### **`header.html`**  
```html
{{ define "header" }}
<header>
    <h1>{{ .Title }}</h1>
</header>
{{ end }}
```

##### **`body.html`**  
```html
{{ define "body" }}
<main>
    <p>{{ .Content }}</p>
</main>
{{ end }}
```

##### **Go 코드에서 템플릿 실행**
```go
tmpl, err := template.ParseFiles("layout.html", "header.html", "body.html")
tmpl.ExecuteTemplate(w, "layout", data)
```

🔹 **한계점:**  
- `{{ template "name" . }}` 형식으로만 다른 템플릿을 포함할 수 있음.  
- **컴포넌트 재사용이 제한적**이며, **동적인 컴포넌트 전달이 불가능**함.  

---

### **2.2 `templ`에서의 컴포넌트 개념**  

`templ`에서는 **컴포넌트를 인자로 전달하거나, `{ children... }`을 사용하여 유연한 레이아웃을 만들 수 있음.**  

#### **📌 예제: `templ`에서의 레이아웃 + 본문 구조**
##### **1️⃣ `Layout.templ` (컴포넌트를 인자로 받는 레이아웃)**
```templ
@use templ

Layout(title string, contents templ.Component) {
    <!DOCTYPE html>
    <html>
    <head>
        <title>{title}</title>
    </head>
    <body>
        {contents}
    </body>
    </html>
}
```

##### **2️⃣ `Header.templ` (독립적인 컴포넌트)**
```templ
@use templ

Header(title string) {
    <header>
        <h1>{title}</h1>
    </header>
}
```

##### **3️⃣ `Body.templ` (독립적인 컴포넌트)**
```templ
@use templ

Body(content string) {
    <main>
        <p>{content}</p>
    </main>
}
```

##### **4️⃣ `Page.templ` (레이아웃을 활용한 페이지)**
```templ
@use templ
@use "./Layout.templ"
@use "./Header.templ"
@use "./Body.templ"

Page(title string, content string) {
    Layout(title, templ.Fragment(
        Header(title),
        Body(content)
    ))
}
```

##### **Go 코드에서 실행**
```go
templates.Page(w, "Welcome", "This is the homepage.")
```

🔹 **장점:**  
- `Layout`이 **`templ.Component` 타입을 인자로 받아** 유연한 구성 가능  
- `templ.Fragment()`를 사용해 여러 컴포넌트를 하나의 요소로 전달 가능  
- 기존의 `text/template`보다 **컴포넌트 기반 아키텍처가 강력함**  

---

## 3. **Children Slot (`{ children... }`) 활용 비교**  

### **3.1 `text/template`에서의 한계**
`text/template`에서는 `{ children... }`과 같은 개념이 없기 때문에, **부모-자식 관계를 직접 구성하는 것이 어렵습니다.**  
즉, **부모 템플릿이 자식 템플릿을 직접 포함해야 하며, 동적으로 자식 컴포넌트를 교체할 수 없습니다.**  

---

### **3.2 `templ`에서 `{ children... }` 활용**  
`templ`에서는 `{ children... }`을 사용하면 **슬롯을 동적으로 채울 수 있는 컴포넌트를 만들 수 있습니다.**  

#### **📌 예제: `templ`에서 `{ children... }` 활용**  
##### **1️⃣ `LayoutWithSlot.templ` (Children 지원)**
```templ
@use templ

LayoutWithSlot(title string) {
    <!DOCTYPE html>
    <html>
    <head>
        <title>{title}</title>
    </head>
    <body>
        { children... }
    </body>
    </html>
}
```

##### **2️⃣ `PageWithSlot.templ` (레이아웃에 동적 콘텐츠 삽입)**
```templ
@use templ
@use "./LayoutWithSlot.templ"
@use "./Header.templ"
@use "./Body.templ"

PageWithSlot(title string, content string) {
    LayoutWithSlot(title) {
        Header(title)
        Body(content)
    }
}
```

##### **Go 코드에서 실행**
```go
templates.PageWithSlot(w, "Dynamic Page", "This content is inside a slot.")
```

🔹 **장점:**  
- `{ children... }`을 사용하면 **부모-자식 관계를 동적으로 변경 가능**  
- `LayoutWithSlot`이 특정 컴포넌트에 의존하지 않고, **유연하게 활용 가능**  
- 기존 `text/template`에서는 불가능한 **동적 슬롯 개념이 도입됨**  

---

## 4. **최종 비교 정리**  

| 비교 항목 | `text/template` | `templ` |
|-----------|----------------|---------|
| 컴포넌트 호출 방식 | `{{ template "name" . }}` | `ComponentName(args...)` |
| 컴포넌트 인자 전달 | 불가능 | `Layout(contents templ.Component)` 가능 |
| 부모-자식 관계 | 부모가 자식을 직접 호출해야 함 | `{ children... }`을 통해 동적 삽입 가능 |
| 컴포넌트 재사용성 | 제한적 | 매우 유연함 |
| 타입 검증 | 없음 (런타임 오류 가능) | 있음 (컴파일 타임 검증) |

---

## 🎯 **결론**
- `templ`은 **컴포넌트 기반 UI를 작성하기에 훨씬 강력한 기능**을 제공  
- `text/template`은 **정적인 구조에서만 유용**하며, **컴포넌트 개념이 부족**  
- `{ children... }`을 활용하면 **부모-자식 관계를 동적으로 변경할 수 있어** 더 유연한 레이아웃 구현 가능  

➡️ **템플릿 시스템을 현대적인 방식으로 사용하려면 `templ`이 훨씬 강력한 선택지가 됨! 🚀**