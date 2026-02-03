package currencies

// AllFiatCurrencies contains all ISO 4217 fiat currency codes
var AllFiatCurrencies = []string{
	// Major World Currencies
	"USD", "EUR", "GBP", "JPY", "CHF", "CAD", "AUD", "NZD",

	// African Currencies
	"XOF", // West African CFA Franc (BCEAO)
	"XAF", // Central African CFA Franc (BEAC)
	"NGN", // Nigerian Naira
	"GHS", // Ghanaian Cedi
	"KES", // Kenyan Shilling
	"ZAR", // South African Rand
	"EGP", // Egyptian Pound
	"MAD", // Moroccan Dirham
	"TND", // Tunisian Dinar
	"DZD", // Algerian Dinar
	"LYD", // Libyan Dinar
	"ETB", // Ethiopian Birr
	"TZS", // Tanzanian Shilling
	"UGX", // Ugandan Shilling
	"RWF", // Rwandan Franc
	"ZMW", // Zambian Kwacha
	"BWP", // Botswana Pula
	"NAD", // Namibian Dollar
	"MUR", // Mauritian Rupee
	"SCR", // Seychellois Rupee
	"GMD", // Gambian Dalasi
	"GNF", // Guinean Franc
	"LRD", // Liberian Dollar
	"SLL", // Sierra Leonean Leone
	"CDF", // Congolese Franc
	"AOA", // Angolan Kwanza
	"MZN", // Mozambican Metical
	"MWK", // Malawian Kwacha
	"LSL", // Lesotho Loti
	"SZL", // Swazi Lilangeni
	"SOS", // Somali Shilling
	"DJF", // Djiboutian Franc
	"ERN", // Eritrean Nakfa
	"BIF", // Burundian Franc
	"KMF", // Comorian Franc
	"MGA", // Malagasy Ariary
	"STN", // São Tomé and Príncipe Dobra
	"CVE", // Cape Verdean Escudo
	"MRU", // Mauritanian Ouguiya
	"SDG", // Sudanese Pound

	// Asian Currencies
	"CNY", // Chinese Yuan
	"INR", // Indian Rupee
	"KRW", // South Korean Won
	"SGD", // Singapore Dollar
	"HKD", // Hong Kong Dollar
	"THB", // Thai Baht
	"MYR", // Malaysian Ringgit
	"PHP", // Philippine Peso
	"IDR", // Indonesian Rupiah
	"VND", // Vietnamese Dong
	"PKR", // Pakistani Rupee
	"BDT", // Bangladeshi Taka
	"LKR", // Sri Lankan Rupee
	"NPR", // Nepalese Rupee
	"MMK", // Myanmar Kyat
	"KHR", // Cambodian Riel
	"LAK", // Lao Kip
	"BND", // Brunei Dollar
	"MOP", // Macanese Pataca
	"TWD", // New Taiwan Dollar
	"KZT", // Kazakhstani Tenge
	"UZS", // Uzbekistani Som
	"KGS", // Kyrgyzstani Som
	"TJS", // Tajikistani Somoni
	"TMT", // Turkmenistani Manat
	"AFN", // Afghan Afghani
	"AMD", // Armenian Dram
	"AZN", // Azerbaijani Manat
	"GEL", // Georgian Lari
	"MNT", // Mongolian Tugrik
	"BTN", // Bhutanese Ngultrum
	"MVR", // Maldivian Rufiyaa

	// Middle Eastern Currencies
	"AED", // UAE Dirham
	"SAR", // Saudi Riyal
	"QAR", // Qatari Riyal
	"KWD", // Kuwaiti Dinar
	"BHD", // Bahraini Dinar
	"OMR", // Omani Rial
	"ILS", // Israeli New Shekel
	"TRY", // Turkish Lira
	"JOD", // Jordanian Dinar
	"LBP", // Lebanese Pound
	"SYP", // Syrian Pound
	"IQD", // Iraqi Dinar
	"IRR", // Iranian Rial
	"YER", // Yemeni Rial

	// European Currencies (non-EUR)
	"NOK", // Norwegian Krone
	"SEK", // Swedish Krona
	"DKK", // Danish Krone
	"PLN", // Polish Zloty
	"CZK", // Czech Koruna
	"HUF", // Hungarian Forint
	"RON", // Romanian Leu
	"BGN", // Bulgarian Lev
	"HRK", // Croatian Kuna
	"RSD", // Serbian Dinar
	"BAM", // Bosnia-Herzegovina Convertible Mark
	"MKD", // Macedonian Denar
	"ALL", // Albanian Lek
	"ISK", // Icelandic Króna
	"BYN", // Belarusian Ruble
	"UAH", // Ukrainian Hryvnia
	"MDL", // Moldovan Leu
	"RUB", // Russian Ruble
	"GEL", // Georgian Lari

	// Latin American Currencies
	"BRL", // Brazilian Real
	"MXN", // Mexican Peso
	"ARS", // Argentine Peso
	"CLP", // Chilean Peso
	"COP", // Colombian Peso
	"PEN", // Peruvian Sol
	"UYU", // Uruguayan Peso
	"PYG", // Paraguayan Guarani
	"BOB", // Bolivian Boliviano
	"VEF", // Venezuelan Bolívar
	"CRC", // Costa Rican Colón
	"GTQ", // Guatemalan Quetzal
	"HNL", // Honduran Lempira
	"NIO", // Nicaraguan Córdoba
	"PAB", // Panamanian Balboa
	"DOP", // Dominican Peso
	"HTG", // Haitian Gourde
	"JMD", // Jamaican Dollar
	"TTD", // Trinidad and Tobago Dollar
	"BSD", // Bahamian Dollar
	"BBD", // Barbadian Dollar
	"BZD", // Belize Dollar
	"GYD", // Guyanese Dollar
	"SRD", // Surinamese Dollar
	"CUP", // Cuban Peso
	"AWG", // Aruban Florin
	"ANG", // Netherlands Antillean Guilder

	// Pacific Currencies
	"FJD", // Fijian Dollar
	"PGK", // Papua New Guinean Kina
	"SBD", // Solomon Islands Dollar
	"TOP", // Tongan Paʻanga
	"VUV", // Vanuatu Vatu
	"WST", // Samoan Tala
	"TVD", // Tuvaluan Dollar

	// Other Currencies
	"XCD", // East Caribbean Dollar
	"XDR", // Special Drawing Rights (IMF)
	"XPF", // CFP Franc (French Pacific)
	"FKP", // Falkland Islands Pound
	"GGP", // Guernsey Pound
	"GIP", // Gibraltar Pound
	"IMP", // Isle of Man Pound
	"JEP", // Jersey Pound
	"KYD", // Cayman Islands Dollar
	"BMD", // Bermudian Dollar
	"SHP", // Saint Helena Pound
	"SPL", // Seborgan Luigino
}

// AllCryptoCurrencies contains major cryptocurrency codes
var AllCryptoCurrencies = []string{
	"BTC",   // Bitcoin
	"ETH",   // Ethereum
	"USDT",  // Tether
	"USDC",  // USD Coin
	"BNB",   // Binance Coin
	"XRP",   // Ripple
	"ADA",   // Cardano
	"SOL",   // Solana
	"DOT",   // Polkadot
	"DOGE",  // Dogecoin
	"MATIC", // Polygon
	"LTC",   // Litecoin
	"TRX",   // TRON
	"AVAX",  // Avalanche
	"UNI",   // Uniswap
	"LINK",  // Chainlink
	"ATOM",  // Cosmos
	"XLM",   // Stellar
	"ALGO",  // Algorand
	"VET",   // VeChain
}

// AllCurrencies combines both fiat and crypto currencies
var AllCurrencies = append(AllFiatCurrencies, AllCryptoCurrencies...)

// CommonCurrencies contains the most frequently used currencies
var CommonCurrencies = []string{
	"USD", "EUR", "GBP", "JPY", "CHF", "CAD", "AUD",
	"XOF", "XAF", "NGN", "GHS", "KES", "ZAR",
	"CNY", "INR", "BRL", "MXN",
	"BTC", "ETH", "USDT", "USDC",
}

// AfricanCurrencies contains all African currency codes
var AfricanCurrencies = []string{
	"XOF", "XAF", "NGN", "GHS", "KES", "ZAR", "EGP", "MAD", "TND",
	"DZD", "LYD", "ETB", "TZS", "UGX", "RWF", "ZMW", "BWP", "NAD",
	"MUR", "SCR", "GMD", "GNF", "LRD", "SLL", "CDF", "AOA", "MZN",
	"MWK", "LSL", "SZL", "SOS", "DJF", "ERN", "BIF", "KMF", "MGA",
	"STN", "CVE", "MRU", "SDG",
}

// IsCurrencySupported checks if a currency code is supported
func IsCurrencySupported(code string) bool {
	for _, c := range AllCurrencies {
		if c == code {
			return true
		}
	}
	return false
}

// IsFiatCurrency checks if a currency is fiat
func IsFiatCurrency(code string) bool {
	for _, c := range AllFiatCurrencies {
		if c == code {
			return true
		}
	}
	return false
}

// IsCryptoCurrency checks if a currency is crypto
func IsCryptoCurrency(code string) bool {
	for _, c := range AllCryptoCurrencies {
		if c == code {
			return true
		}
	}
	return false
}
