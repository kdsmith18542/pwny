// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 kdsmith18542

package cli

import (
	"strings"

	"github.com/chzyer/readline"
)

type ModuleCompleter struct {
	modules []string
}

func NewModuleCompleter(modules []string) *ModuleCompleter {
	return &ModuleCompleter{modules: modules}
}

func (m *ModuleCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	str := string(line[:pos])
	if strings.HasPrefix(str, "use ") {
		prefix := str[4:]
		var suggestions [][]rune
		for _, mod := range m.modules {
			if strings.HasPrefix(mod, prefix) {
				suggestions = append(suggestions, []rune(mod[len(prefix):]))
			}
		}
		return suggestions, len(prefix)
	}
	return nil, 0
}

func GetDefaultCompleter() readline.AutoCompleter {
	// Example modules for stub
	modules := []string{
		"exploits/windows/smb/ms17_010_eternalblue",
		"exploits/windows/smb/psexec",
		"auxiliary/scanner/ssh_brute",
		"post/windows/manage/mimikatz",
	}

	return readline.NewPrefixCompleter(
		readline.PcItem("use", readline.PcItemDynamic(func(line string) []string {
			var result []string
			for _, m := range modules {
				if strings.HasPrefix(m, line) {
					result = append(result, m)
				}
			}
			return result
		})),
		readline.PcItem("set"),
		readline.PcItem("show", readline.PcItem("options")),
		readline.PcItem("run"),
		readline.PcItem("sessions"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
	)
}
