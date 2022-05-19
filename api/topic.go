package api

import (
	"context"
	"net/http"

	dprequest "github.com/ONSdigital/dp-net/request"
	"github.com/ONSdigital/dp-renderer/helper"
	"github.com/ONSdigital/dp-topic-api/apierrors"
	"github.com/ONSdigital/dp-topic-api/assets"
	"github.com/ONSdigital/dp-topic-api/models"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/mux"
)

// getTopicPublicHandler is a handler that gets a topic by its id from MongoDB for Web
func (api *API) getTopicPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicPublicHandler",
	}

	if id == "topic_root" {
		handleError(ctx, w, apierrors.ErrTopicNotFound, logdata)
		return
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document
	if err := WriteJSONBody(ctx, topic.Current, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

// getTopicPrivateHandler is a handler that gets a topic by its id from MongoDB for Publishing
func (api *API) getTopicPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicPrivateHandler",
	}

	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw topic document
	if err := WriteJSONBody(ctx, topic, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

func (api *API) getSubtopicsPublicByID(ctx context.Context, id string, logdata log.Data, w http.ResponseWriter) {
	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		// no topic found to retrieve the subtopics from
		handleError(ctx, w, err, logdata)
		return
	}

	// User is not authenticated and hence has only access to current sub document(s)
	var result models.PublicSubtopics

	if topic.Current == nil {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	if len(topic.Current.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range topic.Current.SubtopicIds {
		// get sub topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(ctx, subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Error(ctx, "missing subtopic for id", err, logdata)
			continue
		}

		if result.PublicItems == nil {
			result.PublicItems = &[]models.Topic{*topic.Current}
		} else {
			*result.PublicItems = append(*result.PublicItems, *topic.Current)
		}

		result.TotalCount++
	}
	if result.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrContentNotFound, logdata)
		return
	}

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

// getSubtopicsPublicHandler is a handler that gets a topic by its id from MongoDB for Web
func (api *API) getSubtopicsPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getSubtopicsPublicHandler",
	}

	if id == "topic_root" {
		handleError(ctx, w, apierrors.ErrTopicNotFound, logdata)
		return
	}

	api.getSubtopicsPublicByID(ctx, id, logdata, w)
}

func (api *API) getSubtopicsPrivateByID(ctx context.Context, id string, logdata log.Data, w http.ResponseWriter) {
	// get topic from mongoDB by id
	topic, err := api.dataStore.Backend.GetTopic(ctx, id)
	if err != nil {
		// no topic found to retrieve the subtopics from
		handleError(ctx, w, err, logdata)
		return
	}

	// User has valid authentication to get raw full topic document(s)
	var result models.PrivateSubtopics

	if topic.Next == nil {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if len(topic.Next.SubtopicIds) == 0 {
		// no subtopics exist for the requested ID
		handleError(ctx, w, apierrors.ErrNotFound, logdata)
		return
	}

	for _, subTopicID := range topic.Next.SubtopicIds {
		// get topic from mongoDB by subTopicID
		topic, err := api.dataStore.Backend.GetTopic(ctx, subTopicID)
		if err != nil {
			logdata["missing subtopic for id"] = subTopicID
			log.Error(ctx, "missing subtopic for id", err, logdata)
			continue
		}

		if result.PrivateItems == nil {
			result.PrivateItems = &[]models.TopicResponse{*topic}
		} else {
			*result.PrivateItems = append(*result.PrivateItems, *topic)
		}

		result.TotalCount++
	}
	if result.TotalCount == 0 {
		handleError(ctx, w, apierrors.ErrInternalServer, logdata)
		return
	}

	if err := WriteJSONBody(ctx, result, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
	log.Info(ctx, "request successful", logdata) // NOTE: name of function is in logdata
}

// getSubtopicsPrivateHandler is a handler that gets a topic by its id from MongoDB for Publishing
func (api *API) getSubtopicsPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	vars := mux.Vars(req)
	id := vars["id"]
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getSubtopicsPrivateHandler",
	}

	api.getSubtopicsPrivateByID(ctx, id, logdata, w)
}

// getTopicsListPublicHandler is a handler that gets a public list of top level topics by a specific id from MongoDB for Web
func (api *API) getTopicsListPublicHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := "topic_root" // access specific document to retrieve list
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicsListPublicHandler",
	}

	// The mongo document with id: `topic_root` contains the list of sobtopics,
	// so we directly return that list
	api.getSubtopicsPublicByID(ctx, id, logdata, w)
}

// getTopicsListPrivateHandler is a handler that gets a private list of top level topics by a specific id from MongoDB for Web
func (api *API) getTopicsListPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	id := "topic_root" // access specific document to retrieve list
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"topic_id":   id,
		"function":   "getTopicsListPrivateHandler",
	}

	// The mongo document with id: `topic_root` contains the list of sobtopics,
	// so we directly return that list
	api.getSubtopicsPrivateByID(ctx, id, logdata, w)
}

// getNavigationPrivateHandler is currently a hard-coded list of topics to be used for site navigation
func (api *API) getNavigationPrivateHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	lang := req.URL.Query().Get("lang")
	logdata := log.Data{
		"request_id": ctx.Value(dprequest.RequestIdKey),
		"function":   "getNavigationPrivateHandler",
	}

	if lang != "en" && lang != "cy" {
		handleError(ctx, w, apierrors.ErrContentUnrecognisedParameter, logdata)
		return
	}
	helper.InitialiseLocalisationsHelper(assets.Asset)

	items := []models.TopicNonReferential{
		{
			Title:       "Business, industry and trade",
			Description: "Activities of businesses and industry in the UK, including data on the production and trade of goods and services, sales by retailers, characteristics of businesses, the construction and manufacturing sectors, and international trade.",
			Name:        "business-industry-and-trade",
			Label:       helper.Localise("BusinessIndustryTrade", lang, 1),
			Href:        "/topics/businessindustryandtrade",
			SubtopicItems: []models.TopicNonReferential{
				{
					Title:       "Business",
					Description: "UK businesses registered for VAT and PAYE with regional breakdowns, including data on size (employment and turnover) and activity (type of industry), research and development, and business services.",
					Name:        "business",
					Label:       helper.Localise("Business", lang, 1),
					Href:        "/topics/business",
				},
				{
					Title:       "Changes to business",
					Description: "UK business growth, survival and change over time. These figures are an informal indicator of confidence in the UK economy.",
					Name:        "changes-to-business",
					Label:       helper.Localise("ChangesToBusiness", lang, 1),
					Href:        "/topics/changestobusiness",
				},
				{
					Title:       "Construction industry",
					Description: "Construction of new buildings and repairs or alterations to existing properties in Great Britain measured by the amount charged for the work, including work by civil engineering companies. ",
					Name:        "construction-industry",
					Label:       helper.Localise("ConstructionIndustry", lang, 1),
					Href:        "http://localhost:25300/topics/constructionindustry",
				},
				{
					Title:       "International trade",
					Description: "Trade in goods and services across the UK's international borders, including total imports and exports, the types of goods and services traded and general trends in international trade. ",
					Name:        "international-trade",
					Label:       helper.Localise("InternationalTrade", lang, 1),
					Href:        "http://localhost:25300/topics/internationaltrade",
				},
				{
					Title:       "IT and internet industry",
					Description: "Internet sales by businesses in the UK (total value and as a percentage of all retail sales) and the percentage of businesses that have a website and broadband connection. These figures indicate the importance of the internet to UK businesses.",
					Name:        "it-and-internet-industry",
					Label:       helper.Localise("ITInternetIndustry", lang, 1),
					Href:        "http://localhost:25300/topics/itandinternetindustry",
				},
				{
					Title:       "Manufacturing and production industry",
					Description: "UK manufacturing and other production industries (such as mining and quarrying, energy supply, water supply and waste management), including total UK production output, and UK manufactures' sales by product and industrial division, with EU comparisons.",
					Name:        "manufacturing-and-production-industry",
					Label:       helper.Localise("ManufacturingProductionIndustry", lang, 1),
					Href:        "http://localhost:25300/topics/manufacturingandproductionindustry",
				},
				{
					Title:       "Retail industry",
					Description: "",
					Name:        "retail-industry",
					Label:       helper.Localise("RetailIndustry", lang, 1),
					Href:        "http://localhost:25300/topics/retailindustry",
				},
				{
					Title:       "Tourism industry",
					Description: "Tourism and travel (including accommodation services, food and beverage services, passenger transport services, vehicle hire, travel agencies and sports, recreational and conference services), employment levels and output of the tourism industry, the number of visitors to the UK and the amount they spend.",
					Name:        "tourism-industry",
					Label:       helper.Localise("TourismIndustry", lang, 1),
					Href:        "http://localhost:25300/topics/tourismindustry",
				},
			},
		},
		{
			Title:       "Economy",
			Description: "UK economic activity covering production, distribution, consumption and trade of goods and services. Individuals, businesses, organisations and governments all affect the development of the economy.",
			Name:        "economy",
			Label:       helper.Localise("Economy", lang, 1),
			Href:        "/topics/economy",
			SubtopicItems: []models.TopicNonReferential{
				{
					Description: "Manufacturing, production and services indices (measuring total economic output) and productivity (measuring efficiency, expressed as a ratio of output to input over a given period of time, for example output per person per hour).",
					Title:       "Economic output and productivity",
					Href:        "/topics/economicoutputandproductivity",
					Name:        "economic-output-and-productivity",
					Label:       helper.Localise("EconomicOutputProductivity", lang, 1),
				},
				{
					Description: "Environmental accounts show how the environment contributes to the economy (for example, through the extraction of raw materials), the impacts that the economy has on the environment (for example, energy consumption and air emissions), and how society responds to environmental issues (for example, through taxation and expenditure on environmental protection).",
					Title:       "Environmental accounts",
					Href:        "/topics/environmentalaccounts",
					Name:        "environmental-accounts",
					Label:       helper.Localise("EnvironmentalAccounts", lang, 1),
				},
				{
					Description: "Public sector spending, tax revenues and investments for the UK, including government debt and deficit (the gap between revenue and spending), research and development, and the effect of taxes.",
					Title:       "Government, public sector and taxes",
					Href:        "/topics/governmentpublicsectorandtaxes",
					Name:        "government-public-sector-and-taxes",
					Label:       helper.Localise("GovernmentPublicSectorTaxes", lang, 1),
				},
				{
					Description: "Estimates of GDP are released on a monthly and quarterly basis. Monthly estimates are released alongside other short-term economic indicators. The two quarterly estimates contain data from all three approaches to measuring GDP and are called the First quarterly estimate of GDP and the Quarterly National Accounts. Data sources feeding into the two types of releases are consistent with each other.",
					Title:       "Gross Domestic Product (GDP)",
					Href:        "/topics/grossdomesticproductgdp",
					Name:        "gross-domestic-product-gdp",
					Label:       helper.Localise("GrossDomesticProduct", lang, 1),
				},
				{
					Description: "Regional gross value added using production (GVA(P)) and income (GVA(I)) approaches. Regional gross value added is the value generated by any unit engaged in the production of goods and services. GVA per head is a useful way of comparing regions of different sizes. It is not, however, a measure of regional productivity.",
					Title:       "Gross Value Added (GVA)",
					Href:        "/topics/grossvalueaddedgva",
					Name:        "gross-value-added-gva",
					Label:       helper.Localise("GrossDomesticValue", lang, 1),
				},
				{
					Description: "The rate of increase in prices for goods and services. Measures of inflation and prices include consumer price inflation, producer price inflation, the house price index, index of private housing rental prices, and construction output price indices. ",
					Title:       "Inflation and price indices",
					Href:        "/topics/inflationandpriceindices",
					Name:        "inflation-and-price-indices",
					Label:       helper.Localise("InflationPriceIndices", lang, 1),
				},
				{
					Description: "Net flows of investment into the UK, the number of people who hold pensions of different types, and investments made by various types of trusts. ",
					Title:       "Investments, pensions and trusts",
					Href:        "/topics/investmentspensionsandtrusts",
					Name:        "investments-pensions-and-trusts",
					Label:       helper.Localise("InvestmentsPensionsTrusts", lang, 1),
				},
				{
					Description: "Core accounts for the UK economy as a whole; individual sectors (sector accounts); accounts for the regions, subregions and local areas of the UK; and satellite accounts that cover activities linked to the economy. The national accounts framework brings units and transactions together to provide a simple and understandable description of production, income, consumption, accumulation and wealth.",
					Title:       "National accounts",
					Href:        "/topics/nationalaccounts",
					Name:        "national-accounts",
					Label:       helper.Localise("NationalAccounts", lang, 1),
				},
				{
					Description: "Accounts for regions, sub-regions and local areas of the UK. These accounts allow comparisons between regions and against a UK average. Statistics include regional gross value added (GVA) and figures on regional gross disposable household income (GDHI).",
					Title:       "Regional accounts",
					Href:        "/topics/regionalaccounts",
					Name:        "regional-accounts",
					Label:       helper.Localise("RegionalAccounts", lang, 1),
				},
			},
		},
		{
			Title:       "Employment and labour market",
			Description: "People in and out of work covering employment, unemployment, types of work, earnings, working patterns and workplace disputes.",
			Name:        "employment-and-labour-market",
			Label:       helper.Localise("EmploymentLabourMarket", lang, 1),
			Href:        "/topics/employmentandlabourmarket",
			SubtopicItems: []models.TopicNonReferential{
				{
					Description: "Employment data covering employment rates, hours of work, claimants and earnings.",
					Title:       "People in work",
					Href:        "http://localhost:25300/topics/peopleinwork",
					Name:        "people-in-work",
					Label:       helper.Localise("PeopleInWork", lang, 1),
				},
				{
					Description: "Unemployed and economically inactive people in the UK including claimants of out-of-work benefits and the number of redundancies.\n",
					Title:       "People not in work",
					Href:        "http://localhost:25300/topics/peoplenotinwork",
					Name:        "people-not-in-work",
					Label:       helper.Localise("PeopleNotInWork", lang, 1),
				},
			},
		},
		{
			Title:       "People, population and community",
			Description: "People living in the UK, changes in the population, how we spend our money, and data on crime, relationships, health and religion. These statistics help us build a detailed picture of how we live.",
			Name:        "people-population-and-community",
			Label:       helper.Localise("PeoplePopulationCommunity", lang, 1),
			Href:        "/topics/peoplepopulationandcommunity",
			SubtopicItems: []models.TopicNonReferential{
				{
					Description: "Life events in the UK including fertility rates, live births and stillbirths, family composition, life expectancy and deaths. This tells us about the health and relationships of the population.",
					Title:       "Births, deaths and marriages",
					Href:        "http://localhost:25300/topics/birthsdeathsandmarriages",
					Name:        "births-deaths-and-marriages",
					Label:       helper.Localise("BirthsDeathsMarriages", lang, 1),
				},
				{
					Description: "Crimes committed and the victims' characteristics, sourced from crimes recorded by the police and from the Crime Survey for England and Wales (CSEW). The outcomes of crime in terms of what happened to the offender are also included.",
					Title:       "Crime and justice",
					Href:        "http://localhost:25300/topics/crimeandjustice",
					Name:        "crime-and-justice",
					Label:       helper.Localise("CrimeJustice", lang, 1),
				},
				{
					Description: "How people in the UK see themselves today in terms of ethnicity, sexual identity, religion and language, and how this has changed over time. We use a diverse range of sources for this data.",
					Title:       "Cultural identity",
					Href:        "http://localhost:25300/topics/culturalidentity",
					Name:        "cultural-identity",
					Label:       helper.Localise("CulturalIdentity", lang, 1),
				},
				{
					Description: "Early years childcare, school and college education, and higher education and adult learning, including qualifications, personnel, and safety and well-being. ",
					Title:       "Education and childcare",
					Href:        "http://localhost:25300/topics/educationandchildcare",
					Name:        "education-and-childcare",
					Label:       helper.Localise("EducationChildcare", lang, 1),
				},
				{
					Description: "",
					Title:       "Elections",
					Href:        "http://localhost:25300/topics/elections",
					Name:        "elections",
					Label:       helper.Localise("Elections", lang, 1),
				},
				{
					Description: "Life expectancy and the impact of factors such as occupation, illness and drug misuse. We collect these statistics from registrations and surveys. ",
					Title:       "Health and social care",
					Href:        "http://localhost:25300/topics/healthandsocialcare",
					Name:        "health-and-social-care",
					Label:       helper.Localise("HealthSocialCare", lang, 1),
				},
				{
					Description: "The composition of households, including those who live alone, overcrowding and under-occupation, as well as internet and social media usage by household.",
					Title:       "Household characteristics",
					Href:        "http://localhost:25300/topics/householdcharacteristics",
					Name:        "household-characteristics",
					Label:       helper.Localise("HouseholdCharacteristics", lang, 1),
				},
				{
					Description: "Property price, private rent and household survey and census statistics, used by government and other organisations for the creation and fulfilment of housing policy in the UK.",
					Title:       "Housing",
					Href:        "http://localhost:25300/topics/housing",
					Name:        "housing",
					Label:       helper.Localise("Housing", lang, 1),
				},
				{
					Description: "Visits and visitors to the UK, the reasons for visiting and the amount of money they spent here. Also UK residents travelling abroad, their reasons for travel and the amount of money they spent. The statistics on UK residents travelling abroad are an informal indicator of living standards.",
					Title:       "Leisure and tourism",
					Href:        "http://localhost:25300/topics/leisureandtourism",
					Name:        "leisure-and-tourism",
					Label:       helper.Localise("LeisureTourism", lang, 1),
				},
				{
					Description: "Earnings and household spending, including household and personal debt, expenditure, and income and wealth. These statistics help build a picture of our spending and saving decisions. ",
					Title:       "Personal and household finances",
					Href:        "http://localhost:25300/topics/personalandhouseholdfinances",
					Name:        "personal-and-household-finances",
					Label:       helper.Localise("PersonalHouseholdFinances", lang, 1),
				},
				{
					Description: "Size, age, sex and geographic distribution of the UK population, and changes in the UK population and the factors driving these changes. These statistics have a wide range of uses. Central government, local government and the health sector use them for planning, resource allocation and managing the economy. They are also used by people such as market researchers and academics.",
					Title:       "Population and migration",
					Href:        "http://localhost:25300/topics/populationandmigration",
					Name:        "population-and-migration",
					Label:       helper.Localise("PopulationMigration", lang, 1),
				},
				{
					Description: "Societal and personal well-being in the UK looking beyond what we produce, to areas such as health, relationships, education and skills, what we do, where we live, our finances and the environment. This data comes from a variety of sources and much of the analysis is new.",
					Title:       "Well-being",
					Href:        "http://localhost:25300/topics/wellbeing",
					Name:        "wellbeing",
					Label:       helper.Localise("Wellbeing", lang, 1),
				},
			},
		},
	}

	if err := WriteJSONBody(ctx, items, w, logdata); err != nil {
		// WriteJSONBody has already logged the error
		return
	}
}
