# ngrok-operator

Ngrok operator provide developer easy access to private Kubernetes cluster for testing purpose via ngrok. Automate the creation of ngrok tunnel via CRD!

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.0](https://img.shields.io/badge/AppVersion-1.0.0-informational?style=flat-square) [![made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen)](http://golang.org) [![Github master branch build](https://img.shields.io/github/workflow/status/zufardhiyaulhaq/ngrok-operator/Master)](https://github.com/zufardhiyaulhaq/ngrok-operator/actions/workflows/master.yml) [![GitHub issues](https://img.shields.io/github/issues/zufardhiyaulhaq/ngrok-operator)](https://github.com/zufardhiyaulhaq/ngrok-operator/issues) [![GitHub pull requests](https://img.shields.io/github/issues-pr/zufardhiyaulhaq/ngrok-operator)](https://github.com/zufardhiyaulhaq/ngrok-operator/pulls)

## Installing

To install the chart with the release name `my-release`:

```console
helm repo add zufardhiyaulhaq https://charts.zufardhiyaulhaq.com/
helm install my-release zufardhiyaulhaq/ngrok-operator --values values.yaml
```

## Usage
1. Apply some example
```console
kubectl apply -f examples/nginx
kubectl apply -f examples/helloworld/namespace.yaml
kubectl apply -f examples/helloworld/
```
2. Check ngrok object
```console
kubectl get ngrok --all-namespaces
NAMESPACE    NAME               STATUS    URL
default      nginx-ngrok        created   https://9496e56ed0bc.ngrok.io
default      nginx-ngrok-full   created   https://ngrok.zufardhiyaulhaq.com
helloworld   helloworld-ngrok   created   https://d00ba8cb0b95.ngrok.io
```
3. access the URL
```console
https://d5150f7c3588.ngrok.io
https://ngrok.zufardhiyaulhaq.com
https://fa03f71fbe18.ngrok.io/hello
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| operator.image | string | `"zufardhiyaulhaq/ngrok-operator"` |  |
| operator.replica | int | `1` |  |
| operator.tag | string | `"v1.0.0"` |  |

