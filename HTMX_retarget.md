## **templ \+ HTMX 를 사용한 GO web server 에러 처리 (HX-Retarget \+ Redirect 사용) 및 Fragment Swap 코드**

요청하신 templ과 HTMX를 사용한 Go 웹 서버에서 에러 발생 시 HX-Retarget 기술과 리다이렉트를 조합하여 특정 영역만 업데이트하는 전체 코드를 제공해 드리겠습니다. 이번에는 someActionHandler에서 HTMX 요청 여부를 확인하고, HTMX 요청인 경우 HX-Retarget 헤더를 설정 **후** /error URL로 리다이렉트하며, errorHandler에서는 HX-Retarget에 의해 지정된 영역에 에러 메시지 Fragment를 반환하는 방식입니다.

**핵심 목표:**

* **에러 핸들링:** 서버 내부 처리 중 에러 발생 시, HTMX 요청에 대해 HX-Retarget 헤더 설정 후 /error URL로 리다이렉트  
* **HX-Retarget 사용:** HX-Retarget 헤더를 사용하여 HTMX 응답의 타겟을 명시적으로 지정 (리다이렉트 전에 설정)  
* **errorHandler Fragment Swap:** /error 핸들러는 HX-Retarget 헤더에 명시된 영역 (id=container) 에 에러 메시지 Fragment HTML 반환  
* **Templ 사용:** HTML 템플릿 및 HTMX Fragment 생성을 위해 templ 패키지 활용

**전체 코드 구조:**

1. **main.go:** Go 웹 서버의 메인 로직 (라우팅, 핸들러, 서버 시작)  
2. **templates 폴더:** templ 템플릿 파일 (index.templ, error.templ)

**1\. main.go 코드:**

Go

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"go-htmx-retarget-redirect-error/templates" // templates 패키지 import (go.mod에 따라 경로 수정)
)

// AppError: 애플리케이션 에러를 나타내는 사용자 정의 타입
type AppError struct {
	Message string
	Code    int
}

// Error: AppError 가 error 인터페이스를 구현하도록 함
func (e AppError) Error() string {
	return fmt.Sprintf("AppError: %d - %s", e.Code, e.Message)
}

// someActionHandler: 에러 발생 가능성이 있는 액션 핸들러 (예시)
func someActionHandler(w http.ResponseWriter, r *http.Request) {
	// ... (내부 처리 로직) ...

	// 예시: id 파라미터를 숫자로 변환 시 에러 발생 가능
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// id 파라미터가 없는 경우, 에러 발생 (예시)
		appErr := AppError{Message: "ID 파라미터가 필요합니다.", Code: http.StatusBadRequest}
		log.Println("[someActionHandler] Error: ", appErr)

		// HTMX 요청인지 확인 (HX-Request 헤더 체크)
		if r.Header.Get("HX-Request") == "true" {
			// HTMX 요청인 경우: HX-Retarget 설정 후 에러 페이지로 리다이렉트
			w.Header().Set("HX-Retarget", "#container") // HX-Retarget 헤더 설정: 타겟 엘리먼트 지정
			redirectToError(w, r, appErr.Message)      // 에러 페이지로 리다이렉트 (errorHandler에서 Fragment 응답)
			return
		} else {
			// 일반적인 HTTP 요청 (HTMX 요청이 아닌 경우): 에러 페이지로 리다이렉트 (기존 방식 유지)
			redirectToError(w, r, appErr.Message)
			return
		}
	}

	_, err := strconv.Atoi(idStr)
	if err != nil {
		// id 파라미터가 숫자가 아닌 경우, 에러 발생 (예시)
		appErr := AppError{Message: "유효하지 않은 ID 입니다.", Code: http.StatusBadRequest}
		log.Println("[someActionHandler] Error: ", appErr)

		// HTMX 요청인지 확인 (HX-Request 헤더 체크)
		if r.Header.Get("HX-Request") == "true" {
			// HTMX 요청인 경우: HX-Retarget 설정 후 에러 페이지로 리다이렉트
			w.Header().Set("HX-Retarget", "#container") // HX-Retarget 헤더 설정: 타겟 엘리먼트 지정
			redirectToError(w, r, appErr.Message)      // 에러 페이지로 리다이렉트 (errorHandler에서 Fragment 응답)
			return
		} else {
			// 일반적인 HTTP 요청 (HTMX 요청이 아닌 경우): 에러 페이지로 리다이렉트 (기존 방식 유지)
			redirectToError(w, r, appErr.Message)
			return
		}
	}

	// ... (정상 처리 로직) ...
	fmt.Fprintln(w, "액션 성공!") // 성공 메시지 (HTMX swap으로 대체 가능)
}

// errorHandler: 에러 페이지 핸들러 (Fragment Swap 응답)
func errorHandler(w http.ResponseWriter, r *http.Request) {
	errorMessage := r.URL.Query().Get("msg") // 리다이렉트 시 전달된 에러 메시지 획득

	// 에러 메시지를 담아 error.templ 렌더링 (HTMX Fragment)
	templates.Error(errorMessage).Render(r.Context(), w) // error.templ 렌더링 (Fragment HTML)
}

// redirectToError: 에러 페이지로 리다이렉트 (에러 메시지 전달)
func redirectToError(w http.ResponseWriter, r *http.Request, msg string) {
	errorURL := "/error?msg=" + msg
	http.Redirect(w, r, errorURL, http.StatusFound) // 302 Found 리다이렉트
}

// indexHandler: 초기 페이지 핸들러 (index.templ 렌더링)
func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.Index().Render(r.Context(), w)
}

func main() {
	http.HandleFunc("/", indexHandler)           // 초기 페이지 핸들러 등록
	http.HandleFunc("/action", someActionHandler) // 액션 핸들러 등록
	http.HandleFunc("/error", errorHandler)     // 에러 핸들러 등록

	fmt.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**2\. templates/index.templ 코드:**

Code snippet

```html
<!DOCTYPE html>
<html>
<head>
	<title>HTMX Error Handling (HX-Retarget + Redirect)</title>
	<script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
	<h1>HTMX 에러 핸들링 데모 (HX-Retarget + Redirect)</h1>

	<div id="container">
		<p>여기에 내용이 표시됩니다.</p>
	</div>

	<button hx-get="/action?id=123" hx-target="#container" hx-swap="innerHTML">정상 액션 실행</button>
	<button hx-get="/action" hx-target="#container" hx-swap="innerHTML">에러 발생 액션 실행 (ID 파라미터 없음)</button>
	<button hx-get="/action?id=abc" hx-target="#container" hx-swap="innerHTML">에러 발생 액션 실행 (잘못된 ID)</button>

</body>
</html>
```

**3\. templates/error.templ 코드:**

Code snippet

```html
{# templates/error.templ #}
@template Error(msg string)

<div id="container">
	<p style="color: red;"><strong>에러 발생 (Redirect):</strong> { msg }</p>
</div>
```

**코드 설명:**

* **someActionHandler 수정:**  
  * **HTMX 요청 감지:** r.Header.Get("HX-Request") \== "true" 조건으로 HTMX 요청인지 확인합니다.  
  * **HTMX 에러 처리 (HX-Retarget \+ Redirect):** HTMX 요청인 경우, w.Header().Set("HX-Retarget", "\#container") 를 사용하여 HX-Retarget 헤더를 설정합니다. **핵심은 리다이렉트 전에 HX-Retarget 헤더를 설정한다는 점입니다.** 이후 redirectToError 함수를 호출하여 /error URL로 리다이렉트하고, 에러 메시지를 쿼리 파라미터로 전달합니다.  
  * **비 HTMX 에러 처리:** HTMX 요청이 아닌 경우, 기존처럼 redirectToError 함수를 사용하여 에러 페이지로 리다이렉트합니다.  
* **errorHandler:** /error 엔드포인트에 대한 핸들러입니다.  
  * redirectToError 함수에서 리다이렉트될 때 쿼리 파라미터로 전달된 에러 메시지를 r.URL.Query().Get("msg") 를 통해 획득합니다.  
  * templates.Error(errorMessage) 를 호출하여 error.templ 템플릿에 에러 메시지를 전달하고 렌더링합니다.  
  * **error.templ 은 \<div id="container"\> 만을 포함하는 HTMX Fragment HTML을 생성합니다.** HTMX는 리다이렉션 응답을 받으면, HX-Retarget 헤더에 지정된 \#container 를 찾아서, 응답으로 받은 Fragment HTML (\<div id="container" ...\> ... \</div\>) 로 innerHTML swap 방식으로 교체합니다.  
* **redirectToError:** 에러 발생 시 /error URL로 리다이렉트하는 유틸리티 함수입니다. 에러 메시지를 쿼리 파라미터에 담아 전달합니다.  
* **indexHandler 및 index.templ**: 초기 페이지 핸들러 및 템플릿은 이전 예시와 동일합니다.  
* **templates/error.templ**: 에러 메시지를 표시하는 HTMX Fragment 템플릿입니다. \<div id="container"\> 를 최상위 엘리먼트로 사용합니다.

**핵심 포인트:**

* **HX-Retarget 헤더 \+ Redirect 조합**: HX-Retarget 헤더를 설정한 상태에서 리다이렉션을 사용하면, 페이지 전체 리다이렉션 후 HTMX가 HX-Retarget 에 지정된 영역만 응답 내용으로 업데이트하는 효과를 낼 수 있습니다.  
* **errorHandler 의 Fragment 응답**: errorHandler는 /error URL 요청에 대해 HX-Retarget 에 의해 타겟팅된 영역을 위한 Fragment HTML 응답을 반환합니다.  
* **URL 변경**: URL이 /error?msg=... 로 변경되므로, 사용자에게 에러 상황 발생 및 에러 페이지로 이동했음을 명확히 인지시킬 수 있습니다.
