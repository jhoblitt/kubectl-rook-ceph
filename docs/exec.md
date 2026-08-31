# Exec

The `exec` command opens an interactive Bash shell in a Rook Ceph component pod. If no component is specified, it opens a shell in the toolbox pod.

The target pod must be running, and the current Kubernetes user must have permission to list pods and execute commands in it. The toolbox and Ceph daemon pods are selected from the Ceph cluster namespace. The operator pod is selected from the operator namespace.

## Usage

Open a toolbox shell:

```bash
kubectl rook-ceph exec
kubectl rook-ceph exec toolbox
```

Open a shell in the Rook Ceph operator:

```bash
kubectl rook-ceph exec operator
```

Open a shell in a specific Ceph daemon:

```bash
kubectl rook-ceph exec mon-a
kubectl rook-ceph exec mgr-b
kubectl rook-ceph exec osd-42
kubectl rook-ceph exec mds-myfs-a
kubectl rook-ceph exec rgw-my-store
```

For MDS, use the daemon ID shown in the pod's `mds` label. For RGW, use the object store name shown in the pod's `rgw` label.

To use a component in a different namespace:

```bash
kubectl rook-ceph --namespace my-rook-ceph-namespace exec mon-a
```

To use an operator in a different namespace:

```bash
kubectl rook-ceph --operator-namespace my-operator-namespace exec operator
```

Run `exit` to leave the shell.
