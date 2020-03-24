import irmaFlow from './irma-flow'
import Thimbleful from 'thimbleful'
let Router

export default {

  load: () => {
    Router = new Thimbleful.Router()

    Router.addRoute('irma', () => {
      openTab('irma')
      irmaFlow.render(
        document.querySelector('#login .card .card-body'),
        {
          header: 'Log in using <i class="irma-web-logo">IRMA</i>',
          cancelled: 'We have not received the signed contract. We\'re sorry, but because of this we can\'t log you in.'
        }
      )
    })

    Thimbleful.Click.instance().register('#login form button', e => {
      e.preventDefault()
      logIn(document.getElementById('username').value)
    })
  },

  render: async (subroute, evnt) => {
    document.getElementById('login').innerHTML = template()
    Router.route(subroute, evnt)
  }

}

function logIn(username) {
  fetch('/api/authentication/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({username})
  })
  .then(() => window.location.hash = 'dashboard')
}

function openTab(tab) {
  document.querySelector('#login .card li a.active').classList.remove('active')
  document.querySelector(`#login .card li a.${tab}`).classList.add('active')
}

const template = () => `
  <h4 class="text-center mb-5">Please log in 🔑</h4>

  <div class="card mx-auto text-center" style="max-width: 33rem;">
    <div class="card-header">
      <ul class="nav nav-tabs card-header-tabs">
        <li class="nav-item">
          <a class="nav-link active" href="#login">Username & Password</a>
        </li>
        <li class="nav-item">
          <a class="nav-link irma" href="#login/irma">IRMA app</a>
        </li>
      </ul>
    </div>
    <div class="card-body">
      <form>
        <p><em>(Any credentials will do 😉)</em></p>
        <div class="form-group row">
          <label for="username" class="col-sm-5 col-form-label text-right">Username:</label>
          <div class="col-sm-7 text-left">
            <input type="text" name="username" id="username"/>
          </div>
        </div>
        <div class="form-group row">
          <label for="password" class="col-sm-5 col-form-label text-right">Password:</label>
          <div class="col-sm-7 text-left">
            <input type="password" name="password" id="password"/>
          </div>
        </div>
        <button class="btn btn-primary">Log in</button>
      </form>
    </div>
  </div>
`
