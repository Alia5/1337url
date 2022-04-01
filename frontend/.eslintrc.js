module.exports = {
    root: true,
    env: {
        node: true,
        browser: true
    },
    extends: [
        'plugin:vue/vue3-essential',
        '@vue/airbnb',
        '@vue/typescript/recommended'
    ],
    parserOptions: {
        ecmaVersion: 2020
    },
    rules: {
        'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'warn',
        'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'warn',
        'max-len': ['error', { code: 140 }],
        '@typescript-eslint/no-explicit-any': 'error',
        '@typescript-eslint/explicit-member-accessibility': ['error'],
        'lines-between-class-members': 'off',
        '@typescript-eslint/lines-between-class-members': ['off'],
        indent: ['error', 4],
        'comma-dangle': ['error', {
            arrays: 'never',
            objects: 'never',
            imports: 'never',
            exports: 'never',
            functions: 'never'
        }],
        'function-paren-newline': 'off',
        'import/extensions': ['error', {
            vue: 'always',
            json: 'always'
        }],
        'import/prefer-default-export': 'off'
    }
};
