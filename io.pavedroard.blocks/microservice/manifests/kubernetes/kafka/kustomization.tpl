{{define "kustomization.tpl"}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - kafka-broker-deployment.yaml
  - kafka-broker-service.yaml
  - zookeeper-deployment.yaml
  - zookeeper-service.yaml
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
