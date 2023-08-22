# Pod Cost Manager

## Overview

This service annotates pods in a cluster based off of https://kubernetes.io/docs/reference/labels-annotations-taints/#pod-deletion-cost. This adjusts the scaling behaviour of k8s to identify which pods to scale down first based on their cost.

The service multiplies the memory consumption of a pod by -1 so pods with higher memory usage are targetted for scale-down first.


## Deployment

The service will require a list of applications to annotate using the `APPLICATIONS` environment variable, which is a comma-separated list of application names.

Pod Cost Manager will filter for these applications using the `app.kubernetes.io/name` label, so this label will need to exist on these pods for them to be annotated.
