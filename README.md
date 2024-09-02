# Event Exporter

A Prometheus exporter for exporting k8s events.

## Build

### Docker

```bash
docker build -t event-exporter .
```

### Binary

```bash
git clone https://github.com/snapp-incubator/event-exporter.git
cd event-exporter
go build
```

## Installation

### Docker

```bash
docker run -p 8080:8080 ghcr.io/snapp-incubator/event-exporter:main
```

### Helm chart

1. Add the Event Exporter Helm Repository:

```bash
helm repo add snapp-cab https://snapp-cab.github.io/event-exporter/charts
helm repo update
```

2. Install with:

```bash
helm install event-exporter snapp-cab/event-exporter
```

### Binary releases

```bash
export VERSION=1.0.0
wget https://github.com/cafebazaar/event-exporter/releases/download/v${VERSION}/event-exporter-${VERSION}.linux-amd64.tar.gz
tar xvzf event-exporter-${VERSION}.linux-amd64.tar.gz event-exporter-${VERSION}.linux-amd64/event-exporter
```

## Metrics and events

|       Metric        | Notes               | Labels                                             |
| :-----------------: | :------------------ | :------------------------------------------------- |
| `event_normal_k8s`  | Normal k8s events.  | `kind`, `namespace`, `reason`, `source_components` |
| `event_warning_k8s` | Warning k8s events. | `kind`, `namespace`, `reason`, `source_components` |

### Event reasons

| Kind                    | Reason                                       |
| ----------------------- | -------------------------------------------- |
| Pod (Normal)            | Scheduled, Pulling, Pulled, Created, Started |
| Pod (Warning)           | BackOff, Unhealthy, FailedMount              |
| ReplicationController   | SuccessfulCreate                             |
| DeploymentConfig        | DeploymentCreated                            |
| DaemonSet               | FailedCreate                                 |
| StatefulSet             | FailedCreate                                 |
| HorizontalPodAutoscaler | FailedGetResourceMetric                      |
| Node                    | Rebooted, NodeNotReady, HostPortConflict     |

## Security

### Reporting security vulnerabilities

If you find a security vulnerability or any security related issues, please DO NOT file a public issue,
instead send your report privately to cloud@snapp.cab.
Security reports are greatly appreciated, and we will publicly thank you for it.

## License

Apache-2.0 License, see [LICENSE](LICENSE).
