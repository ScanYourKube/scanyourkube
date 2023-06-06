{{/*
Expand the name of the chart.
*/}}
{{- define "keel-webhook.name" -}}
{{- default .Chart.Name .Values.scanyourkubewebhook.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "keel-webhook.fullname" -}}
{{- if .Values.scanyourkubewebhook.fullnameOverride }}
{{- .Values.scanyourkubewebhook.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.scanyourkubewebhook.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "keel-webhook.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "keel-webhook.labels" -}}
helm.sh/chart: {{ include "keel-webhook.chart" . }}
{{ include "keel-webhook.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "keel-webhook.selectorLabels" -}}
app.kubernetes.io/name: {{ include "keel-webhook.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "keel-webhook.serviceAccountName" -}}
{{- if .Values.scanyourkubewebhook.serviceAccount.create }}
{{- default (include "keel-webhook.fullname" .) .Values.scanyourkubewebhook.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.scanyourkubewebhook.serviceAccount.name }}
{{- end }}
{{- end }}
