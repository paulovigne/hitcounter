{{- define "hitcounter.name" -}}
hitcounter
{{- end }}

{{- define "hitcounter.labels" -}}
app: {{ include "hitcounter.name" . }}
{{- end }}
