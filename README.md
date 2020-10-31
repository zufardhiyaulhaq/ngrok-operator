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

### Installation
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

## Contributors ✨
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-1-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="http://zufardhiyaulhaq.com"><img src="https://avatars3.githubusercontent.com/u/11990726?v=4" width="100px;" alt=""/><br /><sub><b>Zufar Dhiyaulhaq</b></sub></a><br /><a href="#infra-zufardhiyaulhaq" title="Infrastructure (Hosting, Build-Tools, etc)">🚇</a> <a href="https://github.com/zufardhiyaulhaq/ngrok-operator/commits?author=zufardhiyaulhaq" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome, please check [CONTRIBUTING.md](https://github.com/zufardhiyaulhaq/ngrok-operator/blob/master/.github/CONTRIBUTING.md)!

## Changes
For changes, see the [CHANGELOG.md](CHANGELOG.md).

## License
This program is free software: you can redistribute it and/or modify it under the terms of the [MIT license](LICENSE)
