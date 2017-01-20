const express = require('express')
const app = express()

const sections = [
  '/login',
  '/signup',
  '/home',
  '/stacks',
  '/topics',
  '/stackEdit',
  '/functions',
  '/kv',
  '/users',
  '/organizations',
]

app.use(express.static('public'))
app.use(express.static('node_modules/codemirror/lib'))
app.use('/dist', express.static('dist'))

for (const section of sections) {
  app.use(section, express.static('public'))
}

app.listen(3000)
