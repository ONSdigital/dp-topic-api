const censusTopics = [
  {
    title: 'Census',
    description: 'Census',
    subtopics: [
      { title: 'Ageing', description: 'Facts and figures to help understand the needs of an ageing population in England and Wales including changes over time and the characteristics of certain age groups.' },
      { title: 'Demography', description: 'Facts and figures to help understand the population of England and Wales, including age, sex, household composition and legal partnerships.' },
      { title: 'Education', description: 'Facts and figures to help understand the education level of people in England and Wales, and how it varies by age, sex, area, ethnicity, disability and other characteristics.' },
      { title: 'Equalities', description: 'Facts and figures to help understand equality for people in England and Wales, including what life is like for people with certain characteristics.' },
      { title: 'Ethnic group, national identity, language and religion', description: 'Facts and figures to help understand people\'s ethnicity, national identity, language and religion in England and Wales, including changes over time and Welsh language.' },
      { title: 'Health, disability and unpaid care', description: 'Facts and figures to help understand health, disability and unpaid care for people in England and Wales, and how it varies by areas and other characteristics.' },
      { title: 'Historic census', description: 'Facts and figures to help understand how life has changed over time, using data from the 2011 Census and earlier.' },
      { title: 'Housing', description: 'Facts and figures to help understand the types of housing people live in, in England and Wales, including how it has changed over time, by area and household and property characteristics.' },
      { title: 'International migration', description: 'Facts and figures to help understand people who have moved in and out of the UK within England and Wales, including changes over time, people with more than one passport and second-generation migration.' },
      { title: 'Labour market', description: 'Facts and figures to help understand the labour market for people in England and Wales, including how it varies by area and other characteristics.' },
      { title: 'Sexual orientation and gender identity', description: 'Facts and figures to help understand the sexual orientation and gender identity of people in England and Wales, including how it varies by area and characteristics such as demography, housing, employment and education.' },
      { title: 'Travel to work', description: 'Facts and figures to help understand how people in England and Wales travel to work, including mode of transport, distance and how it changes across rural and urban areas.' },
      { title: 'Veterans', description: 'Facts and figures to help understand the veteran population of England and Wales, including topics such as housing, education, employment and skills.' },
    ],
  },
];

function createCensusTopics(){
  createTopLevelTopics(censusTopics);
}
