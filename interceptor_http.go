package interceptor

import (
	"net/http"
)

type HttpInterceptor func(http.ResponseWriter, *http.Request, http.HandlerFunc) http.HandlerFunc

func ChainHttpInterceptor(interceptors ...HttpInterceptor) HttpInterceptor {
	n := len(interceptors)
	return func(rw http.ResponseWriter, req *http.Request, handler http.HandlerFunc) http.HandlerFunc {
		chainer := func(currInterceptor HttpInterceptor, currHandler http.HandlerFunc) http.HandlerFunc {
			return currInterceptor(rw, req, currHandler)
		}

		currHandler := handler
		for i := n - 1; i >= 0; i-- {
			currHandler = chainer(interceptors[i], currHandler)
		}
		currHandler(rw, req)
		return currHandler
	}
}

func HttpInterceptorWarp(interceptors ...HttpInterceptor) func(http.HandlerFunc) http.HandlerFunc {
	chainInterceptor := ChainHttpInterceptor(interceptors...)
	return func(handler http.HandlerFunc) http.HandlerFunc {
		return func(rw http.ResponseWriter, req *http.Request) {
			chainInterceptor(rw, req, handler)
		}
	}
}

func HttpInterceptorWarpHandleFunc(handler http.HandlerFunc, interceptors ...HttpInterceptor) http.HandlerFunc {
	chainInterceptor := ChainHttpInterceptor(interceptors...)
	return func(rw http.ResponseWriter, req *http.Request) {
		chainInterceptor(rw, req, handler)
	}
}
