# ClusterLendingManager

[![Go Reference][godoc-image]][godoc-link]
[![Coverage Status][cov-image]][cov-link]

ClusterLendingManager is an operator to manage muti-tenant cluster's resources.

Here's a `LendingConfig` example.

```yaml
apiVersion: clusterlendingmanager.dtaniwaki.github.com/v1alpha1
kind: LendingConfig
metadata:
  name: lending-config-sample
spec:
  targets:
  - kind: Deployment
    apiVersion: apps/v1
  - kind: Rollout
    apiVersion: argoproj.io/v1alpha1
  timezone: "Asia/Tokyo"
  schedule:
    default:
      hours:
      - start: "10:00"
        end: "20:00"
    friday:
      hours:
      - start: "10:00"
        end: "17:00" # Happy Friday!
    saturday:
      hours: [] # Ofcourse, no work!
    sunday:
      hours: [] # Ofcourse, no work!
```

## Prerequisites

- [golangci-lint v1.42.1](https://github.com/golangci/golangci-lint)

## Build

Build and load the Docker image to your cluster.

```bash
$ make docker-build

# run a command to load the image to your cluster.
```

If you use a kind cluster, there's a useful shortcut.

```
$ make kind-load
```

## Deployment

Install the CRD to the cluster.

```bash
$ make install
```

Deploy a controller to the cluster.

```bash
$ make deploy
```

## Usage

Now, deploy the samples.

```bash
$ make deploy-samples
```

You will see sample HPA and deployment in the current context, maybe `default` depending on your env. The HPA resource gets updated periodically by the CronHPA.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new [Pull Request](../../pull/new/master)

## Copyright

Copyright (c) 2021 Daisuke Taniwaki. See [LICENSE](LICENSE) for details.


[godoc-image]: https://pkg.go.dev/badge/github.com/dtaniwaki/cluster-lending-manager.svg
[godoc-link]: https://pkg.go.dev/github.com/dtaniwaki/cluster-lending-manager
[cov-image]:   https://coveralls.io/repos/github/dtaniwaki/cluster-lending-manager/badge.svg?branch=main
[cov-link]:    https://coveralls.io/github/dtaniwaki/cluster-lending-manager?branch=main

