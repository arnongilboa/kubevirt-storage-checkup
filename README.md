# kubevirt-storage-checkup

This checkup performs storage checks, validating storage is working correctly for virtual machines. The checkup is using the Kiagnose engine.

## Permissions

Cluster admin should create the following permissions and/for dedicated storage checkup ServiceAccount and namespace:

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: storage-checkup-sa
  namespace: <target-namespace>
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubevirt-storage-checker
rules:
- apiGroups: [ "storage.k8s.io" ]
  resources: [ "storageclasses" ]
  verbs: [ "get", "list" ]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubevirt-storage-checker
subjects:
- kind: ServiceAccount
  name: storage-checkup-sa
  namespace: <target-namespace>
roleRef:
  kind: ClusterRole
  name: kubevirt-storage-checker
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kubevirt-storage-checker-volumesnapshotclasses
rules:
- apiGroups: [ "snapshot.storage.k8s.io" ]
  resources: [ "volumesnapshotclasses" ]
  verbs: [ "get", "list" ]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubevirt-storage-checker-volumesnapshotclasses
subjects:
- kind: ServiceAccount
  name: storage-checkup-sa
  namespace: <target-namespace>
roleRef:
  kind: ClusterRole
  name: kubevirt-storage-checker-volumesnapshotclasses
  apiGroup: rbac.authorization.k8s.io
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: kiagnose-configmap-access
  namespace: <target-namespace>
rules:
- apiGroups: [ "" ]
  resources: [ "configmaps" ]
  verbs: ["get", "update"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: kiagnose-configmap-access
  namespace: <target-namespace>
subjects:
- kind: ServiceAccount
  name: storage-checkup-sa
roleRef:
  kind: Role
  name: kiagnose-configmap-access
  apiGroup: rbac.authorization.k8s.io
```

## Configuration

| Key                                         | Description                                                                                                       | Is Mandatory | Remarks                                                                             |
|---------------------------------------------|-------------------------------------------------------------------------------------------------------------------|--------------|-------------------------------------------------------------------------------------|
| spec.timeout                                | How much time before the checkup will try to close itself                                                         | True         |                                                                                     |

### Example

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: storage-checkup-config
  namespace: <target-namespace>
data:
  spec.timeout: 5m
```

## Execution

In order to execute the checkup, fill in the required data and apply this manifest:
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: storage-checkup
  namespace: <target-namespace>
spec:
  backoffLimit: 0
  template:
    spec:
      serviceAccount: storage-checkup-sa
      restartPolicy: Never
      containers:
        - name: storage-checkup
          image: quay.io/kiagnose/kubevirt-storage-checkup:main
          imagePullPolicy: Always
          env:
            - name: CONFIGMAP_NAMESPACE
              value: <target-namespace>
            - name: CONFIGMAP_NAME
              value: storage-checkup-config
```

You can create the permissions, ConfigMap and Job with:
```bash
export CHECKUP_NAMESPACE=<target-namespace>

envsubst < manifests/storage_checkup_permissions.yaml | kubectl apply -f -
envsubst < manifests/storage_checkup.yaml | kubectl apply -f -
```

and cleanup the ConfigMap and Job:
```bash
envsubst < manifests/storage_checkup.yaml | kubectl delete -f -
```
## Checkup Results Retrieval

After the checkup Job had completed, the results are made available at the user-supplied ConfigMap object:

```bash
kubectl get configmap storage-checkup-config -n <target-namespace> -o yaml
```
| Key                                              | Description                                                       | Remarks  |
|--------------------------------------------------|-------------------------------------------------------------------|----------|
| status.succeeded                                 | Has the checkup succeeded                                         |          |
| status.failureReason                             | Failure reason in case of a failure                               |          |
| status.startTimestamp                            | Checkup start timestamp                                           | RFC 3339 |
| status.completionTimestamp                       | Checkup completion timestamp                                      | RFC 3339 |
| status.result.defaultStorageClass                | Indicates whether there is a default storage class                |          |
| status.result.storageProfilesWithEmptyClaimPropertySets | StorageProfiles with empty claimPropertySets (unknown provisioners) |          |
| status.result.storageProfilesWithSpecClaimPropertySets  | StorageProfiles with spec-overrriden claimPropertySets              |          |
| status.result.storageWithRWX                     | Storage with RWX access mode                                        |          |
| status.result.storageMissingVolumeSnapshotClass  | Storage using snapshot-based clone but missing VolumeSnapshotClass  |          |
