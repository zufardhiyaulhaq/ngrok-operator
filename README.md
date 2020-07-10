# ngrok-operator
Ngrok operator provide developer easy access to private Kubernetes cluster for testing purpose via ngrok. Automate the creation of ngrok tunnel via CRDs!

### Feature
- [x] basic ngrok feature
- [x] support HTTP
- [ ] support TCP
- [ ] support Costum Configuration

### Developing ngrok-operator
This operator build based on [operator-sdk](https://sdk.operatorframework.io/docs/install-operator-sdk/). To build this operator, you need [operator-sdk](https://sdk.operatorframework.io/docs/install-operator-sdk/).

#### Installing ngrok-operator via helm
Please read README.md in charts folder for more information.
```
helm install ./charts --name-template ngrok-operator
```

to insatall without crds
```
--skip-crds
```

to upgrade
```
helm upgrade ngrok-operator ./charts
```

### Example
- Deploy ngrok via Helm
- Apply some example
```
kubectl apply -f examples/nginx
kubectl apply -f examples/helloworld/namespace.yaml
kubectl apply -f examples/helloworld/
```
- Check ngrok object
```
kubectl get ngrok --all-namespaces
NAMESPACE    NAME               STATUS    URL
default      nginx-ngrok        created   https://d5150f7c3588.ngrok.io
helloworld   helloworld-ngrok   created   https://fa03f71fbe18.ngrok.io
```
- access the URL
```
https://d5150f7c3588.ngrok.io
https://fa03f71fbe18.ngrok.io/hello
```
