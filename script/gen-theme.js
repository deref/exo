var fs = require('fs');

var themeVariables = {
  'primary-color': {
    lite: '#000000',
    dark: '#ffffff',
    blck: '#ffffff',
  },
  'primary-bg-color': {
    lite: '#ffffff',
    dark: '#222222',
    blck: '#000000',
  },
};

var childVariables = {
  'log-color': {
    lite: 'var(--light-log-color)',
    dark: 'var(--dark-log-color)',
    blck: 'var(--dark-log-color)',
  },
  'log-bg-color': {
    lite: 'var(--light-log-bg-color)',
    dark: 'var(--dark-log-bg-color)',
    blck: 'var(--dark-log-bg-color)',
  },
  'log-hover-color': {
    lite: 'var(--light-log-hover-color)',
    dark: 'var(--dark-log-hover-color)',
    blck: 'var(--dark-log-hover-color)',
  },
  'log-bg-hover-color': {
    lite: 'var(--light-log-bg-hover-color)',
    dark: 'var(--dark-log-bg-hover-color)',
    blck: 'var(--dark-log-bg-hover-color)',
  },
};

function themeDefinition(indent, theme) {
  return Object.entries(themeVariables).map((entry) => {return `${indent}--${entry[0]}: ${entry[1][theme]};`}).join(`\n`)
}

function childDefinition(indent, theme) {
  return Object.entries(childVariables).map((entry) => {return `${indent}--${entry[0]}: ${entry[1][theme]};`}).join(`\n`)
}

var out = `/* GENERATED FILE */
body.auto {\n${themeDefinition("  ", "lite")}\n}
body.auto * {\n${childDefinition("  ", "lite")}\n}
@media (prefers-color-scheme: dark) {
  body.auto {\n${themeDefinition("  ", "blck")}
  }
  body.auto * {\n${childDefinition("  ", "blck")}
  }
}
body.light {\n${themeDefinition("  ", "lite")}\n}
body.light * {\n${childDefinition("  ", "lite")}\n}
body.dark {\n${themeDefinition("  ", "dark")}\n}
body.dark * {\n${childDefinition("  ", "dark")}\n}
body.black {\n${themeDefinition("  ", "blck")}\n}
body.black * {\n${childDefinition("  ", "blck")}\n}
`;

fs.writeFile('./gui/public/theme-gen.css', out, function(err) {
  if (err) throw err;
  console.log('Generated theme file.')
});
