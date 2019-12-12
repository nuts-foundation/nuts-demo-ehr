/******/ (function(modules) { // webpackBootstrap
/******/ 	// The module cache
/******/ 	var installedModules = {};
/******/
/******/ 	// The require function
/******/ 	function __webpack_require__(moduleId) {
/******/
/******/ 		// Check if module is in cache
/******/ 		if(installedModules[moduleId]) {
/******/ 			return installedModules[moduleId].exports;
/******/ 		}
/******/ 		// Create a new module (and put it into the cache)
/******/ 		var module = installedModules[moduleId] = {
/******/ 			i: moduleId,
/******/ 			l: false,
/******/ 			exports: {}
/******/ 		};
/******/
/******/ 		// Execute the module function
/******/ 		modules[moduleId].call(module.exports, module, module.exports, __webpack_require__);
/******/
/******/ 		// Flag the module as loaded
/******/ 		module.l = true;
/******/
/******/ 		// Return the exports of the module
/******/ 		return module.exports;
/******/ 	}
/******/
/******/
/******/ 	// expose the modules object (__webpack_modules__)
/******/ 	__webpack_require__.m = modules;
/******/
/******/ 	// expose the module cache
/******/ 	__webpack_require__.c = installedModules;
/******/
/******/ 	// define getter function for harmony exports
/******/ 	__webpack_require__.d = function(exports, name, getter) {
/******/ 		if(!__webpack_require__.o(exports, name)) {
/******/ 			Object.defineProperty(exports, name, { enumerable: true, get: getter });
/******/ 		}
/******/ 	};
/******/
/******/ 	// define __esModule on exports
/******/ 	__webpack_require__.r = function(exports) {
/******/ 		if(typeof Symbol !== 'undefined' && Symbol.toStringTag) {
/******/ 			Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });
/******/ 		}
/******/ 		Object.defineProperty(exports, '__esModule', { value: true });
/******/ 	};
/******/
/******/ 	// create a fake namespace object
/******/ 	// mode & 1: value is a module id, require it
/******/ 	// mode & 2: merge all properties of value into the ns
/******/ 	// mode & 4: return value when already ns object
/******/ 	// mode & 8|1: behave like require
/******/ 	__webpack_require__.t = function(value, mode) {
/******/ 		if(mode & 1) value = __webpack_require__(value);
/******/ 		if(mode & 8) return value;
/******/ 		if((mode & 4) && typeof value === 'object' && value && value.__esModule) return value;
/******/ 		var ns = Object.create(null);
/******/ 		__webpack_require__.r(ns);
/******/ 		Object.defineProperty(ns, 'default', { enumerable: true, value: value });
/******/ 		if(mode & 2 && typeof value != 'string') for(var key in value) __webpack_require__.d(ns, key, function(key) { return value[key]; }.bind(null, key));
/******/ 		return ns;
/******/ 	};
/******/
/******/ 	// getDefaultExport function for compatibility with non-harmony modules
/******/ 	__webpack_require__.n = function(module) {
/******/ 		var getter = module && module.__esModule ?
/******/ 			function getDefault() { return module['default']; } :
/******/ 			function getModuleExports() { return module; };
/******/ 		__webpack_require__.d(getter, 'a', getter);
/******/ 		return getter;
/******/ 	};
/******/
/******/ 	// Object.prototype.hasOwnProperty.call
/******/ 	__webpack_require__.o = function(object, property) { return Object.prototype.hasOwnProperty.call(object, property); };
/******/
/******/ 	// __webpack_public_path__
/******/ 	__webpack_require__.p = "";
/******/
/******/
/******/ 	// Load entry module and return exports
/******/ 	return __webpack_require__(__webpack_require__.s = "./client/index.js");
/******/ })
/************************************************************************/
/******/ ({

/***/ "./client/components/header.js":
/*!*************************************!*\
  !*** ./client/components/header.js ***!
  \*************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n// Set organisation specific info and colours\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  render: () => {\n    fetch('/api/organisation/me')\n    .then(response => response.json())\n    .then(json => {\n      const navbar = document.querySelector('nav.navbar');\n\n      navbar.style.backgroundColor = json.colour;\n      navbar.innerHTML = template(json);\n\n      document.title = json.name;\n    });\n  }\n});\n\nconst template = (me) => `\n  <a class=\"navbar-brand\" href=\"#\">${me.name}</a>\n  <span class=\"navbar-text\">Logged in as ${me.user} <i class=\"user-icon\"></i></span>\n`;\n\n\n//# sourceURL=webpack:///./client/components/header.js?");

/***/ }),

/***/ "./client/components/patient-overview.js":
/*!***********************************************!*\
  !*** ./client/components/patient-overview.js ***!
  \***********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  // Fetch all patients and render to the table\n  render: () => {\n    return fetch('/api/patient/all')\n    .then(response => response.json())\n    .then(json => {\n      document.getElementById('patient-overview').innerHTML = template(json);\n    });\n  }\n});\n\nconst template = (patients) => `\n  <h2>Patients in care</h2>\n\n  <table class=\"table table-borderless table-bordered table-hover\">\n\n    <thead class=\"thead-dark\">\n      <tr>\n        <th>Last name</th>\n        <th>First name</th>\n      </tr>\n    </thead>\n\n    <tbody>\n      ${patients.map(patient => `\n        <tr>\n          <td><a href='#patient-details/${patient.id}'>${patient.name.family}</a></td>\n          <td>${patient.name.given}</td>\n        </tr>\n      `).join('')}\n    </tbody>\n\n  </table>\n`;\n\n\n//# sourceURL=webpack:///./client/components/patient-overview.js?");

/***/ }),

/***/ "./client/components/patient/details.js":
/*!**********************************************!*\
  !*** ./client/components/patient/details.js ***!
  \**********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  render: (patient) => {\n    document.getElementById('patient-details').innerHTML = template(patient);\n  }\n});\n\nconst template = (patient) => `\n  <table class=\"table table-borderless table-bordered\">\n    <tbody>\n      <tr>\n        <th>BSN</th>\n        <td>${patient.bsn}</td>\n      </tr>\n      <tr>\n        <th>Last name</th>\n        <td>${patient.name.family}</td>\n      </tr>\n      <tr>\n        <th>First name</th>\n        <td>${patient.name.given}</td>\n      </tr>\n      <tr>\n        <th>Date of birth</th>\n        <td>${patient.birthDate}</td>\n      </tr>\n      <tr>\n        <th>Gender</th>\n        <td>${patient.gender}</td>\n      </tr>\n    </tbody>\n  </table>\n`;\n\n\n//# sourceURL=webpack:///./client/components/patient/details.js?");

/***/ }),

/***/ "./client/components/patient/network.js":
/*!**********************************************!*\
  !*** ./client/components/patient/network.js ***!
  \**********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var thimbleful__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! thimbleful */ \"./node_modules/thimbleful/src/index.js\");\n\n\nlet patientId;\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  render: async (patient) => {\n    patientId = patient.id;\n    await renderNetwork();\n\n    document.getElementById('patient-add-consent-org').addEventListener('keyup', (e) => {\n      const input = e.target.value;\n      if ( input.length >= 2 ) {\n        fetch(`/api/organisation/search/${input}`)\n        .then(response => response.json())\n        .then(result => renderAutoComplete(result));\n      }\n      if ( input.length == 0 ) {\n        renderAutoComplete([]);\n      }\n    });\n  }\n});\n\nasync function renderNetwork() {\n  const received = await fetch(`/api/consent/${patientId}/received`).then(r => r.json());\n  const given    = await fetch(`/api/consent/${patientId}/given`).then(r => r.json());\n\n  document.getElementById('patient-network').innerHTML = template(received, given);\n}\n\nfunction renderAutoComplete(results) {\n  document.getElementById('patient-consent-auto-complete')\n          .innerHTML = results.map(result => `\n    <a class=\"list-group-item list-group-item-action d-flex justify-content-between align-items-center\"\n          data-organisation-id=\"${result.identifier}\" data-organisation-name=\"${result.name}\">\n      ${result.name}\n      <span class=\"badge badge-primary badge-pill\">${result.identifier.split(':').pop()}</span>\n    </a>\n  `).join('');\n}\n\n// Selecting an option from the auto-complete dropdown\nthimbleful__WEBPACK_IMPORTED_MODULE_0__[\"default\"].Click.instance().register('a[data-organisation-id]', (e) => {\n  const id = e.target.attributes['data-organisation-id'].value;\n  const name = e.target.attributes['data-organisation-name'].value;\n  document.querySelector('input[name=\"organisation-id\"]').value = id;\n  document.getElementById('patient-add-consent-org').value = name;\n  document.getElementById('patient-consent-auto-complete').innerHTML = '';\n});\n\n// Storing new consent\nthimbleful__WEBPACK_IMPORTED_MODULE_0__[\"default\"].Click.instance().register('#patient-add-consent-button', (e) => {\n  const organisationURN = document.querySelector('input[name=\"organisation-id\"]').value;\n  const reason = document.getElementById('patient-add-consent-reason').value;\n\n  storeConsent({ organisationURN, reason })\n  .then(() => renderNetwork());\n});\n\nfunction storeConsent(consent) {\n  return fetch(`/api/consent/${patientId}`, {\n    method: 'PUT',\n    headers: {\n      'Content-Type': 'application/json'\n    },\n    body: JSON.stringify(consent)\n  })\n  .then(response => {\n    if ( response.status != 201 ) throw 'Error storing observation';\n  });\n}\n\nconst template = (receivedConsents, givenConsents) => `\n  &nbsp;\n\n  <div class=\"card\">\n    <div class=\"card-body\">\n      <p>Organisations that have shared information with you:</p>\n      <ul>\n      ${receivedConsents.length > 0 ? receivedConsents.map(consent => `\n        <li><a href=\"#patient-network/${patientId}/${consent.identifier}\">${consent.name}</a></li>\n      `).join('') : '<li><i>None</i></li>'}\n      </ul>\n    </div>\n  </div>\n\n  &nbsp;\n\n  <div class=\"card\">\n    <div class=\"card-body\">\n      <p>Organisations you're sharing information with:</p>\n      <ul>\n      ${givenConsents.length > 0 ? givenConsents.map(consent => `\n        <li>${consent.name}</li>\n      `).join('') : '<li><i>None</i></li>'}\n      </ul>\n\n      <p><button class=\"btn btn-primary\" data-toggle=\"#patient-add-consent\">Add</button></p>\n\n      <section id=\"patient-add-consent\" class=\"page\">\n        <div class=\"card\">\n          <div class=\"card-body\">\n            <form id=\"patient-consent-form\">\n\n              <p>Share your information about this patient with another organisation:\n              <div class=\"form-group row\">\n                <label for=\"patient-add-consent-org\" class=\"col-sm-3 col-form-label\">Organisation:</label>\n                <div class=\"col-sm-9\">\n                  <input type=\"hidden\" name=\"organisation-id\"/>\n                  <input type=\"text\" class=\"form-control\" id=\"patient-add-consent-org\" placeholder=\"Organisation name\" autocomplete=\"off\"/>\n                  <div class=\"list-group auto-complete\" id=\"patient-consent-auto-complete\"></div>\n                </div>\n              </div>\n              <div class=\"form-group row\">\n                <label for=\"patient-add-consent-reason\" class=\"col-sm-3 col-form-label\">Legal basis:</label>\n                <div class=\"col-sm-9\">\n                  <input type=\"text\" name=\"reason\" class=\"form-control\" id=\"patient-add-consent-reason\" placeholder=\"Your legal basis for sharing this information\"/>\n                </div>\n              </div>\n              <button id=\"patient-add-consent-button\" type=\"button\" class=\"btn btn-primary float-right\">Share</button>\n\n            </form>\n          </div>\n        </div>\n      </section>\n    </div>\n  </div>\n`;\n\n\n//# sourceURL=webpack:///./client/components/patient/network.js?");

/***/ }),

/***/ "./client/components/patient/observations.js":
/*!***************************************************!*\
  !*** ./client/components/patient/observations.js ***!
  \***************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var thimbleful__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! thimbleful */ \"./node_modules/thimbleful/src/index.js\");\n\n\nlet patientId;\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  render: (patient) => {\n    patientId = patient.id;\n    renderObservations(patientId);\n  }\n});\n\nfunction renderObservations(patientId) {\n  return fetch(`/api/observation/byPatientId/${patientId}`)\n  .then(response => response.json())\n  .then(observations => {\n    document.getElementById('patient-observations').innerHTML = template(observations);\n  });\n}\n\nconst template = (observations) => `\n  &nbsp;\n\n  ${observations.map(observation => `\n    <div class=\"card\"><div class=\"card-body\">\n      <code>${observation.timestamp}</code>\n      <p>${observation.content}</p>\n    </div></div>\n    &nbsp;\n  `).join('')}\n\n  &nbsp;\n\n  <h4>New observation</h4>\n\n  &nbsp;\n\n  <textarea class=\"form-control\" rows=\"5\" id=\"new-observation\"></textarea>\n  <button class=\"btn btn-primary float-right\" id=\"add-observation\">Save</button>\n`;\n\n// Add click handler for storing new observations\nthimbleful__WEBPACK_IMPORTED_MODULE_0__[\"default\"].Click.instance().register('button#add-observation', (e) => {\n  const ta = document.getElementById('new-observation');\n\n  storeObservation({\n    patientId: patientId,\n    content: ta.value\n  })\n  .then(() => {\n    ta.value = '';\n    renderObservations(patientId);\n  });\n});\n\nfunction storeObservation(observation) {\n  return fetch('/api/observation', {\n    method: 'PUT',\n    headers: {\n      'Content-Type': 'application/json'\n    },\n    body: JSON.stringify(observation)\n  })\n  .then(response => {\n    if ( response.status != 201 ) throw 'Error storing observation';\n  });\n}\n\n\n//# sourceURL=webpack:///./client/components/patient/observations.js?");

/***/ }),

/***/ "./client/components/patient/patient.js":
/*!**********************************************!*\
  !*** ./client/components/patient/patient.js ***!
  \**********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _details__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./details */ \"./client/components/patient/details.js\");\n/* harmony import */ var _observations__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./observations */ \"./client/components/patient/observations.js\");\n/* harmony import */ var _network__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./network */ \"./client/components/patient/network.js\");\n\n\n\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  render: (patientId) => {\n    return fetch(`/api/patient/${patientId}`)\n    .then(response => response.json())\n    .then(patient => {\n      document.getElementById('patient').innerHTML = template(patient);\n\n      // Render child components\n      _details__WEBPACK_IMPORTED_MODULE_0__[\"default\"].render(patient);\n      _observations__WEBPACK_IMPORTED_MODULE_1__[\"default\"].render(patient);\n      _network__WEBPACK_IMPORTED_MODULE_2__[\"default\"].render(patient);\n    });\n  }\n});\n\nconst template = (patient) => `\n  <nav aria-label=\"breadcrumb\">\n    <ol class=\"breadcrumb\">\n      <li class=\"breadcrumb-item\"><a href=\"#patient-overview\">Patients in care</a></li>\n      <li class=\"breadcrumb-item active\" aria-current=\"page\">${patient.name.given} ${patient.name.family}</li>\n    </ol>\n  </nav>\n\n  <ul class=\"nav nav-tabs\">\n    <li class=\"nav-item\">\n      <a class=\"nav-link active\" data-open=\"#patient-details\">Details</a>\n    </li>\n    <li class=\"nav-item\">\n      <a class=\"nav-link\" data-open=\"#patient-observations\">Observations</a>\n    </li>\n    <li class=\"nav-item\">\n      <a class=\"nav-link\" data-open=\"#patient-network\">Network</a>\n    </li>\n  </ul>\n\n  <section class=\"tab-pane active\" id=\"patient-details\" data-group=\"patient-tab-panes\" data-follower=\"a[data-open='#patient-details']\"></section>\n  <section class=\"tab-pane\" id=\"patient-observations\" data-group=\"patient-tab-panes\" data-follower=\"a[data-open='#patient-observations']\"></section>\n  <section class=\"tab-pane\" id=\"patient-network\" data-group=\"patient-tab-panes\" data-follower=\"a[data-open='#patient-network']\"></section>\n`;\n\n\n//# sourceURL=webpack:///./client/components/patient/patient.js?");

/***/ }),

/***/ "./client/index.js":
/*!*************************!*\
  !*** ./client/index.js ***!
  \*************************/
/*! no exports provided */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var thimbleful__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! thimbleful */ \"./node_modules/thimbleful/src/index.js\");\n/* harmony import */ var _components_header__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./components/header */ \"./client/components/header.js\");\n/* harmony import */ var _routing__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./routing */ \"./client/routing.js\");\n\n\n\n\n// Render organisation name, colour and user\n_components_header__WEBPACK_IMPORTED_MODULE_1__[\"default\"].render();\n\n// Load the routes\n_routing__WEBPACK_IMPORTED_MODULE_2__[\"default\"].load();\n\n// Enable data attributes for interface components\nnew thimbleful__WEBPACK_IMPORTED_MODULE_0__[\"default\"].Energize(\"#app\");\n\n\n//# sourceURL=webpack:///./client/index.js?");

/***/ }),

/***/ "./client/routing.js":
/*!***************************!*\
  !*** ./client/routing.js ***!
  \***************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var thimbleful__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! thimbleful */ \"./node_modules/thimbleful/src/index.js\");\n/* harmony import */ var _components_patient_overview__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./components/patient-overview */ \"./client/components/patient-overview.js\");\n/* harmony import */ var _components_patient_patient__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./components/patient/patient */ \"./client/components/patient/patient.js\");\n\nconst router = new thimbleful__WEBPACK_IMPORTED_MODULE_0__[\"default\"].Router();\n\n\n\n\n/* harmony default export */ __webpack_exports__[\"default\"] = ({\n  load: () => {\n    // Root redirects to patient overview\n    if ( !window.location.hash ) window.location.hash = 'patient-overview';\n\n    router.addRoute('patient-overview', async link => {\n      await _components_patient_overview__WEBPACK_IMPORTED_MODULE_1__[\"default\"].render();\n      openPage(link);\n    });\n\n    router.addRoute(/patient-details\\/(\\d+)(\\/.*)?/, async (link, matches) => {\n      await _components_patient_patient__WEBPACK_IMPORTED_MODULE_2__[\"default\"].render(matches[1]);\n      openPage('patient');\n    });\n  }\n});\n\n// Show the given page, hide others\nfunction openPage(page) {\n  document.querySelector('.page.active').classList.remove('active');\n  document.querySelector(`#${page}`).classList.add('active');\n  window.scrollTo(0,0);\n}\n\n\n//# sourceURL=webpack:///./client/routing.js?");

/***/ }),

/***/ "./node_modules/thimbleful/src/click.js":
/*!**********************************************!*\
  !*** ./node_modules/thimbleful/src/click.js ***!
  \**********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/*\n * This class installs one single click handler on the whole document, and\n * evaluates which callback to call at click time, based on the element that has\n * been clicked. This allows us to swap out and rerender whole sections of the\n * DOM without having to reinstall a bunch of click handlers each time. This\n * nicely decouples the render logic from the click event management logic.\n *\n * To make sure we really only install a single click handler, you can use the\n * singleton pattern and ask for `Click.instance()` instead of creating a new\n * object.\n */\n\nclass Click {\n\n  constructor() {\n    this._handlers = {};\n\n    document.addEventListener('click',     (e) => this._callHandler('click',     e));\n    document.addEventListener('mousedown', (e) => this._callHandler('mousedown', e));\n    document.addEventListener('mouseup',   (e) => this._callHandler('mouseup',   e));\n  }\n\n  register(selector, handlers = {click: null, mousedown: null, mouseup: null}) {\n    if (typeof handlers == 'function') handlers = { click: handlers };\n    this._handlers[selector] = this._handlers[selector] || [];\n    this._handlers[selector].push(handlers);\n  }\n\n  _callHandler(type, e) {\n    Object.keys(this._handlers).forEach((selector) => {\n      if (e.target.closest(selector) !== null) {\n        const handlers = this._handlers[selector].map((h) => h[type]);\n        handlers.forEach((handler) => {\n          if (typeof handler == 'function' && !e.defaultPrevented)\n            handler(e, selector)\n        });\n      }\n    });\n  }\n\n}\n\nClick.instance = function() {\n  if (!!Click._instance) return Click._instance;\n  return Click._instance = new Click();\n}\n\n/* harmony default export */ __webpack_exports__[\"default\"] = (Click);\n\n\n//# sourceURL=webpack:///./node_modules/thimbleful/src/click.js?");

/***/ }),

/***/ "./node_modules/thimbleful/src/energize.js":
/*!*************************************************!*\
  !*** ./node_modules/thimbleful/src/energize.js ***!
  \*************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"default\", function() { return Energize; });\n/* harmony import */ var _click__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./click */ \"./node_modules/thimbleful/src/click.js\");\n/**\n * Given a scope this class adds a bunch of behaviour to elements that\n * you define through data attributes. This behaviour is based around adding\n * or removing an 'active' class when elements are clicked:\n *\n *  - data-open — A selector to put the 'active' class on when clicked\n *  - data-close — A selector to remove the 'active' class from when clicked\n *  - data-toggle — A selector to toggle the 'active' class on when clicked\n *  - data-group — If I get the 'active' class, remove it from others in my group\n *  - data-timer — If I get the 'active' class, remove it again after this many milliseconds\n *  - data-follower — A selector for another element that follows my behaviour\n *\n * If you wish, you can override the class name and the names of all the\n * attributes as options to the constructor.\n */\n\n\n\nclass Energize {\n\n  constructor(scope, options = {}) {\n    this._scope   = scope;\n    this._options = this._normalizeOptions(options);\n\n    _click__WEBPACK_IMPORTED_MODULE_0__[\"default\"].instance().register(`${scope} [${this._options.open}], ${scope} [${this._options.close}], ${scope} [${this._options.toggle}]`, (e) => this._handleClick(e));\n  }\n\n  _normalizeOptions(options) {\n    return Object.assign({\n      class:    'active',\n      open:     'data-open',\n      close:    'data-close',\n      toggle:   'data-toggle',\n      group:    'data-group',\n      timer:    'data-timer',\n      follower: 'data-follower'\n    }, options);\n  }\n\n  _handleClick(evnt) {\n    // Which element did we click?\n    const target = evnt.target.closest(`[${this._options.open}], [${this._options.close}], [${this._options.toggle}]`);\n\n    // What does the clicked element wish to open, close or toggle?\n    const closeSelector  = target.getAttribute(this._options.close);\n    const openSelector   = target.getAttribute(this._options.open);\n    const toggleSelector = target.getAttribute(this._options.toggle);\n\n    let closeElements = closeSelector ? document.querySelectorAll(`${this._scope} ${closeSelector}`)  : [];\n    let openElements  =  openSelector ? document.querySelectorAll(`${this._scope} ${openSelector}`)   : [];\n\n    // Add elements that need to be toggled\n    closeElements = [...closeElements, ...(toggleSelector ? document.querySelectorAll(`${this._scope} ${toggleSelector}.${this._options.class}`)       : [])];\n    openElements  = [...openElements,  ...(toggleSelector ? document.querySelectorAll(`${this._scope} ${toggleSelector}:not(.${this._options.class})`) : [])];\n\n    this._close(closeElements);\n    this._open(openElements);\n\n    // We're done with this event, don't try to evaluate it any further\n    evnt.preventDefault();\n    evnt.stopPropagation();\n  }\n\n  _close(elements) {\n    elements.forEach((element) => {\n      element.classList.remove(this._options.class);\n      this._close(this._followers(element));\n    });\n  }\n\n  _open(elements) {\n    elements.forEach((element) => {\n      this._close(this._group(element));\n      element.classList.add(this._options.class);\n      this._open(this._followers(element));\n\n      // Set self-destruct timer if needed\n      const delay = element.getAttribute(this._options.timer);\n      if (delay) window.setTimeout(() => this._close([element]), delay);\n    });\n  }\n\n  _group(element) {\n    const group = element.getAttribute(this._options.group);\n    if (!group) return [];\n    return [...document.querySelectorAll(`${this._scope} [${this._options.group}=${group}]`)];\n  }\n\n  _followers(element) {\n    const selector = element.getAttribute(this._options.follower);\n    if (!selector) return [];\n    return [...document.querySelectorAll(`${this._scope} ${selector}`)];\n  }\n\n}\n\n\n//# sourceURL=webpack:///./node_modules/thimbleful/src/energize.js?");

/***/ }),

/***/ "./node_modules/thimbleful/src/filetarget.js":
/*!***************************************************!*\
  !*** ./node_modules/thimbleful/src/filetarget.js ***!
  \***************************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _click__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./click */ \"./node_modules/thimbleful/src/click.js\");\n/*\n * This class installs single drag event handlers on the whole document, and\n * evaluates which element they influence at drag time. If you drop a file the\n * relevant callback gets called, based on the element that the file was dropped\n * on. This allows us to swap out and rerender whole sections of the DOM without\n * having to reinstall a bunch of event handlers each time. This nicely\n * decouples the render logic from the drag event management logic.\n *\n * To make sure we really only install single handlers, you can use the\n * singleton pattern and ask for `FileTarget.instance()` instead of creating a new\n * object.\n */\n\n\n\nclass FileTarget {\n\n  constructor(dragClass = 'dragging') {\n    this._dragClass = dragClass;\n    this._handlers  = {};\n\n    document.addEventListener('dragover',  (e) => this._dragOver(e));\n    document.addEventListener('dragleave', (e) => this._dragLeave(e));\n    document.addEventListener('drop',      (e) => this._drop(e));\n  }\n\n  register(selector, callback) {\n    this._handlers[selector] = callback;\n    _click__WEBPACK_IMPORTED_MODULE_0__[\"default\"].instance().register(selector, (e, s) => this._openFileDialog(e, s));\n  }\n\n  _dragOver(e) {\n    if (!this._isDropTarget(e.target)) return;\n    e.stopPropagation();\n    e.preventDefault();\n    e.dataTransfer.dropEffect = 'copy';\n    e.target.classList.add(this._dragClass);\n  }\n\n  _dragLeave(e) {\n    if (!this._isDropTarget(e.target)) return;\n    e.stopPropagation();\n    e.preventDefault();\n    e.target.classList.remove(this._dragClass);\n  }\n\n  _drop(e) {\n    let selector = this._isDropTarget(e.target);\n    if (!selector) return;\n    e.stopPropagation();\n    e.preventDefault();\n    e.target.classList.remove(this._dragClass);\n    this._handleFile(selector, e, e.dataTransfer.files[0]);\n  }\n\n  _isDropTarget(target) {\n    return Object.keys(this._handlers).find((selector) => {\n      if (target.closest(selector)) return selector;\n    }) || false;\n  }\n\n  _openFileDialog(e, selector) {\n    const input = document.createElement('input');\n    input.type  = 'file';\n    input.addEventListener('change', (c) =>\n      this._handleFile(selector, e, c.target.files[0])\n    );\n    input.click();\n  }\n\n  _handleFile(selector, e, file) {\n    this._readFile(file)\n        .then((r) => this._handlers[selector](file, r, e));\n  }\n\n  _readFile(file) {\n    return new Promise((resolve, reject) => {\n      var reader = new FileReader();\n      reader.addEventListener('load', (e) => resolve(e.target.result));\n      reader.readAsDataURL(file);\n    });\n  }\n\n}\n\nFileTarget.instance = function() {\n  if (!!FileTarget._instance) return FileTarget._instance;\n  return FileTarget._instance = new FileTarget();\n}\n\n/* harmony default export */ __webpack_exports__[\"default\"] = (FileTarget);\n\n\n//# sourceURL=webpack:///./node_modules/thimbleful/src/filetarget.js?");

/***/ }),

/***/ "./node_modules/thimbleful/src/index.js":
/*!**********************************************!*\
  !*** ./node_modules/thimbleful/src/index.js ***!
  \**********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony import */ var _click__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./click */ \"./node_modules/thimbleful/src/click.js\");\n/* harmony import */ var _filetarget__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./filetarget */ \"./node_modules/thimbleful/src/filetarget.js\");\n/* harmony import */ var _router__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! ./router */ \"./node_modules/thimbleful/src/router.js\");\n/* harmony import */ var _energize__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! ./energize */ \"./node_modules/thimbleful/src/energize.js\");\n\n\n\n\n\nconst Thimbleful = {\n  Click: _click__WEBPACK_IMPORTED_MODULE_0__[\"default\"], FileTarget: _filetarget__WEBPACK_IMPORTED_MODULE_1__[\"default\"], Router: _router__WEBPACK_IMPORTED_MODULE_2__[\"default\"], Energize: _energize__WEBPACK_IMPORTED_MODULE_3__[\"default\"]\n};\n\n/* harmony default export */ __webpack_exports__[\"default\"] = (Thimbleful);\nwindow.Thimbleful = Thimbleful;\n\n\n//# sourceURL=webpack:///./node_modules/thimbleful/src/index.js?");

/***/ }),

/***/ "./node_modules/thimbleful/src/router.js":
/*!***********************************************!*\
  !*** ./node_modules/thimbleful/src/router.js ***!
  \***********************************************/
/*! exports provided: default */
/***/ (function(module, __webpack_exports__, __webpack_require__) {

"use strict";
eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export (binding) */ __webpack_require__.d(__webpack_exports__, \"default\", function() { return Router; });\n/* harmony import */ var _click__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! ./click */ \"./node_modules/thimbleful/src/click.js\");\n/**\n * This is a very simple routing class that listens to location hash changes and\n * clicks on links to registered routes.\n *\n * You have to explicitly define the routes that you wish to use, so we don't\n * clash (too much) with deep-linking to named anchors on your page. And also\n * because it enables you to handle different routes with different functions.\n */\n\n\n\nclass Router {\n\n  constructor() {\n    this._routes = [];\n\n    _click__WEBPACK_IMPORTED_MODULE_0__[\"default\"].instance().register('a[href]',  (e) => this._handleClick(e));\n    window.addEventListener('hashchange', (e) => this._handleNavigationEvent(e));\n    window.addEventListener('load',       (e) => this._handleNavigationEvent(e));\n  }\n\n  addRoute(route, handler) {\n    this._routes.push([route, handler]);\n  }\n\n  addRoutes(routes, handler = null) {\n    if (Array.isArray(routes))\n      routes.forEach((route) => this.addRoute(route, handler));\n    else\n      Object.keys(routes).forEach(route => this.addRoute(route, routes[route]));\n  }\n\n  _handleClick(evnt) {\n    let link = evnt.target.getAttribute('href');\n    if (!this._matchingLink(link)) return;\n    window.location.hash = link;\n    evnt.preventDefault();\n  }\n\n  _handleNavigationEvent(evnt) {\n    let link = window.location.hash;\n    if (!(link = this._matchingLink(link))) return;\n    let handler = link.route[1]\n    if (handler) handler(link.route[0], link.matches, evnt);\n  }\n\n  _matchingLink(hash) {\n    if (!hash) return false;\n    if (!hash.substr(0,1) == \"#\") return false;\n    for (const route of this._routes) {\n      if (route[0] instanceof RegExp) {\n        const matches = hash.substr(1).match(route[0]);\n        if (matches) return {route, matches};\n      } else {\n        if (route[0] == hash.substr(1)) return {route, matches: null}\n      }\n    }\n    return false;\n  }\n\n}\n\n\n//# sourceURL=webpack:///./node_modules/thimbleful/src/router.js?");

/***/ })

/******/ });