{{- define "app.name" -}}
{{- .Release.Name -}}
{{- end }}

{{- define "app.name" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{- define "app.labels" -}}
app.kubernetes.io/chart: {{ .Chart.Name }}
app.kubernetes.io/name: {{ .Release.Name }}
app.kubernetes.io/revision: "{{ .Release.Revision }}"
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{- if not (has .Values.exposure.type (list "ingress" "gatewayapi" "route" "istio")) }}
{{- fail "Invalid exposure.type" }}
{{- end }}