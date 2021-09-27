const fs = require("fs");

// Run `/path/to/exo $ node script/gen-icon-list.js` to generate output.

const iconDirectory = "mono";
const iconNames = fs
  .readdirSync(`./gui/src/components/${iconDirectory}`)
  .map((x) => x.substr(0, x.length - "SVG.svelte".length));

// Use maps to format the file.
const out = `<script lang="ts" context="module">
  // GENERATED FILE. DO NOT EDIT.

  // Note, these are internal named SVG icons. Dynamic SVGs should be handled separately.

  export type IconGlyph =
${iconNames.map((x) => `    | '${x}'`).join("\n")};
</script>

<script lang="ts">
${iconNames
  .map((x) => `  import ${x}SVG from './${iconDirectory}/${x}SVG.svelte';`)
  .join("\n")}

  export let glyph: IconGlyph;
</script>

{#${iconNames
  .map((x) => `if glyph === '${x}'}\n  <${x}SVG />`)
  .join("\n{:else ")}
{:else}
  <LayersSVG />
{/if}
`;

// Save generated output.
fs.writeFile("./gui/src/components/Icon.svelte", out, function (err) {
  if (err) throw err;
  console.log("Generated Icon.svelte.");
});
