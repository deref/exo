import fs from 'fs';

const iconDirectory = 'mono';
const iconNames = fs
  .readdirSync(`./src/components/${iconDirectory}`)
  .map((x) => x.substr(0, x.length - '.svelte'.length));

// Use maps to format the file.
const out = `<script lang="ts" context="module">
  // GENERATED FILE. DO NOT EDIT.

  // Note, these are internal named SVG icons. Dynamic SVGs should be handled separately.

  export type IconGlyph =
${iconNames.map((x) => `    | '${x}'`).join('\n')};
</script>

<script lang="ts">
${iconNames
  .map((x) => `  import ${x}Glyph from './${iconDirectory}/${x}.svelte';`)
  .join('\n')}

  export let glyph: IconGlyph;
</script>

{#${iconNames
  .map((x) => `if glyph === '${x}'}\n  <${x}Glyph />`)
  .join('\n{:else ')}
{:else}
  <LayersGlyph />
{/if}
`;

// Save generated output.
fs.writeFile('./src/components/Icon.svelte', out, function (err) {
  if (err) throw err;
  console.log('Generated Icon.svelte.');
});
