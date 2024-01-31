// eslint.config.js
import parserTs from '@typescript-eslint/parser'
import eslintConfigPrettier from 'eslint-config-prettier'

export default [
  eslintConfigPrettier,
  {
    languageOptions: {
      parser: parserTs,
      parserOptions: {
        ecmaFeatures: { modules: true },
        ecmaVersion: 'latest',
        project: './tsconfig.json',
      },
    },
    plugins: {
      '@typescript-eslint': parserTs,
      parserTs,
    },
    files: ['**/*.ts', '**/*.tsx'],
    ignores: ['dist/**', 'node_modules/**', 'bin/**', 'build/**'],
  },
]
