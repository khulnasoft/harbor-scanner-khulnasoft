{{/*
Expand the name of the chart.
*/}}
{{- define "harbor-scanner-khulnasoft.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "harbor-scanner-khulnasoft.fullname" -}}
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
{{- define "harbor-scanner-khulnasoft.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "harbor-scanner-khulnasoft.labels" -}}
app.kubernetes.io/name: {{ include "harbor-scanner-khulnasoft.name" . }}
helm.sh/chart: {{ include "harbor-scanner-khulnasoft.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Return the proper imageRef as used by the init conainer template spec.
*/}}
{{- define "harbor-scanner-khulnasoft.scannerImageRef" -}}
{{- $registryName := .Values.khulnasoft.registry.server -}}
{{- $repositoryName := "scanner" -}}
{{- $tag := .Values.khulnasoft.version | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{/*
Return the proper imageRef as used by the container template spec.
*/}}
{{- define "harbor-scanner-khulnasoft.adapterImageRef" -}}
{{- $registryName := .Values.scanner.image.registry -}}
{{- $repositoryName := .Values.scanner.image.repository -}}
{{- $tag := .Values.scanner.image.tag | toString -}}
{{- printf "%s/%s:%s" $registryName $repositoryName $tag -}}
{{- end -}}

{{- define "imagePullSecret" -}}
{{- printf "{\"auths\": {\"%s\": {\"auth\": \"%s\"}}}" .Values.khulnasoft.registry.server (printf "%s:%s" .Values.khulnasoft.registry.username .Values.khulnasoft.registry.password | b64enc) | b64enc }}
{{- end }}

{{/*
Return the proper scheme for liveness and readiness probe spec.
*/}}
{{- define "probeScheme" -}}
{{- if .Values.scanner.api.tlsEnabled -}}
HTTPS
{{- else -}}
HTTP
{{- end -}}
{{- end -}}
