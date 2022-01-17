<p align="center"><img src="https://raw.githubusercontent.com/mms-gianni/klustair-frontend/master/docs/img/klustair.png" width="200"></p>

# <a href='https://github.com/mms-gianni/klustair'>KlustAIR Client</a>
The Klustair client searches your Kubernetes namespaces for the used images and scans them with Trivy.

### Related Klustair projects: 
- <a href="https://github.com/mms-gianni/klustair-frontend">Klustair Frontend</a> to view the scanner results
- <a href="https://github.com/mms-gianni/klustair-helm">Klustair Helm charts</a> to run Klustair Cronjob, API and Frontend

### Related opensource projects
- <a href="https://github.com/aquasecurity/trivy">trivy</a> A Simple and Comprehensive Vulnerability Scanner for Containers and other Artifacts
- <a href="https://github.com/Shopify/kubeaudit">kubeaudit</a> kubeaudit helps you audit your Kubernetes clusters against common security controls

## Usage
```
klustair [global options]

optional arguments:
   --verbose, -V                          increase output verbosity (default: false) [$KLUSTAIR_VERBOSE]
   --debug, -d                            debug mode (default: false) [$KLUSTAIR_DEBUG]
   --namespaces value, -n value           Coma separated whitelist of Namespaces to check [$KLUSTAIR_NAMESPACES]
   --namespacesblacklist value, -N value  Coma separated whitelist of Namespaces to check [$KLUSTAIR_NAMESPACESBLACKLIST]
   --kubeaudit value, -k value            Coma separated list of audits to run. (disable: "none") [$KLUSTAIR_KUBEAUDIT]
   --trivy, -t                            Run Trivy vulnerability checks (default: false) [$KLUSTAIR_TRIVY]
   --label value, -l value                A optional title for your run [$KLUSTAIR_NAMESPACESBLACKLIST]
   --repocredentialspath value, -c value  Path to repo credentials for trivy [$KLUSTAIR_REPOCREDENTIALSPATH]
   --limitdate value, --ld value          Remove reports older than X days (default: 0) [$KLUSTAIR_LIMITDATE]
   --limitnr value, --ln value            Keep only X reports (default: 0) [$KLUSTAIR_LIMITNR]
   --configkey value, -C value            Load remote configuration from frontend [$KLUSTAIR_CONFIGKEY]
   --apihost value, -H value              Remote API-host address [example: https://localhost:8443] [$KLUSTAIR_APIHOST]
   --apitoken value, -T value             API Access Token from Klustair Frontend [$KLUSTAIR_APITOKEN]
   --help, -h                             show help (default: false)
   --version, -v                          print the version (default: false)
```

## ENV vars (not set by commandline)
```
export TRIVY_USERNAME=....
export TRIVY_PASSWORD=....
export TRIVY_REGISTRY_TOKEN=....
export TRIVY_INSECURE=false
export TRIVY_NON_SSL=false
``

## Installation
```
go get -v github.com/klustair/klustair-cli
```

## develop
```
git clone git@github.com:klustair/klustair-cli.git
cd klustair-cli
go run cmd/klustair/main.go
```

## build
```
go build -o bin/klustair-cli cmd/klustair/main.go
```

## FAQ
Why is the klustair client so big (~80MB)? 
 - it contains the trivy binary(~32MB) and the kubeaudit binary (~30MB).

 