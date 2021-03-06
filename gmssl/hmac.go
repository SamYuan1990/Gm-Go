/*
 * Copyright 2020 The Hyperledger-TWGC Project Authors. All Rights Reserved.
 *
 * Licensed under the Apache License 2.0 (the "License").  You may not use
 * this file except in compliance with the License.  You can obtain a copy
 * in the file LICENSE in the source distribution or at
 * https://www.openssl.org/source/license.html
 */

/* +build cgo */

package gmssl

/*
#include <openssl/hmac.h>
#include <openssl/cmac.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

func GetMacLength(name string) (int, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	md := C.EVP_get_digestbyname(cname)
	if md == nil {
		return 0, GetErrors()
	}
	return int(C.EVP_MD_size(md)), nil
}

type HMACContext struct {
	hctx *C.HMAC_CTX
}

func NewHMACContext(name string, key []byte) (
	*HMACContext, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	md := C.EVP_get_digestbyname(cname)
	if md == nil {
		return nil, GetErrors()
	}
	ctx := C.HMAC_CTX_new()
	if ctx == nil {
		return nil, GetErrors()
	}
	ret := &HMACContext{ctx}
	runtime.SetFinalizer(ret, func(ret *HMACContext) {
		C.HMAC_CTX_free(ret.hctx)
	})
	if 1 != C.HMAC_Init_ex(ctx,
		unsafe.Pointer(&key[0]), C.int(len(key)), md, nil) {
		return nil, GetErrors()
	}
	return ret, nil
}

func (ctx *HMACContext) Update(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if 1 != C.HMAC_Update(ctx.hctx,
		(*C.uchar)(unsafe.Pointer(&data[0])), C.size_t(len(data))) {
		return GetErrors()
	}
	return nil
}

func (ctx *HMACContext) Final() ([]byte, error) {
	outbuf := make([]byte, 64)
	outlen := C.uint(len(outbuf))
	if 1 != C.HMAC_Final(ctx.hctx,
		(*C.uchar)(unsafe.Pointer(&outbuf[0])), &outlen) {
		return nil, GetErrors()
	}
	return outbuf[:outlen], nil
}

func (ctx *HMACContext) Reset() error {
	if 1 != C.HMAC_Init_ex(ctx.hctx, nil, 0, nil, nil) {
		return GetErrors()
	}
	return nil
}
