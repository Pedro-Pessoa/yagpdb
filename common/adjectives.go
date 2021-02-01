// this package provides core functinality to yagpdb, important security stuff here
package common

// this is alphabetically sorted by the english words

var Adjectives = map[string][]string{
	"EN": {
		"abandoned",
		"able",
		"absolute",
		"academic",
		"acceptable",
		"acclaimed",
		"accomplished",
		"accurate",
		"aching",
		"acidic",
		"acrobatic",
		"active",
		"actual",
		"adept",
		"admirable",
		"admired",
		"adolescent",
		"adorable",
		"adored",
		"advanced",
		"adventurous",
		"affectionate",
		"afraid",
		"aged",
		"aggravating",
		"aggressive",
		"agile",
		"agitated",
		"agonizing",
		"agreeable",
		"ajar",
		"alarmed",
		"alarming",
		"alert",
		"alienated",
		"alive",
		"all",
		"altruistic",
		"amazing",
		"ambitious",
		"ample",
		"amused",
		"amusing",
		"anchored",
		"ancient",
		"angelic",
		"angry",
		"anguished",
		"animated",
		"annual",
		"another",
		"antique",
		"anxious",
		"any",
		"apprehensive",
		"appropriate",
		"apt",
		"arctic",
		"arid",
		"aromatic",
		"artistic",
		"ashamed",
		"assured",
		"astonishing",
		"athletic",
		"attached",
		"attentive",
		"attractive",
		"austere",
		"authentic",
		"authorized",
		"automatic",
		"avaricious",
		"average",
		"aware",
		"awesome",
		"awful",
		"awkward",
		"babyish",
		"back",
		"bad",
		"baggy",
		"bare",
		"barren",
		"basic",
		"beautiful",
		"belated",
		"beloved",
		"beneficial",
		"best",
		"better",
		"bewitched",
		"big",
		"big-hearted",
		"biodegradable",
		"bite-sized",
		"bitter",
		"black",
		"black-an-white",
		"bland",
		"blank",
		"blaring",
		"bleak",
		"blind",
		"blissful",
		"blond",
		"blue",
		"blushing",
		"bogus",
		"boiling",
		"bold",
		"bony",
		"boring",
		"bossy",
		"both",
		"bouncy",
		"bountiful",
		"bowed",
		"brave",
		"breakable",
		"brief",
		"bright",
		"brilliant",
		"brisk",
		"broken",
		"bronze",
		"brown",
		"bruised",
		"bubbly",
		"bulky",
		"bumpy",
		"buoyant",
		"burdensome",
		"burly",
		"bustling",
		"busy",
		"buttery",
		"buzzing",
		"calculating",
		"calm",
		"candid",
		"canine",
		"capital",
		"carefree",
		"careful",
		"careless",
		"caring",
		"cautious",
		"cavernous",
		"celebrated",
		"charming",
		"cheap",
		"cheerful",
		"cheery",
		"chief",
		"chilly",
		"chubby",
		"circular",
		"classic",
		"clean",
		"clear",
		"clear-cut",
		"clever",
		"close",
		"closed",
		"cloudy",
		"clueless",
		"clumsy",
		"cluttered",
		"coarse",
		"cold",
		"colorful",
		"colorless",
		"colossal",
		"comfortable",
		"common",
		"compassionate",
		"competent",
		"complete",
		"complex",
		"complicated",
		"composed",
		"concerned",
		"concrete",
		"confused",
		"conscious",
		"considerate",
		"constant",
		"content",
		"conventional",
		"cooked",
		"cool",
		"cooperative",
		"coordinated",
		"corny",
		"corrupt",
		"costly",
		"courageous",
		"courteous",
		"crafty",
		"crazy",
		"creamy",
		"creative",
		"creepy",
		"criminal",
		"crisp",
		"critical",
		"crooked",
		"crowded",
		"cruel",
		"crushing",
		"cuddly",
		"cultivated",
		"cultured",
		"cumbersome",
		"curly",
		"curvy",
		"cute",
		"cylindrical",
		"damaged",
		"damp",
		"dangerous",
		"dapper",
		"daring",
		"dark",
		"darling",
		"dazzling",
		"dead",
		"deadly",
		"deafening",
		"dear",
		"dearest",
		"decent",
		"decimal",
		"decisive",
		"deep",
		"defenseless",
		"defensive",
		"defiant",
		"deficient",
		"definite",
		"definitive",
		"delayed",
		"delectable",
		"delicious",
		"delightful",
		"delirious",
		"demanding",
		"dense",
		"dental",
		"dependable",
		"dependent",
		"descriptive",
		"deserted",
		"detailed",
		"determined",
		"devoted",
		"different",
		"difficult",
		"digital",
		"diligent",
		"dim",
		"dimpled",
		"dimwitted",
		"direct",
		"dirty",
		"disastrous",
		"discrete",
		"disfigured",
		"disguised",
		"disgusting",
		"dishonest",
		"disloyal",
		"dismal",
		"distant",
		"distinct",
		"distorted",
		"dizzy",
		"dopey",
		"doting",
		"double",
		"downright",
		"drab",
		"drafty",
		"dramatic",
		"dreary",
		"droopy",
		"dry",
		"dual",
		"dull",
		"dutiful",
		"each",
		"eager",
		"early",
		"earnest",
		"easy",
		"easy-going",
		"ecstatic",
		"edible",
		"educated",
		"elaborate",
		"elastic",
		"elated",
		"elderly",
		"electric",
		"elegant",
		"elementary",
		"elliptical",
		"embarrassed",
		"embellished",
		"eminent",
		"emotional",
		"empty",
		"enchanted",
		"enchanting",
		"energetic",
		"enlightened",
		"enormous",
		"enraged",
		"entire",
		"envious",
		"equal",
		"equatorial",
		"essential",
		"esteemed",
		"ethical",
		"euphoric",
		"even",
		"evergreen",
		"everlasting",
		"every",
		"evil",
		"exalted",
		"excellent",
		"excitable",
		"excited",
		"exciting",
		"exemplary",
		"exhausted",
		"exotic",
		"expensive",
		"experienced",
		"expert",
		"extra-large",
		"extra-small",
		"extraneous",
		"extroverted",
		"fabulous",
		"failing",
		"faint",
		"fair",
		"faithful",
		"fake",
		"false",
		"familiar",
		"famous",
		"fancy",
		"fantastic",
		"far",
		"far-flung",
		"far-off",
		"faraway",
		"fast",
		"fat",
		"fatal",
		"fatherly",
		"favorable",
		"favorite",
		"fearful",
		"fearless",
		"feisty",
		"feline",
		"female",
		"feminine",
		"few",
		"fickle",
		"filthy",
		"fine",
		"finished",
		"firm",
		"first",
		"firsthand",
		"fitting",
		"fixed",
		"flaky",
		"flamboyant",
		"flashy",
		"flat",
		"flawed",
		"flawless",
		"flickering",
		"flimsy",
		"flippant",
		"flowery",
		"fluffy",
		"fluid",
		"flustered",
		"focused",
		"fond",
		"foolhardy",
		"foolish",
		"forceful",
		"forked",
		"formal",
		"forsaken",
		"forthright",
		"fortunate",
		"fragrant",
		"frail",
		"frank",
		"frayed",
		"free",
		"french",
		"frequent",
		"fresh",
		"friendly",
		"frightened",
		"frightening",
		"frigid",
		"frilly",
		"frivolous",
		"frizzy",
		"front",
		"frosty",
		"frozen",
		"frugal",
		"fruitful",
		"full",
		"fumbling",
		"functional",
		"funny",
		"fussy",
		"fuzzy",
		"gargantuan",
		"gaseous",
		"general",
		"generous",
		"gentle",
		"genuine",
		"giant",
		"giddy",
		"gifted",
		"gigantic",
		"giving",
		"glamorous",
		"glaring",
		"glass",
		"gleaming",
		"gleeful",
		"glistening",
		"glittering",
		"gloomy",
		"glorious",
		"glossy",
		"glum",
		"golden",
		"good",
		"good-natured",
		"gorgeous",
		"graceful",
		"gracious",
		"grand",
		"grandiose",
		"granular",
		"grateful",
		"grave",
		"gray",
		"great",
		"greedy",
		"green",
		"gregarious",
		"grim",
		"grimy",
		"gripping",
		"grizzled",
		"gross",
		"grotesque",
		"grouchy",
		"grounded",
		"growing",
		"growling",
		"grown",
		"grubby",
		"gruesome",
		"grumpy",
		"guilty",
		"gullible",
		"gummy",
		"hairy",
		"half",
		"handmade",
		"handsome",
		"handy",
		"happy",
		"happy-g-lucky",
		"hard",
		"hard-t-find",
		"harmful",
		"harmless",
		"harmonious",
		"harsh",
		"hasty",
		"hateful",
		"haunting",
		"healthy",
		"heartfelt",
		"hearty",
		"heavenly",
		"heavy",
		"hefty",
		"helpful",
		"helpless",
		"hidden",
		"hideous",
		"high",
		"high-level",
		"hilarious",
		"hoarse",
		"hollow",
		"homely",
		"honest",
		"honorable",
		"honored",
		"hopeful",
		"horrible",
		"hospitable",
		"hot",
		"huge",
		"humble",
		"humiliating",
		"humming",
		"humongous",
		"hungry",
		"hurtful",
		"husky",
		"icky",
		"icy",
		"ideal",
		"idealistic",
		"identical",
		"idiotic",
		"idle",
		"idolized",
		"ignorant",
		"ill",
		"ill-fated",
		"ill-informed",
		"illegal",
		"illiterate",
		"illustrious",
		"imaginary",
		"imaginative",
		"immaculate",
		"immaterial",
		"immediate",
		"immense",
		"impartial",
		"impassioned",
		"impeccable",
		"imperfect",
		"imperturbable",
		"impish",
		"impolite",
		"important",
		"impossible",
		"impractical",
		"impressionable",
		"impressive",
		"improbable",
		"impure",
		"inborn",
		"incomparable",
		"incompatible",
		"incomplete",
		"inconsequential",
		"incredible",
		"indelible",
		"indolent",
		"inexperienced",
		"infamous",
		"infantile",
		"infatuated",
		"inferior",
		"infinite",
		"informal",
		"innocent",
		"insecure",
		"insidious",
		"insignificant",
		"insistent",
		"instructive",
		"insubstantial",
		"intelligent",
		"intent",
		"intentional",
		"interesting",
		"internal",
		"international",
		"intrepid",
		"ironclad",
		"irresponsible",
		"irritating",
		"itchy",
		"jaded",
		"jagged",
		"jam-packed",
		"jaunty",
		"jealous",
		"jittery",
		"joint",
		"jolly",
		"jovial",
		"joyful",
		"joyous",
		"jubilant",
		"judicious",
		"juicy",
		"jumbo",
		"jumpy",
		"junior",
		"juvenile",
		"kaleidoscopic",
		"keen",
		"key",
		"kind",
		"kindhearted",
		"kindly",
		"klutzy",
		"knobby",
		"knotty",
		"knowing",
		"knowledgeable",
		"known",
		"kooky",
		"kosher",
		"lame",
		"lanky",
		"large",
		"last",
		"lasting",
		"late",
		"lavish",
		"lawful",
		"lazy",
		"leading",
		"leafy",
		"lean",
		"left",
		"legal",
		"legitimate",
		"light",
		"lighthearted",
		"likable",
		"likely",
		"limited",
		"limp",
		"limping",
		"linear",
		"lined",
		"liquid",
		"little",
		"live",
		"lively",
		"livid",
		"loathsome",
		"lone",
		"lonely",
		"long",
		"long-term",
		"loose",
		"lopsided",
		"lost",
		"loud",
		"lovable",
		"lovely",
		"loving",
		"low",
		"loyal",
		"lucky",
		"lumbering",
		"luminous",
		"lumpy",
		"lustrous",
		"luxurious",
		"mad",
		"made-up",
		"magnificent",
		"majestic",
		"major",
		"male",
		"mammoth",
		"married",
		"marvelous",
		"masculine",
		"massive",
		"mature",
		"meager",
		"mealy",
		"mean",
		"measly",
		"meaty",
		"medical",
		"mediocre",
		"medium",
		"meek",
		"mellow",
		"melodic",
		"memorable",
		"menacing",
		"merry",
		"messy",
		"metallic",
		"mild",
		"milky",
		"mindless",
		"miniature",
		"minor",
		"minty",
		"miserable",
		"miserly",
		"misguided",
		"misty",
		"mixed",
		"modern",
		"modest",
		"moist",
		"monstrous",
		"monthly",
		"monumental",
		"moral",
		"mortified",
		"motherly",
		"motionless",
		"mountainous",
		"muddy",
		"muffled",
		"multicolored",
		"mundane",
		"murky",
		"mushy",
		"musty",
		"muted",
		"mysterious",
		"naive",
		"narrow",
		"nasty",
		"natural",
		"naughty",
		"nautical",
		"near",
		"neat",
		"necessary",
		"needy",
		"negative",
		"neglected",
		"negligible",
		"neighboring",
		"nervous",
		"new",
		"next",
		"nice",
		"nifty",
		"nimble",
		"nippy",
		"nocturnal",
		"noisy",
		"nonstop",
		"normal",
		"notable",
		"noted",
		"noteworthy",
		"novel",
		"noxious",
		"numb",
		"nutritious",
		"nutty",
		"obedient",
		"obese",
		"oblong",
		"obvious",
		"occasional",
		"odd",
		"oddball",
		"offbeat",
		"offensive",
		"official",
		"oily",
		"old",
		"old-fashioned",
		"only",
		"open",
		"optimal",
		"optimistic",
		"opulent",
		"orange",
		"orderly",
		"ordinary",
		"organic",
		"original",
		"ornate",
		"ornery",
		"other",
		"our",
		"outgoing",
		"outlandish",
		"outlying",
		"outrageous",
		"outstanding",
		"oval",
		"overcooked",
		"overdue",
		"overjoyed",
		"overlooked",
		"palatable",
		"pale",
		"paltry",
		"parallel",
		"parched",
		"partial",
		"passionate",
		"past",
		"pastel",
		"peaceful",
		"peppery",
		"perfect",
		"perfumed",
		"periodic",
		"perky",
		"personal",
		"pertinent",
		"pesky",
		"pessimistic",
		"petty",
		"phony",
		"physical",
		"piercing",
		"pink",
		"pitiful",
		"plain",
		"plaintive",
		"plastic",
		"playful",
		"pleasant",
		"pleased",
		"pleasing",
		"plump",
		"plush",
		"pointed",
		"pointless",
		"poised",
		"polished",
		"polite",
		"political",
		"poor",
		"popular",
		"portly",
		"posh",
		"positive",
		"possible",
		"potable",
		"powerful",
		"powerless",
		"practical",
		"precious",
		"present",
		"prestigious",
		"pretty",
		"previous",
		"pricey",
		"prickly",
		"primary",
		"prime",
		"pristine",
		"private",
		"prize",
		"probable",
		"productive",
		"profitable",
		"profuse",
		"proper",
		"proud",
		"prudent",
		"punctual",
		"pungent",
		"puny",
		"pure",
		"purple",
		"pushy",
		"putrid",
		"puzzled",
		"puzzling",
		"quaint",
		"qualified",
		"quarrelsome",
		"quarterly",
		"queasy",
		"querulous",
		"questionable",
		"quick",
		"quick-witted",
		"quiet",
		"quintessential",
		"quirky",
		"quixotic",
		"quizzical",
		"radiant",
		"ragged",
		"rapid",
		"rare",
		"rash",
		"raw",
		"ready",
		"real",
		"realistic",
		"reasonable",
		"recent",
		"reckless",
		"rectangular",
		"red",
		"reflecting",
		"regal",
		"regular",
		"reliable",
		"relieved",
		"remarkable",
		"remorseful",
		"remote",
		"repentant",
		"repulsive",
		"required",
		"respectful",
		"responsible",
		"revolving",
		"rewarding",
		"rich",
		"right",
		"rigid",
		"ringed",
		"ripe",
		"roasted",
		"robust",
		"rosy",
		"rotating",
		"rotten",
		"rough",
		"round",
		"rowdy",
		"royal",
		"rubbery",
		"ruddy",
		"rude",
		"rundown",
		"runny",
		"rural",
		"rusty",
		"sad",
		"safe",
		"salty",
		"same",
		"sandy",
		"sane",
		"sarcastic",
		"sardonic",
		"satisfied",
		"scaly",
		"scarce",
		"scared",
		"scary",
		"scented",
		"scholarly",
		"scientific",
		"scornful",
		"scratchy",
		"scrawny",
		"second",
		"second-hand",
		"secondary",
		"secret",
		"self-assured",
		"self-reliant",
		"selfish",
		"sentimental",
		"separate",
		"serene",
		"serious",
		"serpentine",
		"several",
		"severe",
		"sexy",
		"shabby",
		"shadowy",
		"shady",
		"shallow",
		"shameful",
		"shameless",
		"sharp",
		"shimmering",
		"shiny",
		"shocked",
		"shocking",
		"shoddy",
		"short",
		"short-term",
		"showy",
		"shrill",
		"shy",
		"sick",
		"silent",
		"silky",
		"silly",
		"silver",
		"similar",
		"simple",
		"simplistic",
		"sinful",
		"single",
		"sizzling",
		"skeletal",
		"skinny",
		"sleepy",
		"slight",
		"slim",
		"slimy",
		"slippery",
		"slow",
		"slushy",
		"small",
		"smart",
		"smoggy",
		"smooth",
		"smug",
		"snappy",
		"snarling",
		"sneaky",
		"sniveling",
		"snoopy",
		"sociable",
		"soft",
		"soggy",
		"solid",
		"somber",
		"some",
		"sophisticated",
		"sore",
		"sorrowful",
		"soulful",
		"soupy",
		"sour",
		"spanish",
		"sparkling",
		"sparse",
		"specific",
		"spectacular",
		"speedy",
		"spherical",
		"spicy",
		"spiffy",
		"spirited",
		"spiteful",
		"splendid",
		"spotless",
		"spotted",
		"spry",
		"square",
		"squeaky",
		"squiggly",
		"stable",
		"staid",
		"stained",
		"stale",
		"standard",
		"starchy",
		"stark",
		"starry",
		"steel",
		"steep",
		"sticky",
		"stiff",
		"stimulating",
		"stingy",
		"stormy",
		"straight",
		"strange",
		"strict",
		"strident",
		"striking",
		"striped",
		"strong",
		"studious",
		"stunning",
		"stupendous",
		"stupid",
		"sturdy",
		"stylish",
		"subdued",
		"submissive",
		"substantial",
		"subtle",
		"suburban",
		"sudden",
		"sugary",
		"sunny",
		"super",
		"superb",
		"superficial",
		"superior",
		"supportive",
		"surprised",
		"suspicious",
		"svelte",
		"sweaty",
		"sweet",
		"sweltering",
		"swift",
		"sympathetic",
		"talkative",
		"tall",
		"tame",
		"tan",
		"tangible",
		"tart",
		"tasty",
		"tattered",
		"taut",
		"tedious",
		"teeming",
		"tempting",
		"tender",
		"tense",
		"tepid",
		"terrible",
		"terrific",
		"testy",
		"thankful",
		"that",
		"these",
		"thick",
		"thin",
		"third",
		"thirsty",
		"this",
		"thorny",
		"thorough",
		"those",
		"thoughtful",
		"threadbare",
		"thrifty",
		"thunderous",
		"tidy",
		"tight",
		"timely",
		"tinted",
		"tiny",
		"tired",
		"torn",
		"total",
		"tough",
		"tragic",
		"trained",
		"traumatic",
		"treasured",
		"tremendous",
		"triangular",
		"tricky",
		"trifling",
		"trim",
		"trivial",
		"troubled",
		"true",
		"trusting",
		"trustworthy",
		"trusty",
		"truthful",
		"tubby",
		"turbulent",
		"twin",
		"ugly",
		"ultimate",
		"unacceptable",
		"unaware",
		"uncomfortable",
		"uncommon",
		"unconscious",
		"understated",
		"unequaled",
		"uneven",
		"unfinished",
		"unfit",
		"unfolded",
		"unfortunate",
		"unhappy",
		"unhealthy",
		"uniform",
		"unimportant",
		"unique",
		"united",
		"unkempt",
		"unknown",
		"unlawful",
		"unlined",
		"unlucky",
		"unnatural",
		"unpleasant",
		"unrealistic",
		"unripe",
		"unruly",
		"unselfish",
		"unsightly",
		"unsteady",
		"unsung",
		"untidy",
		"untimely",
		"untried",
		"untrue",
		"unused",
		"unusual",
		"unwelcome",
		"unwieldy",
		"unwilling",
		"unwitting",
		"unwritten",
		"upbeat",
		"upright",
		"upset",
		"urban",
		"usable",
		"used",
		"useful",
		"useless",
		"utilized",
		"utter",
		"vacant",
		"vague",
		"vain",
		"valid",
		"valuable",
		"vapid",
		"variable",
		"vast",
		"velvety",
		"venerated",
		"vengeful",
		"verifiable",
		"vibrant",
		"vicious",
		"victorious",
		"vigilant",
		"vigorous",
		"villainous",
		"violent",
		"violet",
		"virtual",
		"virtuous",
		"visible",
		"vital",
		"vivacious",
		"vivid",
		"voluminous",
		"wan",
		"warlike",
		"warm",
		"warmhearted",
		"warped",
		"wary",
		"wasteful",
		"watchful",
		"waterlogged",
		"watery",
		"wavy",
		"weak",
		"wealthy",
		"weary",
		"webbed",
		"wee",
		"weekly",
		"weepy",
		"weighty",
		"weird",
		"welcome",
		"well-documented",
		"well-groomed",
		"well-informed",
		"well-lit",
		"well-made",
		"well-off",
		"well-worn",
		"wet",
		"which",
		"whimsical",
		"whirlwind",
		"whispered",
		"white",
		"whole",
		"whopping",
		"wicked",
		"wide",
		"wide-eyed",
		"wiggly",
		"wild",
		"willing",
		"wilted",
		"winding",
		"windy",
		"winged",
		"wiry",
		"wise",
		"witty",
		"wobbly",
		"woeful",
		"wonderful",
		"wooden",
		"woozy",
		"wordy",
		"worldly",
		"worn",
		"worried",
		"worrisome",
		"worse",
		"worst",
		"worthless",
		"worthwhile",
		"worthy",
		"wrathful",
		"wretched",
		"writhing",
		"wrong",
		"wry",
		"yawning",
		"yearly",
		"yellow",
		"yellowish",
		"young",
		"youthful",
		"yummy",
		"zany",
		"zealous",
		"zesty",
		"zigzag",
	},
	"PT": {
		"abandonado",
		"capaz",
		"absoluto",
		"acadêmico",
		"aceitável",
		"aclamado",
		"realizado",
		"preciso",
		"dolorido",
		"ácida",
		"acrobático",
		"ativo",
		"atual",
		"adepto",
		"admirável",
		"admirado",
		"adolescente",
		"adorável",
		"adorada",
		"avançado",
		"aventureiro",
		"afectuoso",
		"com-medo",
		"envelhecido",
		"agravante",
		"agressivo",
		"ágil",
		"agitado",
		"agonizante",
		"agradáveis",
		"entreaberto",
		"alarmado",
		"alarmante",
		"alerta",
		"alienado",
		"vivo",
		"todos",
		"altruísta",
		"espantoso",
		"ambicioso",
		"ampla",
		"entretido",
		"divertida",
		"ancorado",
		"antigo",
		"angelical",
		"zangado",
		"angustiado",
		"animado",
		"anual",
		"outro",
		"antiguidade",
		"ansioso",
		"qualquer",
		"apreensivo",
		"apropriado",
		"aptos",
		"ártico",
		"árido",
		"aromático",
		"artística",
		"envergonhados",
		"assegurado",
		"espantoso",
		"atlético",
		"em-anexo",
		"atento",
		"atraente",
		"austero",
		"autêntico",
		"autorizado",
		"automático",
		"avarento",
		"médio",
		"consciente",
		"espectacular",
		"horrível",
		"embaraçoso",
		"infantil",
		"de-volta",
		"mau",
		"folgado",
		"nu",
		"árido",
		"básico",
		"belo",
		"atrasada",
		"amado",
		"benéfico",
		"melhor",
		"melhor-ainda",
		"enfeitiçado",
		"grande",
		"de-coração-grande",
		"biodegradável",
		"pequeno",
		"amargo",
		"preto",
		"preto-e-branco",
		"suave",
		"em-branco",
		"estridente",
		"sombrio",
		"cego",
		"bem-aventurado",
		"loira",
		"azul",
		"ruborizar",
		"falso",
		"fervendo",
		"arrojado",
		"ossudo",
		"aborrecido",
		"mandão",
		"ambos",
		"saltitante",
		"abundante",
		"curvado",
		"corajoso",
		"quebrável",
		"breve",
		"brilhante",
		"brilhantíssimo",
		"rápido",
		"quebrado",
		"bronze",
		"castanho",
		"ferido",
		"borbulhante",
		"volumoso",
		"acidentado",
		"flutuante",
		"oneroso",
		"corpulento",
		"movimentado",
		"ocupado",
		"amanteigado",
		"zumbido",
		"cálculo",
		"calmo",
		"cândida",
		"canino",
		"capital",
		"despreocupado",
		"cuidadoso",
		"descuidado",
		"cuidadosa",
		"cauteloso",
		"cavernoso",
		"celebrado",
		"encantador",
		"barato",
		"alegrado",
		"afortunoso",
		"chefe",
		"friorento",
		"rechonchudo",
		"circular",
		"clássico",
		"limpo",
		"claro",
		"nítida",
		"esperto",
		"fechadíssimo",
		"fechado",
		"nublado",
		"sem-pistas",
		"desajeitado",
		"atulhado",
		"grosseiro",
		"frio",
		"colorido",
		"incolor",
		"colossal",
		"confortável",
		"comum",
		"compassivo",
		"competente",
		"completo",
		"complexo",
		"complicado",
		"composto",
		"preocupado",
		"betão",
		"confuso",
		"consciente",
		"atencioso",
		"constante",
		"conteúdo",
		"convencional",
		"cozinhado",
		"legal",
		"cooperativa",
		"coordenado",
		"piroso",
		"corrupto",
		"dispendioso",
		"corajoso",
		"cortês",
		"astucioso",
		"louco",
		"cremoso",
		"criativo",
		"arrepiante",
		"criminoso",
		"estaladiço",
		"crítico",
		"tortuoso",
		"apinhado",
		"cruel",
		"esmagamento",
		"fofinho",
		"cultivado",
		"culto",
		"incômodo",
		"encaracolado",
		"curvado",
		"giro",
		"cilíndrico",
		"danificado",
		"húmido",
		"perigoso",
		"dapper",
		"ousadia",
		"escuro",
		"querida",
		"deslumbrante",
		"morto",
		"mortífero",
		"ensurdecedor",
		"querido",
		"decente",
		"decente",
		"decimal",
		"decisivo",
		"profundo",
		"indefeso",
		"defensiva",
		"desafiante",
		"deficiente",
		"definitivo",
		"definitivo",
		"atrasado",
		"deleitável",
		"delicioso",
		"encantador",
		"delirante",
		"exigente",
		"denso",
		"dental",
		"de-confiança",
		"dependente",
		"descritiva",
		"deserto",
		"detalhado",
		"determinado",
		"dedicado",
		"diferente",
		"difícil",
		"digital",
		"diligente",
		"escuro",
		"covarde",
		"grosseiro",
		"direto",
		"sujo",
		"desastroso",
		"discreta",
		"desfigurado",
		"disfarçado",
		"nojento",
		"desonesto",
		"desleal",
		"desanimador",
		"distante",
		"distinto",
		"distorcido",
		"tonta",
		"drogado",
		"amoroso",
		"duplo",
		"totalmente",
		"monótono",
		"rascunho",
		"dramático",
		"entediante",
		"caído",
		"seco",
		"duplo",
		"tedioso",
		"obediente",
		"cada",
		"ansioso",
		"cedo",
		"sincero",
		"fácil",
		"facílimo",
		"extático",
		"comestível",
		"educado",
		"elaborado",
		"elástico",
		"eufórico",
		"idoso",
		"elétrico",
		"elegante",
		"elementar",
		"elíptica",
		"envergonhado",
		"embelezado",
		"eminente",
		"emocional",
		"vazio",
		"encantado",
		"encantador",
		"enérgico",
		"iluminado",
		"enorme",
		"enfurecido",
		"inteiro",
		"invejoso",
		"igual",
		"equatorial",
		"essencial",
		"estimado",
		"ético",
		"eufórico",
		"par",
		"perene",
		"eterno",
		"todo",
		"mal",
		"exaltado",
		"excelente",
		"excitável",
		"excitado",
		"excitante",
		"exemplar",
		"exausto",
		"exótico",
		"caro",
		"experiente",
		"especialista",
		"extra-grande",
		"extra-pequeno",
		"estranho",
		"extrovertida",
		"fabuloso",
		"falho",
		"tênue",
		"justo",
		"fiel",
		"aleivoso",
		"falsificado",
		"familiar",
		"famoso",
		"extravagante",
		"fantástico",
		"longe",
		"longínquo",
		"distante",
		"afastado",
		"rápido",
		"gordo",
		"fatal",
		"paternal",
		"favorável",
		"favorito",
		"temeroso",
		"sem-medo",
		"combativo",
		"felino",
		"feminino",
		"feminal",
		"pouco",
		"inconstante",
		"imundo",
		"multado",
		"terminado",
		"firme",
		"primeiro",
		"em-primeira-mão",
		"encaixe",
		"fixo",
		"escamoso",
		"flambado",
		"cintilante",
		"planejador",
		"defeituoso",
		"irrepreensível",
		"cintilante",
		"frágil",
		"irreverente",
		"florido",
		"fofo",
		"fluido",
		"nervoso",
		"focalizado",
		"afeiçoado",
		"imprudente",
		"tolo",
		"enérgico",
		"bifurcado",
		"formal",
		"abandonado",
		"directo",
		"afortunado",
		"perfumado",
		"frágil",
		"franco",
		"desgastado",
		"livre",
		"francês",
		"frequente",
		"fresco",
		"amigável",
		"assustado",
		"assustador",
		"frígido",
		"de-babado",
		"frívolo",
		"crespo",
		"frente",
		"gelado",
		"congelado",
		"frugal",
		"frutuosa",
		"cheio",
		"fumegante",
		"funcional",
		"engraçado",
		"picuinhas",
		"difusa",
		"gigantesco",
		"gasoso",
		"geral",
		"generoso",
		"gentil",
		"genuíno",
		"gigante",
		"tonto",
		"dotado",
		"gigantesco",
		"doador",
		"glamoroso",
		"gritante",
		"quebradiço",
		"cintilante",
		"contente",
		"reluzente",
		"cintilante",
		"sombrio",
		"glorioso",
		"brilhante",
		"taciturno",
		"dourado",
		"bom",
		"bem-disposto",
		"deslumbrante",
		"gracioso",
		"aliciante",
		"grandioso",
		"epopeico",
		"granular",
		"grato",
		"falecido",
		"cinzento",
		"grande",
		"ganancioso",
		"verde",
		"gregário",
		"sombrio",
		"sujo",
		"apaixonante",
		"cinzento",
		"bruto",
		"grotesco",
		"ranzinza",
		"fundamentado",
		"em-crescimento",
		"rosnado",
		"cultivado",
		"sujo",
		"horripilante",
		"mal-humorado",
		"culpado",
		"ingénuo",
		"grudento",
		"peludo",
		"metade",
		"feito-à-mão",
		"bonitão",
		"à-mão",
		"feliz",
		"sortudo",
		"duro",
		"difícil-de-encontrar",
		"prejudicial",
		"inofensivo",
		"harmonioso",
		"áspero",
		"apressado",
		"odioso",
		"assombroso",
		"saudável",
		"sincero",
		"caloroso",
		"celestial",
		"pesado",
		"insultuoso",
		"útil",
		"indefeso",
		"escondido",
		"hediondo",
		"elevado",
		"de-alto-nível",
		"hilariante",
		"rouco",
		"oco",
		"caseiro",
		"honesto",
		"honorável",
		"honrado",
		"esperançoso",
		"horrível",
		"hospitaleiro",
		"quente",
		"enorme",
		"humilde",
		"humilhante",
		"zumbido",
		"encorpado",
		"faminto",
		"doloroso",
		"rouco",
		"nojento",
		"gelado",
		"ideal",
		"idealista",
		"idêntico",
		"idiota",
		"ocioso",
		"idolatrado",
		"ignorante",
		"doente",
		"malfadado",
		"mal-informado",
		"ilegal",
		"analfabeto",
		"ilustre",
		"imaginário",
		"imaginativo",
		"imaculado",
		"imaterial",
		"imediato",
		"imenso",
		"imparcial",
		"apaixonado",
		"impecável",
		"imperfeito",
		"imperturbável",
		"travesso",
		"indelicado",
		"importante",
		"impossível",
		"impraticável",
		"impressionável",
		"impressionante",
		"improvável",
		"impuro",
		"inato",
		"incomparável",
		"incompatível",
		"incompleto",
		"inconsequente",
		"incrível",
		"indelével",
		"indolente",
		"inexperiente",
		"infame",
		"infantil",
		"apaixonado",
		"inferior",
		"infinito",
		"informal",
		"inocente",
		"insegura",
		"insidioso",
		"insignificante",
		"insistente",
		"instrutivo",
		"insubstancial",
		"inteligente",
		"intenção",
		"intencional",
		"interessante",
		"interno",
		"internacional",
		"intrépido",
		"revestido-a-ferro",
		"irresponsável",
		"irritante",
		"coceira",
		"exausto",
		"recortado",
		"encravado",
		"alegre",
		"ciumento",
		"trêmulo",
		"conjunto",
		"bem-aventurado",
		"jovial",
		"bem-afortunado",
		"afortunado",
		"de-júbilo",
		"judicioso",
		"suculento",
		"jumbo",
		"saltitante",
		"júnior",
		"juvenil",
		"caleidoscópico",
		"entusiasta",
		"chavoso",
		"bondoso",
		"de-bom-coração",
		"amável",
		"desastrada",
		"saliente",
		"nodoso",
		"sabido",
		"conhecedor",
		"conhecido",
		"excêntrico",
		"kosher",
		"coxo",
		"magro",
		"grande",
		"último",
		"duradouro",
		"tarde",
		"pródigo",
		"legal",
		"preguiçoso",
		"líder",
		"folhoso",
		"magra",
		"esquerda",
		"legal",
		"legítimo",
		"luz",
		"despreocupado",
		"simpático",
		"provável",
		"limitado",
		"coxear",
		"a-coxear",
		"linear",
		"forrado",
		"líquido",
		"pouco",
		"ao-vivo",
		"animado",
		"lívido",
		"repugnante",
		"solitário",
		"sozinho",
		"longo",
		"a-longo-prazo",
		"solto",
		"desequilibrado",
		"perdido",
		"alto",
		"adorável",
		"adorável",
		"amoroso",
		"baixo",
		"leal",
		"sortuda",
		"serração",
		"luminosa",
		"grumoso",
		"lustroso",
		"luxuoso",
		"louco",
		"maquilhado",
		"magnífico",
		"majestático",
		"maior",
		"masculino",
		"mamute",
		"casado",
		"maravilhoso",
		"masculino",
		"massivo",
		"maduro",
		"parco",
		"farinhento",
		"mau",
		"miserável",
		"carnudo",
		"médico",
		"medíocre",
		"médio",
		"manso",
		"suave",
		"melódico",
		"memorável",
		"ameaçador",
		"bendito",
		"desarrumado",
		"metálico",
		"suave",
		"leitosa",
		"sem-sentido",
		"miniatura",
		"menor",
		"menta",
		"miserável",
		"miseravelmente",
		"mal-orientado",
		"nebuloso",
		"misto",
		"moderno",
		"modesto",
		"húmido",
		"monstruoso",
		"mensalmente",
		"monumental",
		"moral",
		"mortificado",
		"maternal",
		"imóvel",
		"montanhoso",
		"lamacento",
		"abafado",
		"multicolorido",
		"mundano",
		"obscuro",
		"pastoso",
		"bafiento",
		"mudo",
		"misterioso",
		"ingênuo",
		"estreito",
		"desagradável",
		"natural",
		"maroto",
		"náutico",
		"perto",
		"puro",
		"necessário",
		"carenciado",
		"negativo",
		"negligenciado",
		"insignificante",
		"vizinho",
		"nervoso",
		"novo",
		"seguinte",
		"simpático",
		"elegante",
		"ágil",
		"frio",
		"noturno",
		"ruidoso",
		"sem-parar",
		"normal",
		"notável",
		"anotado",
		"digno-de-nota",
		"romântico",
		"nocivo",
		"entorpecido",
		"nutritivo",
		"louco",
		"resignado",
		"obeso",
		"oblongo",
		"óbvio",
		"ocasional",
		"estranho",
		"estapafúrdio",
		"excêntrico",
		"ofensivo",
		"oficial",
		"oleoso",
		"antigo",
		"antiquado",
		"apenas",
		"aberto",
		"ótimo",
		"otimista",
		"opulento",
		"laranja",
		"ordeiro",
		"ordinário",
		"orgânico",
		"original",
		"ornamentado",
		"ornamental",
		"outro",
		"nosso",
		"extrovertido",
		"estranho",
		"periférico",
		"ultrajante",
		"notável",
		"oval",
		"cozido-em-demasia",
		"atrasado",
		"muito-contente",
		"negligenciado",
		"palatável",
		"pálido",
		"miserável",
		"paralelo",
		"ressecado",
		"parcial",
		"apaixonado",
		"passado",
		"pastel",
		"pacífico",
		"apimentado",
		"perfeito",
		"perfumado",
		"periódico",
		"animado",
		"pessoal",
		"pertinente",
		"irritante",
		"pessimista",
		"mesquinho",
		"fingido",
		"físico",
		"piercing",
		"rosa",
		"lastimosa",
		"simples",
		"queixoso",
		"plástico",
		"brincalhão",
		"agradável",
		"satisfeito",
		"agradável",
		"roliço",
		"pelúcia",
		"pontiagudo",
		"inútil",
		"em-posição",
		"polido",
		"educado",
		"político",
		"pobre",
		"popular",
		"adiposo",
		"chique",
		"positivo",
		"possível",
		"potável",
		"poderoso",
		"impotente",
		"prático",
		"precioso",
		"presente",
		"prestigioso",
		"bonito",
		"anterior",
		"caro",
		"picadinho",
		"primária",
		"primo",
		"imaculado",
		"privado",
		"prêmio",
		"provável",
		"produtivo",
		"rentável",
		"profuso",
		"próprio",
		"orgulhoso",
		"prudente",
		"pontual",
		"pungente",
		"insignificante",
		"puro",
		"púrpura",
		"insistente",
		"pútrido",
		"intrigado",
		"intrigante",
		"pitoresco",
		"qualificado",
		"briguento",
		"trimestralmente",
		"enjoado",
		"guerrento",
		"questionável",
		"rápido",
		"perspicaz",
		"sossegado",
		"quintessencial",
		"peculiar",
		"quixotesco",
		"curioso",
		"radiante",
		"esfarrapado",
		"rápido",
		"raro",
		"erupção",
		"em-bruto",
		"pronto",
		"real",
		"realista",
		"razoável",
		"recente",
		"imprudente",
		"retangular",
		"vermelho",
		"refletor",
		"régio",
		"regular",
		"de-confiança",
		"aliviado",
		"notável",
		"com-remorso",
		"remoto",
		"arrependido",
		"repulsivo",
		"necessário",
		"respeitoso",
		"responsável",
		"rotativo",
		"gratificante",
		"rico",
		"direito",
		"rígido",
		"anelado",
		"maduro",
		"assado",
		"robusto",
		"cor-de-rosa",
		"rotativo",
		"apodrecido",
		"áspero",
		"redondo",
		"bagunceiro",
		"real",
		"emborrachado",
		"vermelho",
		"indelicado",
		"atropelado",
		"escorregadio",
		"rural",
		"enferrujado",
		"triste",
		"seguro",
		"salgado",
		"mesmo",
		"arenoso",
		"são",
		"sarcástico",
		"sardônico",
		"satisfeito",
		"escamoso",
		"escasso",
		"assustado",
		"assustador",
		"perfumado",
		"erudito",
		"científico",
		"desdenhoso",
		"riscado",
		"esquelético",
		"segundo",
		"em-segunda-mão",
		"secundário",
		"secreto",
		"seguro-de-si",
		"auto-suficiente",
		"egoísta",
		"sentimental",
		"separado",
		"sereno",
		"sério",
		"serpentino",
		"variável",
		"grave",
		"sexy",
		"maltrapilho",
		"sombrio",
		"ignoto",
		"rasa",
		"vergonhoso",
		"sem-vergonha",
		"afiado",
		"tremeluzente",
		"brilhante",
		"chocado",
		"chocante",
		"de-má-qualidade",
		"curto",
		"a-curto-prazo",
		"vistoso",
		"estriduloso",
		"tímido",
		"doente",
		"silencioso",
		"sedoso",
		"tonto",
		"prata",
		"semelhante",
		"simples",
		"simplista",
		"pecaminoso",
		"único",
		"abrasador",
		"esquelético",
		"magricela",
		"adormecido",
		"ligeiro",
		"magro",
		"viscoso",
		"escorregadio",
		"lento",
		"lamacento",
		"pequeno",
		"inteligente",
		"poluído",
		"suave",
		"presunçoso",
		"rápido",
		"rosnador",
		"sorrateiro",
		"ranhento",
		"bisbilhoteiro",
		"sociável",
		"suave",
		"ensopado",
		"sólido",
		"sombrio",
		"algum",
		"sofisticado",
		"dorido",
		"doloroso",
		"cheio-de-alma",
		"com-alma",
		"azedo",
		"espanhol",
		"cintilante",
		"esparsa",
		"específico",
		"espectacular",
		"rápido",
		"esférico",
		"picante",
		"elegante",
		"espirituoso",
		"vingativo",
		"esplêndido",
		"imaculado",
		"manchado",
		"ágil",
		"quadrado",
		"guinchante",
		"rabugento",
		"estável",
		"sóbrio",
		"manchado",
		"envelhecido",
		"padrão",
		"amiláceo",
		"duro",
		"estrelado",
		"aço",
		"íngreme",
		"pegajoso",
		"rígido",
		"estimulante",
		"mesquinho",
		"tempestuoso",
		"recto",
		"estranho",
		"rigoroso",
		"estridente",
		"impressionante",
		"às-riscas",
		"forte",
		"estudioso",
		"deslumbrante",
		"estupendo",
		"estúpido",
		"robusto",
		"elegante",
		"subjugado",
		"submisso",
		"substancial",
		"subtil",
		"suburbano",
		"repentino",
		"açucarado",
		"ensolarado",
		"super",
		"soberbo",
		"superficial",
		"superior",
		"de-apoio",
		"surpreendido",
		"suspeito",
		"svelte",
		"suado",
		"doce",
		"sufocante",
		"rápido",
		"simpático",
		"tagarela",
		"alto",
		"domesticado",
		"bronzeado",
		"tangível",
		"torta",
		"saboroso",
		"esfarrapado",
		"esticado",
		"enfadonho",
		"fervilhante",
		"tentador",
		"terno",
		"tenso",
		"tépido",
		"terrível",
		"fantástico",
		"impacientes",
		"grato",
		"safado",
		"estralado",
		"espesso",
		"fino",
		"terceiro",
		"sedento",
		"insuportável",
		"espinhoso",
		"meticuloso",
		"carioca",
		"pensativo",
		"sem-fios",
		"parcimonioso",
		"trovejante",
		"arrumado",
		"apertado",
		"oportuno",
		"colorido",
		"minúsculo",
		"cansado",
		"rasgado",
		"total",
		"duro",
		"trágico",
		"treinados",
		"traumático",
		"tesouro",
		"tremendo",
		"triangular",
		"complicado",
		"insignificante",
		"aparar",
		"trivial",
		"perturbado",
		"verdadeiro",
		"confiante",
		"digno-de-confiança",
		"de-confiança",
		"Verdadeiro",
		"atarracado",
		"turbulento",
		"gêmeo",
		"feio",
		"último",
		"inaceitável",
		"inconsciente",
		"desconfortável",
		"invulgar",
		"inconsciente",
		"subestimado",
		"inigualável",
		"desigual",
		"inacabado",
		"impróprio",
		"desdobrado",
		"desafortunado",
		"infeliz",
		"insalubres",
		"uniforme",
		"sem-importância",
		"único",
		"unido",
		"desgrenhado",
		"desconhecido",
		"ilegal",
		"sem-forro",
		"azarado",
		"não-natural",
		"desagradável",
		"irrealista",
		"verde",
		"indisciplinado",
		"altruísta",
		"de-má-aparência",
		"instável",
		"não-cantado",
		"desarrumado",
		"inoportuno",
		"não-experimentado",
		"enganador",
		"não-utilizado",
		"invulgar",
		"indesejável",
		"pesado",
		"relutante",
		"involuntário",
		"não-escrito",
		"otimista",
		"de-pé",
		"chateado",
		"urbano",
		"utilizável",
		"usado",
		"útil",
		"inútil",
		"utilizado",
		"absoluto",
		"vago",
		"em-vão",
		"vaidoso",
		"válido",
		"valioso",
		"insípido",
		"variável",
		"vasta",
		"aveludado",
		"venerada",
		"vingativo",
		"verificável",
		"vibrante",
		"vicioso",
		"vitorioso",
		"vigilante",
		"vigoroso",
		"vilão",
		"violento",
		"violeta",
		"virtual",
		"virtuoso",
		"visível",
		"vital",
		"vivo",
		"vívido",
		"volumosa",
		"amarelento",
		"bélico",
		"quente",
		"de-coração-quente",
		"deformado",
		"cauteloso",
		"esbanjador",
		"vigilante",
		"alagado",
		"aguado",
		"ondulado",
		"fraco",
		"rico",
		"cansado",
		"alado",
		"wee",
		"semanalmente",
		"chorão",
		"maciço",
		"esquisito",
		"bem-vindo",
		"bem-documentado",
		"bem-apetrechada",
		"bem-informado",
		"bem-iluminado",
		"bem-feito",
		"próspero",
		"bem-gasto",
		"molhado",
		"querido",
		"caprichoso",
		"redemoinho",
		"sussurrado",
		"branco",
		"inteiro",
		"tagarelador",
		"perverso",
		"largo",
		"de-olhos-arregalados",
		"perspicaz",
		"selvagem",
		"disposto",
		"murcha",
		"enrolado",
		"ventoso",
		"alado",
		"rijo",
		"sábio",
		"espirituoso",
		"trêmulo",
		"abominável",
		"maravilhoso",
		"amadeirado",
		"tonto",
		"verboso",
		"mundano",
		"desgastado",
		"preocupada",
		"inquietante",
		"pior",
		"pior-ainda",
		"sem-valor",
		"que-valha-a-pena",
		"digno",
		"furioso",
		"miserável",
		"contorcionista",
		"errado",
		"irônico",
		"bocejador",
		"anualmente",
		"amarelo",
		"amarelado",
		"jovem",
		"jovem",
		"delicioso",
		"maluco",
		"zeloso",
		"picante",
		"ziguezague",
	},
}
