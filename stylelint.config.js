const kebabCasePattern = /^([a-z][a-z0-9]*)(-[a-z0-9]+)*$/

module.exports = {
  plugins: ['stylelint-scss', 'stylelint-prettier'],
  customSyntax: 'postcss-scss',
  ignoreFiles: [
    'node_modules',
  ],
  rules: {
    'prettier/prettier': true,
    'max-nesting-depth': [6, { ignore: ['pseudo-classes'] }],
    'no-empty-source': null,
    'no-descending-specificity': null,
    'selector-class-pattern': kebabCasePattern,
    'keyframes-name-pattern': kebabCasePattern,
    'at-rule-empty-line-before': [
      'always',
      {
        except: ['first-nested'],
        ignore: [
          'after-comment',
          'blockless-after-blockless',
          'blockless-after-same-name-blockless',
        ],
        ignoreAtRules: ['else'],
      },
    ],
    'at-rule-no-unknown': null,
    'at-rule-no-vendor-prefix': true,
    'rule-empty-line-before': [
      'always',
      {
        except: ['first-nested'],
        ignore: ['after-comment'],
      },
    ],
    'property-no-vendor-prefix': true,
    'scss/at-rule-no-unknown': true,
    'scss/dollar-variable-first-in-block': [true, { ignore: ['comments', 'imports'] }],
    'scss/dollar-variable-colon-space-after': 'always',
    'scss/dollar-variable-pattern': kebabCasePattern,
    'scss/at-if-closing-brace-newline-after': 'always-last-in-chain',
  },
}
