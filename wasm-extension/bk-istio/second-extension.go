package main

import (
    "encoding/binary"

    "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
    "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

const (
    sharedDataKey = "my_key"
)

func main() {
    proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
    // Embed the default VM context here,
    // so that we don't need to reimplement all the methods.
    types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
    return &pluginContext{}
}

type pluginContext struct {
    // Embed the default plugin context here,
    // so that we don't need to reimplement all the methods.
    types.DefaultPluginContext
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
    proxywasm.LogInfo("[second-extension]: NewHttpContext")

    value, _, err := proxywasm.GetSharedData(sharedDataKey)
    if err != nil {
        proxywasm.LogWarnf("[second-extension]: error getting shared data on NewHttpContext: %v", err)
    }
    buf := make([]byte, 8)
    ret := binary.LittleEndian.Uint64(value)
    binary.LittleEndian.PutUint64(buf, ret)
    proxywasm.LogInfof("[second-extension]: Reading shared data: %d", ret)

    return &httpContext{contextID: contextID}
}

type httpContext struct {
    // Embed the default http context here,
    // so that we don't need to reimplement all the methods.
    types.DefaultHttpContext
    contextID uint32
}

