package cmd

import (
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type completionCmd struct{}

func init() {
	var c completionCmd

	cmd := cobra.Command{
		Use:     "completion",
		Short:   "generate shell completion",
		Example: "eval \"$(" + name + " completion)\"",
		RunE: func(_ *cobra.Command, _ []string) error {
			return c.run()
		},
	}

	rootCmd.AddCommand(&cmd)
}

func (c completionCmd) run() error {
	switch shell := filepath.Base(os.Getenv("SHELL")); shell {
	case "zsh":
		return runCompletionZsh(os.Stdout)
	case "bash":
		return rootCmd.GenBashCompletion(os.Stdout)
	default:
		return errors.Errorf("unsupported shell: %s", shell)
	}
}

func runCompletionZsh(out io.Writer) error {
	if err := zshCompletionTmpl.Execute(out, name); err != nil {
		return err
	}

	if err := rootCmd.GenBashCompletion(out); err != nil {
		return err
	}

	return zshCompletionTailTmpl.Execute(out, name)
}

var (
	zshCompletionTmpl     = template.Must(template.New("zsh").Parse(zshCompletionTmplStr))
	zshCompletionTailTmpl = template.Must(template.New("zshTail").Parse(zshCompletionTailTmplStr))
)

const zshCompletionTailTmplStr = `
BASH_COMPLETION_EOF
}
__{{.}}_bash_source <(__{{.}}_convert_bash_to_zsh)
`

const zshCompletionTmplStr = `
__{{.}}_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__{{.}}_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__{{.}}_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__{{.}}_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__{{.}}_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__{{.}}_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}
__{{.}}_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__{{.}}_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__{{.}}_filedir() {
	local RET OLD_IFS w qw
	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__{{.}}_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__{{.}}_quote() {
    if [[ $1 == \'* || $1 == \"* ]]; then
        # Leave out first character
        printf %q "${1:1}"
    else
    	printf %q "$1"
    fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__{{.}}_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__{{.}}_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__{{.}}_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__{{.}}_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__{{.}}_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__{{.}}_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__{{.}}_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__{{.}}_type/g" \
	<<'BASH_COMPLETION_EOF'
`
