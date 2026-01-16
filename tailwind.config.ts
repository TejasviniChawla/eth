import type { Config } from 'tailwindcss';

const config: Config = {
  content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        canvas: '#0B0F19',
        card: '#111827',
        accent: '#7C3AED'
      }
    }
  },
  plugins: []
};

export default config;
