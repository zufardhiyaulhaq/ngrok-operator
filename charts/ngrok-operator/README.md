# ngrok-operator

ngrok-operator for managing ngrok lifecycle

![Version: 1.0.0](https://img.shields.io/badge/Version-1.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.0.0](https://img.shields.io/badge/AppVersion-1.0.0-informational?style=flat-square) [![made with Go](https://img.shields.io/badge/made%20with-Go-brightgreen)](http://golang.org) [![Github master branch build](https://img.shields.io/github/workflow/status/zufardhiyaulhaq/ngrok-operator/Master)](https://github.com/zufardhiyaulhaq/ngrok-operator/actions/workflows/master.yml) [![GitHub issues](https://img.shields.io/github/issues/zufardhiyaulhaq/ngrok-operator)](https://github.com/zufardhiyaulhaq/ngrok-operator/issues) [![GitHub pull requests](https://img.shields.io/github/issues-pr/zufardhiyaulhaq/ngrok-operator)](https://github.com/zufardhiyaulhaq/ngrok-operator/pulls)

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm repo add zufardhiyaulhaq https://charts.zufardhiyaulhaq.com/
helm install my-release zufardhiyaulhaq/ngrok-operator --values values.yaml
```

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| operator.image | string | `"zufardhiyaulhaq/ngrok-operator"` |  |
| operator.replica | int | `1` |  |
| operator.tag | string | `"v1.0.0"` |  |

