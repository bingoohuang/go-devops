package main

func PsAuxGrep(keywords ...string) bool {
	shellCmd := `ps axo args`

	count := 0
	for _, keyword := range keywords {
		if keyword != "" {
			count += 1
			shellCmd += `|grep ` + keyword
		}
	}

	if count > 0 {
		shellCmd += `|grep -v grep`
	}

	greped := false
	ExecuteBash(shellCmd, func(line string) bool {
		greped = true
		return false // break continue
	})
	return greped
}
