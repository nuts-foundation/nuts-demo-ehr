const colors = require('tailwindcss/colors')

module.exports = {
  content: ['./web/src/index.html', './web/src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        green: colors.emerald,
        yellow: colors.amber,
        purple: colors.violet,
        gray: colors.neutral,
      },
    },
  },
  important: true,
  plugins: [],
}
