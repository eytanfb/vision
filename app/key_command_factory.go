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
	case "o":
		return OKeyCommand{}
	case "n":
		return NKeyCommand{}
	case "f":
		return FKeyCommand{}
	case "C":
		return UppercaseCKeyCommand{}
	case "Q":
		return UppercaseQKeyCommand{}
	case "L":
		return UppercaseLKeyCommand{}
	case "1":
		return OneKeyCommand{}
	case "2":
		return TwoKeyCommand{}
	case "3":
		return ThreeKeyCommand{}
	case "a":
		return AKeyCommand{}
	case "A":
		return UppercaseAKeyCommand{}
	case "w":
		return WKeyCommand{}
	case "-":
		return MinusKeyCommand{}
	case "+":
		return PlusKeyCommand{}
	case "g":
		return GKeyCommand{}
	case "D":
		return UppercaseDKeyCommand{}
	case "d":
		return DKeyCommand{}
	case "W":
		return UppercaseWKeyCommand{}
	case "s":
		return SKeyCommand{}
	case "S":
		return UppercaseSKeyCommand{}
	case "/":
		return SlashKeyCommand{}
	}
	return NilKeyCommand{}
}
