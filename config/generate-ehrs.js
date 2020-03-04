#!/usr/bin/env node

const faker = require('faker/locale/nl')
const uuid = require('uuid/v4')

let patients = new Array(10)
let observations = new Array(20)

patients = patients.fill().map(() => ({
  id: uuid(),
  bsn: faker.fake('{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}{{random.number(9)}}'),
  name: {
    given: faker.name.firstName(),
    family: faker.name.lastName()
  },
  gender: faker.random.arrayElement(['male', 'female']),
  birthDate: faker.date.past().toLocaleDateString('nl-NL')
}))

observations = observations.fill().map(() => ({
  id: uuid(),
  patientId: faker.random.arrayElement(patients).id,
  timestamp: faker.date.past().toLocaleString('nl-NL'),
  content: faker.lorem.paragraphs()
}))

console.log(JSON.stringify({
  patients, observations
}, null, 2))
