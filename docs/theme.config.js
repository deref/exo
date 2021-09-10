export default {
  github: "https://github.com/deref/exo", // GitHub link in the navbar
  docsRepositoryBase: "https://github.com/deref/exo/blob/main/docs/pages", // base URL for the docs repository
  titleSuffix: " – exo docs",
  nextLinks: true,
  prevLinks: true,
  search: true,
  customSearch: null, // customizable, you can use algolia for example
  darkMode: true,
  footer: true,
  footerText: `© ${new Date().getFullYear()} Deref Inc.`,
  footerEditLink: `Edit this page on GitHub`,
  logo: (
    <>
      {/* <svg>...</svg> */}
      <span>exo</span>
    </>
  ),
  head: (
    <>
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta
        name="description"
        content="exo: process manager &amp; log viewer for dev"
      />
      <meta
        name="og:title"
        content="exo: process manager &amp; log viewer for dev"
      />
    </>
  ),
};
