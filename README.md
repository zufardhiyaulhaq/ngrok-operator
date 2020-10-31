<div style="text-align:center">
<img src="./logo.png" width="200">
</div>
Ngrok operator provide developer easy access to private Kubernetes cluster for testing purpose via ngrok. Automate the creation of ngrok tunnel via CRDs!

### Feature
- [x] support HTTP
- [x] support TCP
- [x] support costum configuration
  - [x] custom domain
  - [x] custom TCP address
  - [x] custom region
  - [x] enable/disable inspection
  - [x] support HTTP auth
  - [ ] service for ngrok object (dashboard related)

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
default      nginx-ngrok        created   https://9496e56ed0bc.ngrok.io
default      nginx-ngrok-full   created   https://ngrok.zufardhiyaulhaq.com
helloworld   helloworld-ngrok   created   https://d00ba8cb0b95.ngrok.io
```
- access the URL
```
https://d5150f7c3588.ngrok.io
https://ngrok.zufardhiyaulhaq.com
https://fa03f71fbe18.ngrok.io/hello
```

### Contributors
