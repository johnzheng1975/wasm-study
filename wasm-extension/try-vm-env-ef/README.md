 2181  tinygo build -o main.wasm -scheduler=none -target=wasi main.go
 2186  tinygo build -o  second-extension.wasm -scheduler=none -target=wasi   second-extension.go
 2188  func-e run -c envoy.yaml
 

