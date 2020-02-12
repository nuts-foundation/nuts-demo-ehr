export default function call(url, element) {

  setState(element, 'loading');

  return fetch(url)
  .then(response => {

    if ( response.ok ) {
      setState(element, 'loaded');
      return response.json();
    } else {
      console.info("Response is not OK");
      console.error(response);
      setState(element, 'error');
      return Promise.reject();
    }

  })
  .catch(error => {
    console.info("Response is not OK");
    console.error(error);
    setState(element, 'error');
    return Promise.reject();
  });
}

function setState(element, state) {
  element.classList.remove('loaded');
  element.classList.remove('loading');
  element.classList.remove('error');
  element.classList.add(state);
}
