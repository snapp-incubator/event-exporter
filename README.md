# Event-exporter

This tool is used to export k8s events.

## Build

`docker build -t event-exporter .`

## Prometheus Endpoint

Prometheus endpoint uses port `8090` and `/metric`
Metric name is `snappcloud_event_openshift` by default.

## Events

https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/events/event.go
event_reason:
    Pod:
        Scheduled
        Pulling
        Pulled
        Created
        Started
        --
        BackOff
        Unhealthy
        FailedMount
    ReplicationController:
        SuccessfulCreate
    DeploymentConfig:
        DeploymentCreated
    DaemonSet
        FailedCreate
    StatefulSet
        FailedCreate
    HorizontalPodAutoscaler
        FailedGetResourceMetric
    Node
        Rebooted
        NodeNotReady
        HostPortConflict
FailedScheduling
FailedAttachVolume
Failed
Evicted


---
source component:
    default-scheduler
    kubelet
    replication-controller
    deploymentconfig-controller
    a/d controller
