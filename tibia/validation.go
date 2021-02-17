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

	TibiaWorlds = []string{
		"Adra",
		"Antica",
		"Assombra",
		"Astera",
		"Belluma",
		"Belobra",
		"Bona",
		"Calmera",
		"Carnera",
		"Celebra",
		"Celesta",
		"Concorda",
		"Cosera",
		"Damora",
		"Descubra",
		"Dibra",
		"Duna",
		"Emera",
		"Endebra",
		"Endera",
		"Endura",
		"Epoca",
		"Estela",
		"Faluna",
		"Ferobra",
		"Firmera",
		"Funera",
		"Furia",
		"Garnera",
		"Gentebra",
		"Gladera",
		"Harmonia",
		"Helera",
		"Honbra",
		"Impera",
		"Inabra",
		"Javibra",
		"Jonera",
		"Kalibra",
		"Kenora",
		"Libertabra",
		"Lobera",
		"Luminera",
		"Lutabra",
		"Macabra",
		"Menera",
		"Mitigera",
		"Monza",
		"Nefera",
		"Noctera",
		"Nossobra",
		"Olera",
		"Ombra",
		"Pacembra",
		"Pacera",
		"Peloria",
		"Premia",
		"Pyra",
		"Quelibra",
		"Quintera",
		"Ragna",
		"Refugia",
		"Relania",
		"Relembra",
		"Secura",
		"Serdebra",
		"Serenebra",
		"Solidera",
		"Talera",
		"Torpera",
		"Tortura",
		"Unica",
		"Utobra",
		"Velocera",
		"Velocibra",
		"Velocita",
		"Venebra",
		"Vitia",
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
