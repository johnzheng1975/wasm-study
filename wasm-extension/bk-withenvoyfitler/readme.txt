https://raw.githubusercontent.com/tetratelabs/proxy-wasm-go-sdk/main/examples/vm_plugin_configuration/main.go
https://github.com/johnzheng1975/proxy-wasm-go-sdk
https://github.com/tetratelabs/proxy-wasm-go-sdk/blob/main/examples/vm_plugin_configuration/main.go



k exec -ti httpbin-5dcc45899f-slfgz  -c istio-proxy -- bash
curl localhost:15000/stats/prometheus | grep hello

root@nginx:/# for i in {1..100};do  curl  -H "hello: something" httpbin:8000  ; done;

