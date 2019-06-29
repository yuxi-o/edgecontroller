```text
Copyright 2019 Intel Corporation and Smart-Edge.com, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

# 3GPP CUPS Management API Client

## Prerequisites
- Node & NPM installed (v10.15.3, or V10 LTS)
  - recommended to use NVM https://github.com/nvm-sh/nvm to manage your Node versions
- Yarn installed globally `npm install -g yarn`
- Install dependencies via `yarn install` within the project

## Environment Setup

### Development
A development .env under `.env.development` is already configured with the default URLs
for the CUPs UI local development

The local development server is proxied via create-react-app's proxy functionality.
This is to resolve CORS local dev concerns.

### Production

**Any client web browser using the Controller CE web user interface must have network access 
to the listening address and port of the Controller CE REST API.**

## Available Scripts

In the project directory, you can run:

### `yarn install`

Downloads and install all project dependencies defined in the `package.json`
file.

### `yarn start`

Runs the app in the development mode.<br> Open
[http://localhost:3000](http://localhost:3000) to view it in the browser.

The page will reload if you make edits.<br> You will also see any lint errors in
the console.

### `yarn test`

Launches the test runner in the interactive watch mode.<br> See the section
about [running
tests](https://facebook.github.io/create-react-app/docs/running-tests) for more
information.

### `yarn run build`

Builds the app for production to the `build` folder.<br> It correctly bundles
React in production mode and optimizes the build for the best performance.

The build is minified and the filenames include the hashes.<br> Your app is
ready to be deployed!

See the section about
[deployment](https://facebook.github.io/create-react-app/docs/deployment) for
more information.

