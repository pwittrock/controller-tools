/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import "text/template"

var (
	configTemplate = template.Must(template.New("template").Parse(`
{{- if not .DisableCreateRBAC }}
#
# RBAC: Leader election.
#
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ .Name }}-leader-election-role
  namespace: {{ .Name }}-system
rules:
- apiGroups:
  - ""
  resources:
  - configmaps
  verbs:
  - get
  - list
  - watch
  - create
  - update
  - patch
  - delete
- apiGroups:
  - ""
  resources:
  - configmaps/status
  verbs:
  - get
  - update
  - patch
- apiGroups:
  - ""
  resources:
  - events
  verbs:
  - create
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRole
metadata:
  name: {{ .Name }}-metrics-reader
rules:
- nonResourceURLs: ["/metrics"]
  verbs: ["get"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ .Name }}-leader-election-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: {{ .Name }}-leader-election-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: {{ .Name }}-system
---
#
# RBAC: Manager permissions
#
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ .Name }}-manager-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{ .Name }}-manager-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: {{ .Name }}-system
---
{{- if not .DisableAuthProxy}}
#
# RBAC: Metrics auth proxy
#
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: {{ .Name }}-proxy-role
rules:
- apiGroups: ["authentication.k8s.io"]
  resources:
  - tokenreviews
  verbs: ["create"]
- apiGroups: ["authorization.k8s.io"]
  resources:
  - subjectaccessreviews
  verbs: ["create"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: {{ .Name }}-proxy-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: {{ .Name }}-proxy-role
subjects:
- kind: ServiceAccount
  name: default
  namespace: {{ .Name }}-system
---
{{- end }}
{{- end }}
{{- if .EnablePrometheus }}
#
# Prometheus Monitor Service (Metrics)
#
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  namespace: {{.Name}}-system
  name: controller-manager-metrics-monitor
  labels:
    instance: {{.Name}}
    control-plane: controller-manager
spec:
  endpoints:
    - path: /metrics
      port: https
  selector:
    matchLabels:
      instance: {{.Name}}
      control-plane: controller-manager
---
{{- end }}
{{- if not .DisableCreateNamespace }}
apiVersion: v1
kind: Namespace
metadata:
  name: {{.Name}}-system
  labels:
    instance: {{.Name}}
    control-plane: controller-manager
---
{{- end }}
{{- if .EnableCertManager}}
# The following manifests contain a self-signed issuer CR and a certificate CR.
# More document can be found at https://docs.cert-manager.io
# WARNING: Targets CertManager 0.11 check https://docs.cert-manager.io/en/latest/tasks/upgrading/index.html for 
# breaking changes
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  name: {{.Name}}-selfsigned-issuer
  namespace: {{.Name}}-system
spec:
  selfSigned: {}
---
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  name: {{.Name}}-serving-cert
  namespace: {{.Name}}-system
spec:
  dnsNames:
  - {{.Name}}-webhook-service.{{.Name}}-system.svc
  - {{.Name}}-webhook-service.{{.Name}}-system.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: selfsigned-issuer
  secretName: webhook-server-cert # this secret will not be prefixed, since it's not managed by kustomize
---
{{- end }}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: {{.Name}}-system
  labels:
    instance: {{ .Name }}
    control-plane: controller-manager
spec:
  selector:
    matchLabels:
      instance: {{ .Name }}
      control-plane: controller-manager
  replicas: 1
  template:
    metadata:
      labels:
        instance: {{ .Name }}
        control-plane: controller-manager
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: manager
        image: {{ .Image }}
        command: [ "/manager" ]
        args:
        - --enable-leader-election
{{- if not .DisableAuthProxy }}
        - --metrics-addr=127.0.0.1:8080
{{- end }}
        resources:
          limits:
            cpu: 100m
            memory: 30Mi
          requests:
            cpu: 100m
            memory: 20Mi
{{- if .EnableWebhooks }}
        ports:
        - containerPort: 9443
          name: webhook-server
          protocol: TCP
        volumeMounts:
        - mountPath: /tmp/k8s-webhook-server/serving-certs
          name: cert
          readOnly: true
{{- end }}
{{- if not .DisableAuthProxy }}
      - name: kube-rbac-proxy
        image: gcr.io/kubebuilder/kube-rbac-proxy:v0.5.0
        args:
        - "--secure-listen-address=0.0.0.0:8443"
        - "--upstream=http://127.0.0.1:8080/"
        - "--logtostderr=true"
        - "--v=10"
        ports:
        - containerPort: 8443
          name: https
{{- end }}
{{- if .EnableWebhooks }}
    volumes:
    - name: cert
      secret:
        defaultMode: 420
        secretName: webhook-server-cert
{{- end }}
{{- if .EnableWebhooks }}
---
apiVersion: v1
kind: Service
metadata: 
  namespace: {{.Name}}-system
  name: {{.Name}}-webhook-service
  labels:
    instance: {{.Name}}
    control-plane: webhook
spec:
  ports:
  - port: 443
    targetPort: webhook-server
  selector:
    instance: {{.Name}}
    control-plane: controller-manager
{{- end }}
{{- if not .DisableAuthProxy }}
---
apiVersion: v1
kind: Service
metadata:
  namespace: {{ .Name }}-system
  name: {{ .Name }}-metrics-service
  labels:
    control-plane: controller-manager
    instance: {{ .Name }}
spec:
  ports:
  - name: https
    port: 8443
    targetPort: https
  selector:
    control-plane: controller-manager
    instance: {{ .Name }}
{{ end -}}
`))
)
