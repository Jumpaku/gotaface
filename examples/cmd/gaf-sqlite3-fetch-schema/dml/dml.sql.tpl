{{- /* Go Template */ -}}
CREATE TABLE {{.Name}}(
{{ range $Index, $Column := .Columns }}
    {{$Column.Name}} {{$Column.Type}} {{ if (not $Column.Nullable) }}NOT NULL{{ end }},
{{ end }}
{{ range $Index, $ForeignKey := .ForeignKeys }}
    FOREIGN KEY (
        {{ range $Index, $Name := $ForeignKey.ReferencingKey }}
            {{ if $Index }}, {{ end }}{{$Name}}
        {{ end }}
    ) {{$ForeignKey.ReferencedTable}} (
        {{ range $Index, $Name := $ForeignKey.ReferencedKey }}
            {{ if $Index }}, {{ end }}{{$Name}}
        {{ end }}
    ),
{{ end }}
{{ range $Index, $UniqueKey := .UniqueKeys }}
    {{if not $UniqueKey.Name}}
    UNIQUE (
        {{ range $Index, $Name := $UniqueKey.Key }}
        {{ if $Index }}, {{ end }}{{$Name}}
        {{ end }}
    ),
    {{ end }}
{{ end }}
    PRIMARY KEY (
        {{ range $Index, $Name := .PrimaryKey }}
        {{ if $Index }}, {{ end }}{{$Name}}
        {{ end }}
    )
);

{{- $Table := .Name -}}
{{ range $Index, $UniqueKey := .UniqueKeys }}
{{if $UniqueKey.Name}}
CREATE UNIQUE INDEX {{$UniqueKey.Name}} ON {{$Table}}(
    {{ range $Index, $Name := $UniqueKey.Key }}
    {{ if $Index }}, {{ end }}{{$Name}}
    {{ end }}
);
{{ end }}
{{ end }}