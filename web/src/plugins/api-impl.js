export default createApi;
function createApi(options) {
  const basePath = '/web';
  const endpoint = options.endpoint || '';
  const cors = !!options.cors;
  const mode = cors ? 'cors' : 'basic';
  const securityHandlers = options.securityHandlers || {};
  const handleSecurity = (security, headers, params, operationId) => {
    for (let i = 0, ilen = security.length; i < ilen; i++) {
      let scheme = security[i];
      let schemeParts = Object.keys(scheme);
      for (let j = 0, jlen = schemeParts.length; j < jlen; j++) {
        let schemePart = schemeParts[j];
        let fulfilsSecurityRequirements = securityHandlers[schemePart](
            headers, params, schemePart);
        if (fulfilsSecurityRequirements) {
          return;
        }

      }
    }
    throw new Error('No security scheme was fulfilled by the provided securityHandlers for operation ' + operationId);
  };
  const ensureRequiredSecurityHandlersExist = () => {
    let requiredSecurityHandlers = ['bearerAuth'];
    for (let i = 0, ilen = requiredSecurityHandlers.length; i < ilen; i++) {
      let requiredSecurityHandler = requiredSecurityHandlers[i];
      if (typeof securityHandlers[requiredSecurityHandler] !== 'function') {
        throw new Error('Expected to see a security handler for scheme "' +
            requiredSecurityHandler + '" in options.securityHandlers');
      }
    }
  };
  ensureRequiredSecurityHandlersExist();
  const buildQuery = (obj) => {
    return Object.keys(obj)
      .filter(key => typeof obj[key] !== 'undefined')
      .map((key) => {
        const value = obj[key];
        if (value === undefined) {
          return '';
        }
        if (value === null) {
          return key;
        }
        if (Array.isArray(value)) {
          if (value.length) {
            return key + '=' + value.map(encodeURIComponent).join('&' + key + '=');
          } else {
            return '';
          }
        } else {
          return key + '=' + encodeURIComponent(value);
        }
      }).join('&');
    };
  return {
    setCustomer(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'setCustomer');
      return fetch(endpoint + basePath + '/auth'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    authenticateWithIRMA(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'authenticateWithIRMA');
      return fetch(endpoint + basePath + '/auth/irma/session'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    getIRMAAuthenticationResult(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getIRMAAuthenticationResult');
      return fetch(endpoint + basePath + '/auth/irma/session/' + params['sessionToken'] + '/result'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    authenticateWithPassword(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'authenticateWithPassword');
      return fetch(endpoint + basePath + '/auth/passwd'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    listCustomers(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'listCustomers');
      return fetch(endpoint + basePath + '/customers'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    checkSession(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'checkSession');
      return fetch(endpoint + basePath + '/private'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    getCustomer(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getCustomer');
      return fetch(endpoint + basePath + '/private/customer'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    getDossier(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getDossier');
      return fetch(endpoint + basePath + '/private/dossier'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    searchOrganizations(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'searchOrganizations');
      return fetch(endpoint + basePath + '/private/network/organizations' + '?' + buildQuery({
          'query': params['query'],
          'didServiceType': params['didServiceType'],
        })

        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    getPatient(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getPatient');
      return fetch(endpoint + basePath + '/private/patient/' + params['patientID'] + ''
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    updatePatient(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'updatePatient');
      return fetch(endpoint + basePath + '/private/patient/' + params['patientID'] + ''
        , {
          method: 'PUT',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    getPatients(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getPatients');
      return fetch(endpoint + basePath + '/private/patients'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    newPatient(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'newPatient');
      return fetch(endpoint + basePath + '/private/patients'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    getPatientTransfers(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getPatientTransfers');
      return fetch(endpoint + basePath + '/private/transfer' + '?' + buildQuery({
          'patientID': params['patientID'],
        })

        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    createTransfer(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'createTransfer');
      return fetch(endpoint + basePath + '/private/transfer'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    getTransfer(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getTransfer');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + ''
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    listTransferNegotiations(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'listTransferNegotiations');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + '/negotiation'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    startTransferNegotiation(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'startTransferNegotiation');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + '/negotiation/' + params['organizationDID'] + ''
        , {
          method: 'POST',
          headers,
          mode,
        });
    },
    assignTransferNegotiation(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'assignTransferNegotiation');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + '/negotiation/' + params['organizationDID'] + '/assign'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },

  };
}
