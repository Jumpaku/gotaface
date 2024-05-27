{{- /* Go Template */ -}}
{{ range . }}

{{ if not .View }}
CREATE TABLE {{.Name}}(
{{ range $Index, $Column := .Columns }}
    {{$Column.Name}} {{$Column.Type}} {{ if (not $Column.Nullable) }}NOT NULL{{ end }},
{{ end }}
{{ range $Index, $ForeignKey := .ForeignKeys }}
    CONSTRAINT {{$ForeignKey.Name}} FOREIGN KEY (
        {{ range $Index, $Name := $ForeignKey.ReferencingKey }}
            {{ if $Index }}, {{ end }}{{$Name}}
        {{ end }}
    ) {{$ForeignKey.ReferencedTable}} (
        {{ range $Index, $Name := $ForeignKey.ReferencedKey }}
            {{ if $Index }}, {{ end }}{{$Name}}
        {{ end }}
    ),
{{ end }}
) PRIMARY KEY (
    {{ range $Index, $Name := .PrimaryKey }}
    {{ if $Index }},{{ end }}{{$Name}}
    {{ end }}
){{ if .Parent }}, INTERLEAVE IN PARENT {{.Parent}} ON DELETE CASCADE{{ end }};

{{ $table := . }}
{{ range $i, $Index := .Indexes }}
CREATE UNIQUE INDEX {{$Index.Name}} ON {{$table.Name}} (
    {{ range $Index, $Key := $Index.Key }}
    {{ if $Index }}, {{ end }}{{$Key.Name}} {{ if $Key.Desc }}DESC{{ end }}
    {{ end }}
);
{{ end }}

{{ end }}

{{ end }}
