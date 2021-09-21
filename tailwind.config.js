module.exports = {
  purge: ['./web/src/index.html', './web/src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'media', // or 'media' or 'class'
  theme: {
    extend: {},
  },
  important: true,
  variants: {
    extend: {
      opacity: ['disabled'],
    },
  },
  plugins: [],
}
