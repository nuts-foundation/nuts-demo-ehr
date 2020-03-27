// Set organisation specific info and colours
export default {
  render: () => {
    fetch('/api/organisation/me')
      .then(response => {
        if (response.ok) { return response.json() }
        throw response
      })
      .then(json => {
        const navbar = document.querySelector('nav.navbar')

        navbar.style.backgroundColor = json.colour
        navbar.innerHTML = template(json)

        document.title = json.name
      })
      .catch(reason => {
        if ('status' in reason && reason.status === 403) {
          window.localStorage.setItem('afterLoginReturnUrl', 'dashboard')
          window.location.hash = 'login'
        }
      })
  }
}

const template = (me) => `
  <a class="navbar-brand" href="#">${me.name}</a>
  <span class="navbar-text"><a href="/#public/logout" title="Click to log out">Logged in as ${me.user} <i class="user-icon"></i></a></span>
`
