package utils

import "os"

func GetCommand() string {
	if len(os.Args) < 2 {
		return ""
	}

	return os.Args[1]
}

func GetCommandArgs() []string {
	if len(os.Args) < 3 {
		return []string{}
	}

	return os.Args[2:]
}

func GetCommandArg(index int) string {
	args := GetCommandArgs()
	if index < 0 || index >= len(args) {
		return ""
	}

	return args[index]
}

func GetCommandArgInt(index int) int {
	arg := GetCommandArg(index)
	if arg == "" {
		return 0
	}

	return StringToInt(arg, 0)
}

func GetCommandArgOrDefault(index int, defaultValue string) string {
	arg := GetCommandArg(index)
	if arg == "" {
		return defaultValue
	}

	return arg
}

func GetCommandArgOrDefaultInt(index int, defaultValue int) int {
	arg := GetCommandArg(index)
	if arg == "" {
		return defaultValue
	}

	return StringToInt(arg, defaultValue)
}
