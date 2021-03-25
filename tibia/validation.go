package tibia

import (
	"regexp"
	"strings"

	"emperror.dev/errors"
)

var (
	ErrLenZero     = errors.New("O nome fornecido não pode ser uma empty string")
	ErrSmallName   = errors.New("O nome fornecido é pequeno demais")
	ErrInvalidName = errors.New("O nome fornecido é inválido")
	NameRegex      = regexp.MustCompile(`[^\s'a-zA-Z\-\.]`)

	TibiaWorlds = [83]string{
		"Adra",
		"Antica",
		"Astera",
		"Belobra",
		"Bona",
		"Calmera",
		"Carnera",
		"Celebra",
		"Celesta",
		"Concorda",
		"Damora",
		"Descubra",
		"Dibra",
		"Emera",
		"Endebra",
		"Endera",
		"Endura",
		"Epoca",
		"Estela",
		"Fera",
		"Ferobra",
		"Fervora",
		"Firmera",
		"Garnera",
		"Gentebra",
		"Gladera",
		"Harmonia",
		"Honbra",
		"Impera",
		"Inabra",
		"Javibra",
		"Juva",
		"Kalibra",
		"Karna",
		"Kenora",
		"Libertabra",
		"Lobera",
		"Luminera",
		"Lutabra",
		"Menera",
		"Mercera",
		"Mitigera",
		"Monza",
		"Mudabra",
		"Nefera",
		"Nexa",
		"Nossobra",
		"Ombra",
		"Optera",
		"Pacembra",
		"Pacera",
		"Peloria",
		"Premia",
		"Quelibra",
		"Quintera",
		"Ragna",
		"Refugia",
		"Reinobra",
		"Relania",
		"Relembra",
		"Secura",
		"Serdebra",
		"Serenebra",
		"Solidera",
		"Talera",
		"Unica",
		"Unisera",
		"Utobra",
		"Velocera",
		"Velocibra",
		"Velocita",
		"Venebra",
		"Visabra",
		"Vunira",
		"Wintera",
		"Wizera",
		"Xandebra",
		"Xylona",
		"Yonabra",
		"Ysolera",
		"Zenobra",
		"Zuna",
		"Zunera",
	}
)

func validateName(name string) error {
	lenName := len(name)
	switch {
	case lenName <= 0:
		return ErrLenZero
	case lenName < 2:
		return ErrSmallName
	}

	matched := NameRegex.MatchString(name)
	if matched {
		return ErrInvalidName
	}

	return nil
}

func validWorld(world string) bool {
	for _, w := range TibiaWorlds {
		if strings.EqualFold(world, w) {
			return true
		}
	}

	return false
}
