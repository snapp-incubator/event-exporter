# openshift-event-exporter

This tool is used to export Openshift events. It effectively runs a watch on
the API server, detecting as granular as possible all changes to the event
objects.

From git@gitlab.snapp.ir:sajad.orouji/openshift-event-exporter.git

## Build

`docker build -t event-exporter .`

## Run
For running event exporter container:

You can run this container for testing on your local system using docker command:

`docker run -d --name event-exporter -e "OPENSHIFT_API_URL=<Your API Url>" -e "OPENSHIFT_TOKEN=<Your Openshift token>" -e "OPENSHIFT_NAMESPACE=<Your project name>" event-exporter`

For using it on your openshift project:
## OKD RESOURCE
File located in `okd` folder should be applied with `oc command`

### RoleBinding
First edit `rbacRole.yaml` file and set the namespace to your namespace name and apply `rbacRole.yaml` file.
The file contains 2 resources, service account, and role binding, this 2 together give you permission to access openshift API.

### Make

1. Set your namespace (i.e. project name) and repository URL in .makerc.dist
2. For repository source secret you can read [Add gitlab secret](http://docs.snappcloud.io/quickstart/ssh-keys.html)
3. Make sure you are logged in via oc CLI.
4. make

Wait to building of image become finished then from OKD dashboard > deployment > openshift-event-exporter check for any running pod
If any error acuard you can redeploy it.

## Environment Variable
This exporter accepts 3 Environment variable:

### OPENSHIFT_API_URL
For openshift API URL endpoint, default variable set to `https://okd.private.teh-1.snappcloud.io`, if you want to use different URL set it in a deployment config file.
### OPENSHIFT_TOKEN
Service account token that you need to access the event uri in the openshift api.
If this variable not is set it will use `/var/run/secrets/kubernetes.io/serviceaccount/token` inside the pod, so you must add the service account that has access to event API.
### OPENSHIFT_NAMESPACE
This is your project name that must set if this variable no set application won't be run

## Prometheus Endpoint

Prometheus endpoint uses port `8090` and `/metric`
Metric name is `snappcloud_event_openshift` by default.
