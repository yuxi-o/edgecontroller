// Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const express = require('express');
const dotenv = require('dotenv');
const nunjucks = require('nunjucks');
const helmet = require('helmet');
const uuidv4 = require('uuid/v4');
const path = require('path');

const port = 80;

const app = express();

dotenv.config();

nunjucks.configure('build', {
  autoescape: true,
  express: app
});

app.use(helmet());

app.use((req, res, next) => {
  res.locals.styleNonce = Buffer.from(uuidv4()).toString('base64');
  next();
});

app.use(
  helmet.contentSecurityPolicy({
    directives: {
      defaultSrc: ["'self'"],
      styleSrc: ["'self'", (req, res) => `'nonce-${res.locals.styleNonce}'`],
      connectSrc: [`${process.env.REACT_APP_CUPS_API}`]
    }
  })
);

app.get('/', (req, res) => {
  res.render('index.html', { styleNonce: res.locals.styleNonce });
});

app.use(express.static(path.join(__dirname, 'build')));

app.get('*', (req, res) => {
  res.render('index.html', { styleNonce: res.locals.styleNonce });
});

app.listen(port, () => console.log(`Server listening on port ${port}.`));
