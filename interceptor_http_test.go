package interceptor

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func FooHttpInterceptor(rw http.ResponseWriter, req *http.Request, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Println("FooHttpInterceptor", req.RequestURI)
		fmt.Fprint(rw, "foo")
		handler(rw, req)
	}
}

func BarHttpInterceptor(rw http.ResponseWriter, req *http.Request, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Println("BarHttpInterceptor", req.RequestURI)
		fmt.Fprint(rw, "bar")
		handler(rw, req)
	}
}

func ForbiddenHttpInterceptor(rw http.ResponseWriter, req *http.Request, handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Println("ForbiddenHttpInterceptor", req.RequestURI)
		rw.WriteHeader(http.StatusForbidden)
		fmt.Fprint(rw, "forbidden")
		return
		handler(rw, req)
	}
}

func TestChainHttpInterceptor(t *testing.T) {
	helloHandler := func(rw http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rw, `helloworld!`)
	}

	testcases := []struct {
		warp func(http.HandlerFunc) http.HandlerFunc
		want string
	}{
		{
			warp: HttpInterceptorWarp(FooHttpInterceptor, BarHttpInterceptor),
			want: "foobarhelloworld!",
		},
		{
			// 设置Header 需要在write之前，foo会写body，故把forbidden放在前面
			warp: HttpInterceptorWarp(ForbiddenHttpInterceptor, FooHttpInterceptor),
			want: "forbidden",
		},
	}

	for _, it := range testcases {
		ts := httptest.NewServer(it.warp(helloHandler))
		resp, err := http.Get(ts.URL)
		if err != nil {
			log.Fatal(err)
		}
		greeting, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("err: %+v", err)
		}
		ts.Close()
		resp.Body.Close()

		got := string(greeting)
		fmt.Printf("resp: %s\n", got)
		if got != it.want {
			t.Errorf("got: %s want: %s", got, it.want)
		}
	}
}
