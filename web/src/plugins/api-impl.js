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
    authenticateWithDummy(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'authenticateWithDummy');
      return fetch(endpoint + basePath + '/auth/dummy'
        , {
          method: 'POST',
          headers,
          mode,
        });
    },
    getDummyAuthenticationResult(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getDummyAuthenticationResult');
      return fetch(endpoint + basePath + '/auth/dummy/session/' + params['sessionToken'] + '/result'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    authenticateWithIRMA(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'authenticateWithIRMA');
      return fetch(endpoint + basePath + '/auth/irma/session'
        , {
          method: 'POST',
          headers,
          mode,
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
    notifyTransferUpdate(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'notifyTransferUpdate');
      return fetch(endpoint + basePath + '/external/transfer/notify'
        , {
          method: 'POST',
          headers,
          mode,
        });
    },
    taskUpdate(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'taskUpdate');
      return fetch(endpoint + basePath + '/internal/customer/' + params['customerID'] + '/task/' + params['taskID'] + ''
        , {
          method: 'PUT',
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
    createDossier(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'createDossier');
      return fetch(endpoint + basePath + '/private/dossier'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    getDossier(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getDossier');
      return fetch(endpoint + basePath + '/private/dossier/' + params['patientID'] + ''
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    createEpisode(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'createEpisode');
      return fetch(endpoint + basePath + '/private/episode'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    getEpisode(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getEpisode');
      return fetch(endpoint + basePath + '/private/episode/' + params['episodeID'] + ''
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    getInbox(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getInbox');
      return fetch(endpoint + basePath + '/private/network/inbox'
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    getInboxInfo(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getInboxInfo');
      return fetch(endpoint + basePath + '/private/network/inbox/info'
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
      return fetch(endpoint + basePath + '/private/patients' + '?' + buildQuery({
          'name': params['name'],
        })

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
    getReports(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getReports');
      return fetch(endpoint + basePath + '/private/reports/' + params['patientID'] + ''
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    createReport(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'createReport');
      return fetch(endpoint + basePath + '/private/reports/' + params['patientID'] + ''
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
    getTransferRequest(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'getTransferRequest');
      return fetch(endpoint + basePath + '/private/transfer-request/' + params['requestorDID'] + '/' + params['fhirTaskID'] + ''
        , {
          method: 'GET',
          headers,
          mode,
        });
    },
    changeTransferRequestState(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'changeTransferRequestState');
      return fetch(endpoint + basePath + '/private/transfer-request/' + params['requestorDID'] + '/' + params['fhirTaskID'] + ''
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    cancelTransfer(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'cancelTransfer');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + ''
        , {
          method: 'DELETE',
          headers,
          mode,
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
    updateTransfer(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'updateTransfer');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + ''
        , {
          method: 'PUT',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    assignTransferDirect(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'assignTransferDirect');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + '/assign'
        , {
          method: 'PUT',
          headers,
          mode,
          body: JSON.stringify(params['body']),

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
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'startTransferNegotiation');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + '/negotiation'
        , {
          method: 'POST',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },
    updateTransferNegotiationStatus(parameters) {
      const params = typeof parameters === 'undefined' ? {} : parameters;
      let headers = {
        'content-type': 'application/json',

      };
      handleSecurity([{"bearerAuth":[]}]
          , headers, params, 'updateTransferNegotiationStatus');
      return fetch(endpoint + basePath + '/private/transfer/' + params['transferID'] + '/negotiation/' + params['negotiationID'] + ''
        , {
          method: 'PUT',
          headers,
          mode,
          body: JSON.stringify(params['body']),

        });
    },

  };
}
