package api

import (
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
)

const (
	english = "english"
	welsh   = "welsh"
)

var description = map[string]string{
	english: "A list of topical areas and their subtopics in english to generate the website navbar.",
	welsh:   "A list of topical areas and their subtopics in welsh to generate the website navbar in welsh.",
}

var labels = map[string]map[string]string{
	"business-industry-and-trade": {
		english: "Business, industry and trade",
		welsh:   "Busnes, diwydiant a masnach",
	},
	"business": {
		english: "Business",
		welsh:   "Busnes",
	},
	"changes-to-business": {
		english: "Changes to business",
		welsh:   "Newidiadau i fusnesau",
	},
	"construction-industry": {
		english: "Construction industry",
		welsh:   "Diwydiant adeiladu",
	},
	"international-trade": {
		english: "International trade",
		welsh:   "Masnach ryngwladol",
	},
	"it-and-internet-industry": {
		english: "IT and internet industry",
		welsh:   "Y diwydiant TG a'r rhyngrwyd",
	},
	"manufacturing-and-production-industry": {
		english: "Manufacturing and production industry",
		welsh:   "Y diwydiant gweithgynhyrchu a chynhyrchu",
	},
	"retail-industry": {
		english: "Retail industry",
		welsh:   "Y diwydiant manwethu",
	},
	"tourism-industry": {
		english: "Tourism industry",
		welsh:   "Y diwydiant twristiaeth",
	},
	"economy": {
		english: "Economy",
		welsh:   "Yr economi",
	},
	"economic-output-and-productivity": {
		english: "Economic output and productivity",
		welsh:   "Allgynnyrch economaidd a chynhyrchiant",
	},
	"environmental-accounts": {
		english: "Environmental accounts",
		welsh:   "Cyfrifon amgylcheddol",
	},
	"government-public-sector-and-taxes": {
		english: "Government, public sector and taxes",
		welsh:   "Llwodraeth, y sector cyhoeddus a threthi",
	},
	"gross-domestic-product-gdp": {
		english: "Gross Domestic Product (GDP)",
		welsh:   "Cynnyrch Domestig Gros (CDG)",
	},
	"gross-value-added-gva": {
		english: "Gross Value Added (GVA)",
		welsh:   "Gwerth Ychwanegol Gros",
	},
	"inflation-and-price-indices": {
		english: "Inflation and price price indicess",
		welsh:   "Mynegeion chwyddiant a phrisiau",
	},
	"investments-pensions-and-trusts": {
		english: "Investments, pensions and trusts",
		welsh:   "Buddsoddiadau, pensiynau ac ymddiriedolaethau",
	},
	"national-accounts": {
		english: "National accounts",
		welsh:   "Cyfrifon gwladol",
	},
	"regional-accounts": {
		english: "Regional accounts",
		welsh:   "Cyfrifon rhanbarthol",
	},
	"employment-and-labour-market": {
		english: "Employment and labour market",
		welsh:   "Cyflogaeth a'r farchnad lafur",
	},
	"people-in-work": {
		english: "People in work",
		welsh:   "Pobl mewn gwaith",
	},
	"people-not-in-work": {
		english: "People not in work",
		welsh:   "Pobl nad ydynt mewn gwaith",
	},
	"people-population-and-community": {
		english: "People, population and community",
		welsh:   "Pobl, y boblogaeth a chymunedau",
	},
	"births-deaths-and-marriages": {
		english: "Births, deaths and marriages",
		welsh:   "Genedigaethau, marwolaethau a phriodasau",
	},
	"crime-and-justice": {
		english: "Crime and justice",
		welsh:   "Troseddu a chyfiawnder",
	},
	"cultural-identity": {
		english: "Cultural identity",
		welsh:   "Hunaniaeth ddiwylliannol",
	},
	"education-and-childcare": {
		english: "Education and childcare",
		welsh:   "Addysg a gofal plant",
	},
	"elections": {
		english: "Elections",
		welsh:   "Etholiadau",
	},
	"health-and-social-care": {
		english: "Health and social care",
		welsh:   "Iechyd a gofal cymdeithasol",
	},
	"household-characteristics": {
		english: "Household characteristics",
		welsh:   "Nodweddion aelwydydd",
	},
	"housing": {
		english: "Housing",
		welsh:   "Tai",
	},
	"leisure-and-tourism": {
		english: "Leisure and tourism",
		welsh:   "Hamdden a thwristiaeth",
	},
	"personal-and-household-finances": {
		english: "Personal and household finances",
		welsh:   "Cyllid personol a chyllid aelwydydd",
	},
	"population-and-migration": {
		english: "Population and migration",
		welsh:   "Poblogaeth ac ymfudo",
	},
	"wellbeing": {
		english: "Well-being",
		welsh:   "Lles",
	},
	"census": {
		english: "Census",
		welsh:   "Cyfrifiad",
	},
	"taking-part-in-a-survey": {
		english: "Taking part in a survey?",
		welsh:   "Cymryd rhan mewn arolwg?",
	},
}

// getNavigationHandler is currently a hard-coded list of topics to be used for site navigation
func (api *API) getNavigationHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lang := req.URL.Query().Get("lang")
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"function":   "getNavigationPrivateHandler",
	}

	nav := models.Navigation{
		Links: &models.TopicLinks{
			Self: &models.LinkObject{
				HRef: "/navigation",
			},
		},
	}

	if lang == "cy" {
		nav.Items, nav.Description = getNavItems(welsh)
	} else {
		nav.Items, nav.Description = getNavItems(english)
	}

	if err := WriteJSONBody(ctx, nav, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}

}

func getNavItems(lang string) (*[]models.TopicNonReferential, string) {

	return &[]models.TopicNonReferential{
		{
			Title:       "Business, industry and trade",
			Description: "Activities of businesses and industry in the UK, including data on the production and trade of goods and services, sales by retailers, characteristics of businesses, the construction and manufacturing sectors, and international trade.",
			Name:        "business-industry-and-trade",
			Label:       labels["business-industry-and-trade"][lang],
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					ID:   "businessindustryandtrade",
					HRef: "/topics/businessindustryandtrade",
				},
			},
			Uri: "/businessindustryandtrade",
			SubtopicItems: &[]models.TopicNonReferential{
				{
					Title:       "Business",
					Description: "UK businesses registered for VAT and PAYE with regional breakdowns, including data on size (employment and turnover) and activity (type of industry), research and development, and business services.",
					Name:        "business",
					Label:       labels["business"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "business",
							HRef: "/topics/business",
						},
					},
					Uri: "/businessindustryandtrade/business",
				},
				{
					Title:       "Changes to business",
					Description: "UK business growth, survival and change over time. These figures are an informal indicator of confidence in the UK economy.",
					Name:        "changes-to-business",
					Label:       labels["changes-to-business"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "changestobusiness",
							HRef: "/topics/changestobusiness",
						},
					},
					Uri: "/businessindustryandtrade/changestobusiness",
				},
				{
					Title:       "Construction industry",
					Description: "Construction of new buildings and repairs or alterations to existing properties in Great Britain measured by the amount charged for the work, including work by civil engineering companies. ",
					Name:        "construction-industry",
					Label:       labels["construction-industry"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "constructionindustry",
							HRef: "/topics/constructionindustry",
						},
					},
					Uri: "/businessindustryandtrade/constructionindustry",
				},
				{
					Title:       "International trade",
					Description: "Trade in goods and services across the UK's international borders, including total imports and exports, the types of goods and services traded and general trends in international trade. ",
					Name:        "international-trade",
					Label:       labels["international-trade"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "internationaltrade",
							HRef: "/topics/internationaltrade",
						},
					},
					Uri: "/businessindustryandtrade/internationaltrade",
				},
				{
					Title:       "IT and internet industry",
					Description: "Internet sales by businesses in the UK (total value and as a percentage of all retail sales) and the percentage of businesses that have a website and broadband connection. These figures indicate the importance of the internet to UK businesses.",
					Name:        "it-and-internet-industry",
					Label:       labels["it-and-internet-industry"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "itandinternetindustry",
							HRef: "/topics/itandinternetindustry",
						},
					},
					Uri: "/businessindustryandtrade/itandinternetindustry",
				},
				{
					Title:       "Manufacturing and production industry",
					Description: "UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and UK manufactures' sales by product and industrial division, with EU comparisons.",
					Name:        "manufacturing-and-production-industry",
					Label:       labels["manufacturing-and-production-industry"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "manufacturingandproductionindustry",
							HRef: "/topics/manufacturingandproductionindustry",
						},
					},
					Uri: "/businessindustryandtrade/manufacturingandproductionindustry",
				},
				{
					Title:       "Retail industry",
					Description: "",
					Name:        "retail-industry",
					Label:       labels["retail-industry"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "retailindustry",
							HRef: "/topics/retailindustry",
						},
					},
					Uri: "/businessindustryandtrade/retailindustry",
				},
				{
					Title:       "Tourism industry",
					Description: "Tourism and travel (including accommodation services, food and beverage services, passenger transport services, vehicle hire, travel agencies and sports, recreational and conference services), employment levels and output of the tourism industry, the number of visitors to the UK and the amount they spend.",
					Name:        "tourism-industry",
					Label:       labels["tourism-industry"][lang],
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "tourismindustry",
							HRef: "/topics/tourismindustry",
						},
					},
					Uri: "/businessindustryandtrade/tourismindustry",
				},
			},
		},
		{
			Title:       "Economy",
			Description: "UK economic activity covering production, distribution, consumption and trade of goods and services. Individuals, businesses, organisations and governments all affect the development of the economy.",
			Name:        "economy",
			Label:       labels["economy"][lang],
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					ID:   "economy",
					HRef: "/topics/economy",
				},
			},
			Uri: "/economy",
			SubtopicItems: &[]models.TopicNonReferential{
				{
					Description: "Manufacturing, production and services indices (measuring total economic output) and productivity (measuring efficiency, expressed as a ratio of output to input over a given period of time, for example output per person per hour).",
					Title:       "Economic output and productivity",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "economicoutputandproductivity",
							HRef: "/topics/economicoutputandproductivity",
						},
					},
					Uri:   "/economy/economicoutputandproductivity",
					Name:  "economic-output-and-productivity",
					Label: labels["economic-output-and-productivity"][lang],
				},
				{
					Description: "Environmental accounts show how the environment contributes to the economy (for example, through the extraction of raw materials), the impacts that the economy has on the environment (for example, energy consumption and air emissions), and how society responds to environmental issues (for example, through taxation and expenditure on environmental protection).",
					Title:       "Environmental accounts",
					Uri:         "/economy/environmentalaccounts",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "environmentalaccounts",
							HRef: "/topics/environmentalaccounts",
						},
					},
					Name:  "environmental-accounts",
					Label: labels["environmental-accounts"][lang],
				},
				{
					Description: "Public sector spending, tax revenues and investments for the UK, including government debt and deficit (the gap between revenue and spending), research and development, and the effect of taxes.",
					Title:       "Government, public sector and taxes",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "governmentpublicsectorandtaxes",
							HRef: "/topics/governmentpublicsectorandtaxes",
						},
					},
					Uri:   "/economy/governmentpublicsectorandtaxes",
					Name:  "government-public-sector-and-taxes",
					Label: labels["government-public-sector-and-taxes"][lang],
				},
				{
					Description: "Estimates of GDP are released on a monthly and quarterly basis. Monthly estimates are released alongside other short-term economic indicators. The two quarterly estimates contain data from all three approaches to measuring GDP and are called the First quarterly estimate of GDP and the Quarterly National Accounts. Data sources feeding into the two types of releases are consistent with each other.",
					Title:       "Gross Domestic Product (GDP)",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "grossdomesticproductgdp",
							HRef: "/topics/grossdomesticproductgdp",
						},
					},
					Uri:   "/economy/grossdomesticproductgdp",
					Name:  "gross-domestic-product-gdp",
					Label: labels["gross-domestic-product-gdp"][lang],
				},
				{
					Description: "Regional gross value added using production (GVA(P)) and income (GVA(I)) approaches. Regional gross value added is the value generated by any unit engaged in the production of goods and services. GVA per head is a useful way of comparing regions of different sizes. It is not, however, a measure of regional productivity.",
					Title:       "Gross Value Added (GVA)",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "grossvalueaddedgva",
							HRef: "/topics/grossvalueaddedgva",
						},
					},
					Uri:   "/economy/grossvalueaddedgva",
					Name:  "gross-value-added-gva",
					Label: labels["gross-value-added-gva"][lang],
				},
				{
					Description: "The rate of increase in prices for goods and services. Measures of inflation and prices include consumer price inflation, producer price inflation, the house price index, index of private housing rental prices, and construction output price indices. ",
					Title:       "Inflation and price indices",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "inflationandpriceindices",
							HRef: "/topics/inflationandpriceindices",
						},
					},
					Uri:   "/economy/inflationandpriceindices",
					Name:  "inflation-and-price-indices",
					Label: labels["inflation-and-price-indices"][lang],
				},
				{
					Description: "Net flows of investment into the UK, the number of people who hold pensions of different types, and investments made by various types of trusts. ",
					Title:       "Investments, pensions and trusts",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "investmentspensionsandtrusts",
							HRef: "/topics/investmentspensionsandtrusts",
						},
					},
					Uri:   "/economy/investmentspensionsandtrusts",
					Name:  "investments-pensions-and-trusts",
					Label: labels["investments-pensions-and-trusts"][lang],
				},
				{
					Description: "Core accounts for the UK economy as a whole; individual sectors (sector accounts); accounts for the regions, subregions and local areas of the UK; and satellite accounts that cover activities linked to the economy. The national accounts framework brings units and transactions together to provide a simple and understandable description of production, income, consumption, accumulation and wealth.",
					Title:       "National accounts",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "nationalaccounts",
							HRef: "/topics/nationalaccounts",
						},
					},
					Uri:   "/economy/nationalaccounts",
					Name:  "national-accounts",
					Label: labels["national-accounts"][lang],
				},
				{
					Description: "Accounts for regions, sub-regions and local areas of the UK. These accounts allow comparisons between regions and against a UK average. Statistics include regional gross value added (GVA) and figures on regional gross disposable household income (GDHI).",
					Title:       "Regional accounts",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "regionalaccounts",
							HRef: "/topics/regionalaccounts",
						},
					},
					Uri:   "/economy/regionalaccounts",
					Name:  "regional-accounts",
					Label: labels["regional-accounts"][lang],
				},
			},
		},
		{
			Title:       "Employment and labour market",
			Description: "People in and out of work covering employment, unemployment, types of work, earnings, working patterns and workplace disputes.",
			Name:        "employment-and-labour-market",
			Label:       labels["employment-and-labour-market"][lang],
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					ID:   "employmentandlabourmarket",
					HRef: "/topics/employmentandlabourmarket",
				},
			},
			Uri: "/employmentandlabourmarket",
			SubtopicItems: &[]models.TopicNonReferential{
				{
					Description: "Employment data covering employment rates, hours of work, claimants and earnings.",
					Title:       "People in work",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "peopleinwork",
							HRef: "/topics/peopleinwork",
						},
					},
					Uri:   "/employmentandlabourmarket/peopleinwork",
					Name:  "people-in-work",
					Label: labels["people-in-work"][lang],
				},
				{
					Description: "Unemployed and economically inactive people in the UK including claimants of out-of-work benefits and the number of redundancies.\n",
					Title:       "People not in work",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "peoplenotinwork",
							HRef: "/topics/peoplenotinwork",
						},
					},
					Uri:   "/employmentandlabourmarket/peoplenotinwork",
					Name:  "people-not-in-work",
					Label: labels["people-bot-in-work"][lang],
				},
			},
		},
		{
			Title:       "People, population and community",
			Description: "People living in the UK, changes in the population, how we spend our money, and data on crime, relationships, health and religion. These statistics help us build a detailed picture of how we live.",
			Name:        "people-population-and-community",
			Label:       labels["people-population-and-community"][lang],
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					ID:   "peoplepopulationandcommunity",
					HRef: "/topics/peoplepopulationandcommunity",
				},
			},
			Uri: "/peoplepopulationandcommunity",
			SubtopicItems: &[]models.TopicNonReferential{
				{
					Description: "Life events in the UK including fertility rates, live births and stillbirths, family composition, life expectancy and deaths. This tells us about the health and relationships of the population.",
					Title:       "Births, deaths and marriages",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "birthsdeathsandmarriages",
							HRef: "/topics/birthsdeathsandmarriages",
						},
					},
					Uri:   "/peoplepopulationandcommunity/birthsdeathsandmarriages",
					Name:  "births-deaths-and-marriages",
					Label: labels["births-deaths-and-marriages"][lang],
				},
				{
					Description: "Crimes committed and the victims' characteristics, sourced from crimes recorded by the police and from the Crime Survey for England and Wales (CSEW). The outcomes of crime in terms of what happened to the offender are also included.",
					Title:       "Crime and justice",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "crimeandjustice",
							HRef: "/topics/crimeandjustice",
						},
					},
					Uri:   "/peoplepopulationandcommunity/crimeandjustice",
					Name:  "crime-and-justice",
					Label: labels["crime-and-justice"][lang],
				},
				{
					Description: "How people in the UK see themselves today in terms of ethnicity, sexual identity, religion and language, and how this has changed over time. We use a diverse range of sources for this data.",
					Title:       "Cultural identity",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "culturalidentity",
							HRef: "/topics/culturalidentity",
						},
					},
					Uri:   "/peoplepopulationandcommunity/culturalidentity",
					Name:  "cultural-identity",
					Label: labels["cultural-identity"][lang],
				},
				{
					Description: "Early years childcare, school and college education, and higher education and adult learning, including qualifications, personnel, and safety and well-being. ",
					Title:       "Education and childcare",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "educationandchildcare",
							HRef: "/topics/educationandchildcare",
						},
					},
					Uri:   "/peoplepopulationandcommunity/educationandchildcare",
					Name:  "education-and-childcare",
					Label: labels["education-and-childcare"][lang],
				},
				{
					Description: "",
					Title:       "Elections",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "elections",
							HRef: "/topics/elections",
						},
					},
					Uri:   "/peoplepopulationandcommunity/elections",
					Name:  "elections",
					Label: labels["elections"][lang],
				},
				{
					Description: "Life expectancy and the impact of factors such as occupation, illness and drug misuse. We collect these statistics from registrations and surveys. ",
					Title:       "Health and social care",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "healthandsocialcare",
							HRef: "/topics/healthandsocialcare",
						},
					},
					Uri:   "/peoplepopulationandcommunity/healthandsocialcare",
					Name:  "health-and-social-care",
					Label: labels["health-and-social-care"][lang],
				},
				{
					Description: "The composition of households, including those who live alone, overcrowding and under-occupation, as well as internet and social media usage by household.",
					Title:       "Household characteristics",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "householdcharacteristics",
							HRef: "/topics/householdcharacteristics",
						},
					},
					Uri:   "/peoplepopulationandcommunity/householdcharacteristics",
					Name:  "household-characteristics",
					Label: labels["household-characteristics"][lang],
				},
				{
					Description: "Property price, private rent and household survey and census statistics, used by government and other organisations for the creation and fulfilment of housing policy in the UK.",
					Title:       "Housing",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "housing",
							HRef: "/topics/housing",
						},
					},
					Uri:   "/peoplepopulationandcommunity/housing",
					Name:  "housing",
					Label: labels["housing"][lang],
				},
				{
					Description: "Visits and visitors to the UK, the reasons for visiting and the amount of money they spent here. Also UK residents travelling abroad, their reasons for travel and the amount of money they spent. The statistics on UK residents travelling abroad are an informal indicator of living standards.",
					Title:       "Leisure and tourism",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "leisureandtourism",
							HRef: "/topics/leisureandtourism",
						},
					},
					Uri:   "/peoplepopulationandcommunity/leisureandtourism",
					Name:  "leisure-and-tourism",
					Label: labels["leisure-and-tourism"][lang],
				},
				{
					Description: "Earnings and household spending, including household and personal debt, expenditure, and income and wealth. These statistics help build a picture of our spending and saving decisions. ",
					Title:       "Personal and household finances",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "personalandhouseholdfinances",
							HRef: "/topics/personalandhouseholdfinances",
						},
					},
					Uri:   "/peoplepopulationandcommunity/personalandhouseholdfinances",
					Name:  "personal-and-household-finances",
					Label: labels["personal-and-household-finances"][lang],
				},
				{
					Description: "Size, age, sex and geographic distribution of the UK population, and changes in the UK population and the factors driving these changes. These statistics have a wide range of uses. Central government, local government and the health sector use them for planning, resource allocation and managing the economy. They are also used by people such as market researchers and academics.",
					Title:       "Population and migration",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "populationandmigration",
							HRef: "/topics/populationandmigration",
						},
					},
					Uri:   "/peoplepopulationandcommunity/populationandmigration",
					Name:  "population-and-migration",
					Label: labels["population-and-migration"][lang],
				},
				{
					Description: "Societal and personal well-being in the UK looking beyond what we produce, to areas such as health, relationships, education and skills, what we do, where we live, our finances and the environment. This data comes from a variety of sources and much of the analysis is new.",
					Title:       "Well-being",
					Links: &models.TopicLinks{
						Self: &models.LinkObject{
							ID:   "wellbeing",
							HRef: "/topics/wellbeing",
						},
					},
					Uri:   "/peoplepopulationandcommunity/wellbeing",
					Name:  "wellbeing",
					Label: labels["wellbeing"][lang],
				},
			},
		},
		{
			Title:       "Census",
			Description: "",
			Name:        "census",
			Label:       labels["census"][lang],
			Links: &models.TopicLinks{
				Self: &models.LinkObject{
					ID:   "census",
					HRef: "/topics/census",
				},
			},
			Uri: "/census",
		},
		{
			Title:       "Survey",
			Description: "",
			Name:        "taking-part-in-a-survey",
			Label:       labels["taking-part-in-a-survey"][lang],
			Uri:         "/surveys",
		},
	}, description[lang]
}
