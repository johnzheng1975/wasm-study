package main

import (
    "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
    "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
    "github.com/valyala/fastjson"
    "encoding/binary"
)

func main() {
    proxywasm.SetVMContext(&vmContext{})
}

const (
  sharedDataKey                 = "my_key"
  sharedDataInitialValue uint64 = 1
)


type vmContext struct {
    // Embed the default VM context here,
    // so that we don't need to reimplement all the methods.
    types.DefaultVMContext
}

// Override types.VMContext.
func (*vmContext) OnVMStart(vmConfigurationSize int) types.OnVMStartStatus {
  initialValueBuf := make([]byte, 8)
  binary.LittleEndian.PutUint64(initialValueBuf, sharedDataInitialValue) // 
  if err := proxywasm.SetSharedData(sharedDataKey, initialValueBuf, 0); err != nil {
    proxywasm.LogWarnf("error setting shared data on OnVMStart: %v", err)
  }
  proxywasm.LogInfof("[wasm-extension]: Setting initial shared value %v", sharedDataInitialValue)
  return types.OnVMStartStatusOK
}


func (ctx *httpContext) incrementData() (uint64, error) {
  value, cas, err := proxywasm.GetSharedData(sharedDataKey) // 
  if err != nil {
    proxywasm.LogWarnf("error getting shared data: %v", err)
    return 0, err
  }

  buf := make([]byte, 8)
  ret := binary.LittleEndian.Uint64(value) + 1 // 
  binary.LittleEndian.PutUint64(buf, ret)
  if err := proxywasm.SetSharedData(sharedDataKey, buf, cas); err != nil { // 
    proxywasm.LogWarnf("error setting shared data: %v", err)
    return 0, err
  }
  return ret, err
}


// Override types.DefaultVMContext.

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
    return &pluginContext{contextID: contextID, additionalHeaders: map[string]string{}, helloHeaderCounter: proxywasm.DefineCounterMetric("hello_header_counter")}
}


type pluginContext struct {
    // Embed the default plugin context here,
    // so that we don't need to reimplement all the methods.
    types.DefaultPluginContext
    additionalHeaders map[string]string
    contextID         uint32
    helloHeaderCounter proxywasm.MetricCounter
}


// Override types.DefaultPluginContext.
func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
  return &httpContext{contextID: contextID, additionalHeaders: ctx.additionalHeaders, helloHeaderCounter: ctx.helloHeaderCounter}
}


type httpContext struct {
    // Embed the default http context here,
    // so that we don't need to reimplement all the methods.
    types.DefaultHttpContext
    contextID uint32
    additionalHeaders map[string]string
    helloHeaderCounter proxywasm.MetricCounter
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
  data, err := proxywasm.GetPluginConfiguration() // 
  if err != nil {
    proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
  }

  var p fastjson.Parser
  v, err := p.ParseBytes(data)
  if err != nil {
    proxywasm.LogCriticalf("error parsing configuration: %v", err)
  }

  obj, err := v.Object()
  if err != nil {
    proxywasm.LogCriticalf("error getting object from json value: %v", err)
  }

  obj.Visit(func(k []byte, v *fastjson.Value) {
    ctx.additionalHeaders[string(k)] = string(v.GetStringBytes()) // 
  })

  return types.OnPluginStartStatusOK
}



func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
  proxywasm.LogInfo("OnHttpRequestHeaders")

  _, err := proxywasm.GetHttpRequestHeader("hello") // 
  if err != nil {
    // Ignore if header is not set
    return types.ActionContinue
  }

  ctx.helloHeaderCounter.Increment(1) // 
  proxywasm.LogInfo("hello_header_counter incremented")
  return types.ActionContinue
}

