package pgfumlsd

import (
	"fmt"
	"strings"
	"text/template"
)

func expandEOL(a string) string {
	if strings.Contains(a, "\n") {
		return fmt.Sprintf("\\shortstack{%s}", strings.Replace(a, "\n", "\\\\", -1))
	}
	return a
}

func anchor(list []string, src, dst, fromOrTo string) string {

	for _, v := range list {
		if v == src {
			switch fromOrTo {
			case "from":
				return "east"
			default:
				return "west"
			}
		}
		if v == dst {
			switch fromOrTo {
			case "from":
				return "west"
			default:
				return "east"
			}
		}
	}

	return "east"

}

func instSize(list []string, abbr string) int {

	if list[0] == abbr {
		return 0
	}
	return 6
}

// GetTemplate return a parsed template
func GetTemplate() *template.Template {

	funcMap := template.FuncMap{
		"expandEOL": expandEOL,
		"anchor":    anchor,
		"instSize":  instSize,
	}

	return template.Must(template.New("pgfumlsd").Funcs(funcMap).Delims("##", "##").Parse(theTemplate))

}

const theTemplate = `
##- template "header" ##
##- template "document" . ##

##- define "header" ##
\documentclass[tikz,border=3mm]{standalone}
\usepackage[underline=false,roundedcorners=true]{pgf-umlsd}
\usepackage{underscore}
\usepackage{syntax}
\usetikzlibrary{shadows,positioning}
\tikzset{every shadow/.style={fill=none,shadow xshift=0pt,shadow yshift=0pt}}
##- end ##

##- define "document" ##

\begin{document}

	\sffamily
	\small

    \begin{sequencediagram}
    \tikzstyle{inststyle}+=[drop shadow={opacity=0.9,fill=lightgray}]

	\def\unitfactor{1.2}

##- template "actors" . ##
##- template "sequences" . ##

    \end{sequencediagram}
\end{document}
##- end ##

##- define "actors" ##

	##- $actorList := .ActorList ##

	##- range $idx, $actor := .Actor ##
		##- range $abbr, $name := $actor ##
\newinst[## instSize $actorList $abbr ##]{##$abbr##}{##$name##}
		##- end ##
	##- end ##

##- end ##

##- define "sequences" ##

	##- $actorList := .ActorList ##
	##- range .Sequence ##
	\mess[## .Delay ##]{## .Source ##}{## expandEOL .Label ##}{## .Destination ##}
	\node [anchor=## anchor $actorList .Source .Destination "from"  ##] at (mess from) {## expandEOL .AnnotationFrom ##};
	\node [anchor=## anchor $actorList .Source .Destination "to"  ##] at (mess to) {## expandEOL .AnnotationTo ##};
	##- end##

##- end ##
`
