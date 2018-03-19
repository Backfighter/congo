#!/usr/bin/env bash

format="{{ range .TestImports }}{{ println . }}{{ end }}{{ range .Imports }}{{ println . }}{{ end }}"

go list -f "${format}" ../... | sort | uniq | grep -v "gitlab.com/silentteacup/congo/*"
