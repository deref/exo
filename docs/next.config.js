/** @type {import('next').NextConfig} */
const withNextra = require("nextra")({
  theme: "nextra-theme-docs",
  themeConfig: "./theme.config.js",
});

module.exports = withNextra({
  reactStrictMode: true,
  redirects: () => {
    return [
      // This redirects /[name].md extension pages to /[name]
      // such that we can write links with ".md" in the source
      // code and navigate between docs pages in both the
      // docs website and the GitHub previewer.
      {
        source: "/:slug.md",
        destination: "/:slug",
        permanent: false,
      },
    ];
  },
});
