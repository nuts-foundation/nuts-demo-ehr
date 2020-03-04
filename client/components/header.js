// Set organisation specific info and colours
export default {
  render: () => {
    fetch('/api/organisation/me')
      .then(response => response.json())
      .then(json => {
        const navbar = document.querySelector('nav.navbar')

        navbar.style.backgroundColor = json.colour
        navbar.innerHTML = template(json)

        document.title = json.name
      })
  }
}

const template = (me) => `
  <a class="navbar-brand" href="#">${me.name}</a>
  <span class="navbar-text">Logged in as ${me.user} <i class="user-icon"></i></span>
`
