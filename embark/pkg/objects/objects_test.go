package objects

import (
	"bytes"
	"testing"
	"text/template"

	"gopkg.in/yaml.v2"

	"k8s.io/helm/pkg/proto/hapi/chart"
)

type M = map[string]interface{}

func TestObjects(test *testing.T) {
	var objects = Objects{
		"svc":      bytes.NewBufferString(testSvc),
		"deploy":   bytes.NewBufferString(testDepl),
		"_helpers": bytes.NewBufferString(testHelpers),
	}
	var result = make(map[string]*bytes.Buffer, 3)
	var v map[string]interface{}
	yaml.Unmarshal([]byte(testValues), &v)
	var values = M{
		"Values": v,
		"Release": M{
			"Name": "containerum",
		},
		"Chart":        chart.Metadata{},
		"Files":        map[string][]byte{},
		"Capabilities": M{},
		"mongodb": M{
			"image":    "mongodb",
			"fullname": "name",
		},
		"containerum": M{
			"image": "containerum",
		},
	}
	if err := objects.Render(template.New("testObjects"), result, values); err != nil {
		test.Log(err)
	}
	for name, body := range result {
		test.Log(name, body)
	}
}

const testValues = `
image:
  registry: docker.io
  ## Bitnami MongoDB image name
  ##
  repository: bitnami/mongodb
  ## Bitnami MongoDB image tag
  ## ref: https://hub.docker.com/r/bitnami/mongodb/tags/
  ##
  tag: 3.6.6-debian-9

  ## Specify a imagePullPolicy
  ## ref: http://kubernetes.io/docs/user-guide/images/#pre-pulling-images
  ##
  pullPolicy: Always
  ## Optionally specify an array of imagePullSecrets.
  ## Secrets must be manually created in the namespace.
  ## ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/
  ##
  # pullSecrets:
  #   - myRegistrKeySecretName

## Enable authentication
## ref: https://docs.mongodb.com/manual/tutorial/enable-authentication/
#
usePassword: true
# existingSecret: name-of-existing-secret

## MongoDB admin password
## ref: https://github.com/bitnami/bitnami-docker-mongodb/blob/master/README.md#setting-the-root-password-on-first-run
##
# mongodbRootPassword:

## MongoDB custom user and database
## ref: https://github.com/bitnami/bitnami-docker-mongodb/blob/master/README.md#creating-a-user-and-database-on-first-run
##
# mongodbUsername: username
# mongodbPassword: password
# mongodbDatabase: database

## MongoDB additional command line flags
##
## Can be used to specify command line flags, for example:
##
## mongodbExtraFlags:
##  - "--wiredTigerCacheSizeGB=2"
mongodbExtraFlags: []

## Kubernetes service type
service:
  type: ClusterIP
  port: 27017

  ## Specify the nodePort value for the LoadBalancer and NodePort service types.
  ## ref: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
  ##
  # nodePort:


replicaSet:
  ## Whether to create a MongoDB replica set for high availability or not
  enabled: false
  ## Name of the replica set
  ##
  name: rs0

  ## Key used for replica set authentication
  ##
  # key: key

  ## Number of replicas per each node type
  ##
  replicas:
    secondary: 1
    arbiter: 1
  ## Pod Disruption Budget
  ## ref: https://kubernetes.io/docs/concepts/workloads/pods/disruptions/
  pdb:
    minAvailable:
      primary: 1
      secondary: 1
      arbiter: 1

# Annotations to be added to MongoDB pods
podAnnotations: {}

## Configure resource requests and limits
## ref: http://kubernetes.io/docs/user-guide/compute-resources/
##
resources: {}
# limits:
#   cpu: 500m
#   memory: 512Mi
# requests:
#   cpu: 100m
#   memory: 256Mi

## Node selector
## ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector
nodeSelector: {}

## Affinity
## ref: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity
affinity: {}

## Tolerations
## ref: https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
tolerations: []

## Enable persistence using Persistent Volume Claims
## ref: http://kubernetes.io/docs/user-guide/persistent-volumes/
##
persistence:
  enabled: true
  ## A manually managed Persistent Volume and Claim
  ## Requires persistence.enabled: true
  ## If defined, PVC must be created manually before volume will be bound
  # existingClaim:

  ## mongodb data Persistent Volume Storage Class
  ## If defined, storageClassName: <storageClass>
  ## If set to "-", storageClassName: "", which disables dynamic provisioning
  ## If undefined (the default) or set to null, no storageClassName spec is
  ##   set, choosing the default provisioner.  (gp2 on AWS, standard on
  ##   GKE, AWS & OpenStack)
  ##
  # storageClass: "-"
  accessModes:
    - ReadWriteOnce
  size: 8Gi
  annotations: {}

## Configure extra options for liveness and readiness probes
## ref: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/#configure-probes)
livenessProbe:
  enabled: true
  initialDelaySeconds: 30
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 6
  successThreshold: 1
readinessProbe:
  enabled: true
  initialDelaySeconds: 5
  periodSeconds: 10
  timeoutSeconds: 5
  failureThreshold: 6
  successThreshold: 1

# Entries for the MongoDB config file
configmap:
#  # Where and how to store data.
#  storage:
#    dbPath: /opt/bitnami/mongodb/data/db
#    journal:
#      enabled: true
#    #engine:
#    #wiredTiger:
#  # where to write logging data.
#  systemLog:
#    destination: file
#    logAppend: true
#    path: /opt/bitnami/mongodb/logs/mongodb.log
#  # network interfaces
#  net:
#    port: 27017
#    bindIp: 0.0.0.0
#    unixDomainSocket:
#      enabled: true
#      pathPrefix: /opt/bitnami/mongodb/tmp
#  # replica set options
#  #replication:
#  #  replSetName: replicaset
#  # process management options
#  processManagement:
#     fork: false
#     pidFilePath: /opt/bitnami/mongodb/tmp/mongodb.pid
#  # set parameter options
#  setParameter:
#     enableLocalhostAuthBypass: true
#  # security options
#  security:
#    authorization: enabled
#    #keyFile: /opt/bitnami/mongodb/conf/keyfile
`
const testSvc = `{{- if not .Values.replicaSet.enabled }}
apiVersion: v1
kind: Service
metadata:
  name: {{ template "mongodb.fullname" . }}
  labels:
    app: {{ template "mongodb.name" . }}
    chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
    release: "{{ .Release.Name }}"
    heritage: "{{ .Release.Service }}"
spec:
  type: {{ .Values.service.type }}
  ports:
  - name: mongodb
    port: 27017
    targetPort: mongodb
{{- if .Values.service.nodePort }}
    nodePort: {{ .Values.service.nodePort }}
{{- end }}
  selector:
    app: {{ template "mongodb.name" . }}
    release: "{{ .Release.Name }}"
{{- end }}
`

const testDepl = `{{- if not .Values.replicaSet.enabled }}
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: {{ template "mongodb.fullname" . }}
  labels:
    app: {{ template "mongodb.name" . }}
    chart: "{{ .Chart.Name }}-{{ .Chart.Version }}"
    release: "{{ .Release.Name }}"
    heritage: "{{ .Release.Service }}"
spec:
  template:
    metadata:
      labels:
        app: {{ template "mongodb.name" . }}
        release: "{{ .Release.Name }}"
    spec:
      {{- if .Values.nodeSelector }}
      nodeSelector:
{{ toYaml .Values.nodeSelector | indent 8 }}
      {{- end }}
      {{- if .Values.image.pullSecrets }}
      imagePullSecrets:
      {{- range .Values.image.pullSecrets }}
        - name: {{ . }}
      {{- end}}
      {{- end }}
      containers:
      - name: {{ template "mongodb.fullname" . }}
        image: {{ template "mongodb.image" . }}
        imagePullPolicy: {{ .Values.image.pullPolicy | quote }}
        env:
        {{- if .Values.usePassword }}
        - name: MONGODB_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: {{ template "mongodb.fullname" . }}
              key: mongodb-root-password
        {{- end }}
        - name: MONGODB_USERNAME
          value: {{ default "" .Values.mongodbUsername | quote }}
        {{- if and .Values.mongodbUsername .Values.mongodbDatabase }}
        - name: MONGODB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: {{ template "mongodb.fullname" . }}
              key: mongodb-password
        {{- end }}
        - name: MONGODB_DATABASE
          value: {{ default "" .Values.mongodbDatabase | quote }}
        - name: MONGODB_EXTRA_FLAGS
          value: {{ default "" .Values.mongodbExtraFlags | join " " }}
        ports:
        - name: mongodb
          containerPort: 27017
        {{- if .Values.livenessProbe.enabled }}
        livenessProbe:
          exec:
            command:
            - mongo
            - --eval
            - "db.adminCommand('ping')"
          initialDelaySeconds: {{ .Values.livenessProbe.initialDelaySeconds }}
          periodSeconds: {{ .Values.livenessProbe.periodSeconds }}
          timeoutSeconds: {{ .Values.livenessProbe.timeoutSeconds }}
          successThreshold: {{ .Values.livenessProbe.successThreshold }}
          failureThreshold: {{ .Values.livenessProbe.failureThreshold }}
        {{- end }}
        {{- if .Values.readinessProbe.enabled }}
        readinessProbe:
          exec:
            command:
            - mongo
            - --eval
            - "db.adminCommand('ping')"
          initialDelaySeconds: {{ .Values.readinessProbe.initialDelaySeconds }}
          periodSeconds: {{ .Values.readinessProbe.periodSeconds }}
          timeoutSeconds: {{ .Values.readinessProbe.timeoutSeconds }}
          successThreshold: {{ .Values.readinessProbe.successThreshold }}
          failureThreshold: {{ .Values.readinessProbe.failureThreshold }}
        {{- end }}
        volumeMounts:
        - name: data
          mountPath: /bitnami/mongodb
        {{- if .Values.configmap }}
        - name: config
          mountPath: /opt/bitnami/mongodb/conf/mongodb.conf
          subPath: mongodb.conf
        {{- end }}
        resources:
{{ toYaml .Values.resources | indent 10 }}
      volumes:
      - name: data
      {{- if .Values.persistence.enabled }}
        persistentVolumeClaim:
          claimName: {{ if .Values.persistence.existingClaim }}{{ .Values.persistence.existingClaim }}{{- else }}{{ template "mongodb.fullname" . }}{{- end }}

      {{- else }}
        emptyDir: {}
      {{- end -}}
      {{- if .Values.configmap }}
      - name: config
        configMap:
          name: {{ template "mongodb.fullname" . }}
      {{- end }}
 {{- end -}}
`

const testHelpers = `{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "mongodb.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "mongodb.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "mongodb.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create the name for the admin secret.
*/}}
{{- define "mongodb.adminSecret" -}}
    {{- if .Values.auth.existingAdminSecret -}}
        {{- .Values.auth.existingAdminSecret -}}
    {{- else -}}
        {{- template "mongodb.fullname" . -}}-admin
    {{- end -}}
{{- end -}}

{{/*
Create the name for the key secret.
*/}}
{{- define "mongodb.keySecret" -}}
    {{- if .Values.auth.existingKeySecret -}}
        {{- .Values.auth.existingKeySecret -}}
    {{- else -}}
        {{- template "mongodb.fullname" . -}}-keyfile
    {{- end -}}
{{- end -}}

{{/*
Return the proper image name
*/}}
{{- define "mongodb.image" -}}
{{- $registryName :=  .Values.image.registry -}}
{{- $repositoryName := .Values.image.repository -}}
{{- $tag := .Values.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}
`
