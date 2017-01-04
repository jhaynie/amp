const SwaggerParser = require('swagger-parser')
const fs = require('fs-extra')

function findSpecs () {
  return new Promise((resolve, reject) => {
    const specs = []
    fs.walk('../../api/rpc/')
      .on('error', reject)
      .on('data', item => {
        if (/build/.test(item.path)) {
          return false
        }
        if (/\.swagger\.json$/.test(item.path)) {
          specs.push(item.path)
        }
      })
      .on('end', () => resolve(specs))
  })
}

async function main () {
  const specs = await findSpecs()
  const fullapi = {
    swagger: '2.0',
    info: {
      title: 'github.com/appcelerator/amp',
      version: 'v1'
    },
    schemes: [
      'http',
      'https'
    ],
    consumes: [
      'application/json'
    ],
    produces: [
      'application/json'
    ],
    paths: {},
    definitions: {}
  }
  for (const spec of specs) {
    const api = await await SwaggerParser.validate(spec)
    fullapi.paths = Object.assign(fullapi.paths, api.paths)
    fullapi.definitions = Object.assign(fullapi.definitions, api.definitions)
  }
  await SwaggerParser.validate(fullapi)
  const bundle = await SwaggerParser.bundle(fullapi)
  fs.writeJsonSync('../../cmd/amplifier-gateway/ui/v1/swagger.json', bundle)
}

main().catch(function (error) {
  throw error
})
