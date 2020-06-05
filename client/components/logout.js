export default {
  render: async () => {
    return fetch('/api/authentication/logout')
      .then(response => {
        if (!response.ok || response.status !== 204)
          throw Error('Error logging you out' + response)

        window.location.hash = 'public/login'
      })
  }
}
