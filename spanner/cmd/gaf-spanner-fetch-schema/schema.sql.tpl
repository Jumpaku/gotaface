{{- /* Go Template */ -}}
CREATE TABLE {{.Name}}(
{{- range $Index, $Column := .Columns -}}
    {{"\n    "}}{{$Column.Name}} {{$Column.Type}} {{- if (not $Column.Nullable) -}}{{" "}}NOT NULL{{- end -}},
{{- end -}}
{{- range $Index, $ForeignKey := .ForeignKeys -}}
    {{"\n    "}}CONSTRAINT {{$ForeignKey.Name}} FOREIGN KEY (
        {{- range $Index, $Name := $ForeignKey.ReferencingKey -}}
            {{- if $Index -}},{{- end -}}{{$Name}}
        {{- end -}}
    ) {{$ForeignKey.ReferencedTable}} (
        {{- range $Index, $Name := $ForeignKey.ReferencedKey -}}
            {{- if $Index -}},{{- end -}}{{$Name}}
        {{- end -}}
    ),
{{- end -}}
{{"\n"}}) PRIMARY KEY (
{{- range $Index, $Name := .PrimaryKey -}}
    {{- if $Index -}},{{- end -}}{{$Name}}
{{- end -}}
    )
{{- if .Parent -}}
    ,{{"\n    "}}INTERLEAVE IN PARENT {{.Parent}} ON DELETE CASCADE
{{- end -}}
;{{"\n"}}

{{- range $Index, $UniqueKey := .UniqueKeys -}}
    CREATE UNIQUE INDEX {{$UniqueKey.Name}} ON {{.Name}} (
    {{- range $Index, $Name := $UniqueKey.Key -}}
        {{- if $Index -}},{{- end -}}{{$Name}}
    {{- end -}}
    );{{"\n"}}
{{- end -}}