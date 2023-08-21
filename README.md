# Triton Monitor

- Read memory from pods in cluster
- Annotate pod based off of https://kubernetes.io/docs/reference/labels-annotations-taints/#pod-deletion-cost
- Mulitply by -1 so higher memory usage gets lower score and more likely to be scaled down
