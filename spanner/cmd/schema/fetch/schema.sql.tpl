{{- /* Go Template */ -}}
CREATE TABLE(
{{- range $Index, $Column := .Columns -}}
    {{$Column.Name}} {{$Column.Type}} {{- if (not $Column.Nullable) -}}{{" "}}NOT NULL{{- end -}},{{"\n"}}
{{- end -}}
{{- range $Index, $ForeignKey := .ForeignKeys -}}
    CONSTRAINT {{$ForeignKey.Name}} FOREIGN KEY (
        {{- range $Index, $Name := $ForeignKey.ReferencingKey -}}
            {{- if  -}},{{- end -}}{{$Name}}
        {{- end -}}
    ) {{$ForeignKey.ReferencedTable}} (
        {{- range $Index, $Name := $ForeignKey.ReferencedKey -}}
            {{- if  -}},{{- end -}}{{$Name}}
        {{- end -}}
    ),{{"\n"}}
{{- end -}}
) PRIMARY KEY (
{{- range $Index, $Name := .PrimaryKey -}}
    {{- if  -}},{{- end -}}{{$Name}}
{{- end -}}
    )
{{- if .PrimaryKey.Parent -}}
    ,{{"\n    "}}INTERLEAVE IN PARENT {{.PrimaryKey.Parent}} ON DELETE CASCADE
{{- end -}}