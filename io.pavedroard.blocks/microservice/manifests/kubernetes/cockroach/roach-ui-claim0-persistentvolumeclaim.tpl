{{define "roach-ui-claim0-persistentvolumeclaim.tpl"}}
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: roach-ui-claim0
  labels:
    pavedroad.service: roach-ui-claim0
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 100Mi
status: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
