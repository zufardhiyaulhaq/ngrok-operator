# ngrok-operator charts
Helm chart for ngrok-operators

### Installing the charts
```
helm repo add zufardhiyaulhaq https://charts.zufardhiyaulhaq.com/
helm install zufardhiyaulhaq/ngrok-operator --name-template ngrok-operator
```

### Configuration

| Parameter | Description | Default |
|-|-| -|
| operator.image | Image for ngrok-operator | zufardhiyaulhaq/ngrok-operator |
| operator.tag | Tag for image ngrok-operator | 0.0.2 |
| operator.pullPolicy | pullPolicy | Always |
| operator.replica | number of replica | 1 |

Specify each parameter using the `--set key=value[,key=value]` argument to helm install.
