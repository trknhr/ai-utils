package provider

// stripFlagWithValue removes occurrences of a flag that consumes the next value.
// Example: ["--model","gpt-5","-s"] -> ["-s"] when flag="--model".
func stripFlagWithValue(args []string, flag string) []string {
	out := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == flag {
			if i+1 < len(args) {
				i++
			}
			continue
		}
		out = append(out, args[i])
	}
	return out
}

