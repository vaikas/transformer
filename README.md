# Simple Knative Transformer

This simple repo is used to demonstrate Knative Sequences. It creates a container
which listens for `CloudEvents` and assumes the incoming data is produced by
[Knative CronJobSource](https://knative.dev/docs/eventing/samples/cronjob-source/).

It's very simple and is just meant to demonstrate the `Sequence` concept. It takes
in an environment variable string "STEP" and upon receiving an event will read
in the data, and append " - Handled by $STEP" to it.

## Deploying

```shell
kubectl -n <NAMESPACE> create -f ./config/transformer.yaml
```

The above mentioned file can be re-created with (or if you need to modify the code, etc.):

```shell
ko resolve -f ./transformer.yaml > config/transformer.yaml
```

