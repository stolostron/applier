{{- define "myshortname.name" -}}
{{- .Values.name | trunc 4 -}}
{{- end -}}
