package app

type KeyCommandFactory struct{}

func (kcf KeyCommandFactory) CreateKeyCommand(key string) KeyCommand {
	switch key {
	case "j":
		return JKeyCommand{}
	case "k":
		return KKeyCommand{}
	case "enter":
		return EnterKeyCommand{}
	case "e":
		return EKeyCommand{}
	case "m":
		return MKeyCommand{}
	case "t":
		return TKeyCommand{}
	case "l":
		return LKeyCommand{}
	case "h":
		return HKeyCommand{}
	case "esc":
		return EscKeyCommand{}
	case "s":
		return SKeyCommand{}
	case "o":
		return OKeyCommand{}
	case "n":
		return NKeyCommand{}
	}
	return NilKeyCommand{}
}
