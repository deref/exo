package template

import (
	"fmt"

	"github.com/deref/exo/internal/core/api"
)

var bucketURL = "https://exo-starter-templates.s3.us-west-2.amazonaws.com"

func GetTemplateDescriptions() []api.TemplateDescription {
	// TODO: make this dynamic
	templates := []api.TemplateDescription{
		{Name: "ackee", DisplayName: "Ackee"},
		{Name: "blitzjs", DisplayName: "Blitz.js"},
		{Name: "calendso", DisplayName: "Calendso"},
		{Name: "code-server", DisplayName: "code-server"},
		{Name: "cusdis", DisplayName: "Cusdis"},
		{Name: "deno", DisplayName: "Deno", IconGlyph: "LogoDeno"},
		{Name: "discordjs-typescript", DisplayName: "discord.js Typescript", IconGlyph: "LogoDiscord"},
		{Name: "discordjs", DisplayName: "discord.js", IconGlyph: "LogoDiscord"},
		{Name: "discordpy", DisplayName: "discord.py", IconGlyph: "LogoDiscord"},
		{Name: "djangopy", DisplayName: "DjangoPy"},
		{Name: "elixir-phoenix", DisplayName: "Elixir/Phoenix"},
		{Name: "eris", DisplayName: "Eris"},
		{Name: "expressjs-mongoose", DisplayName: "Express.js + Mongoose", IconGlyph: "LogoExpress"},
		{Name: "expressjs-postgres", DisplayName: "Express.js + Postgres", IconGlyph: "LogoExpress"},
		{Name: "expressjs-prisma", DisplayName: "Express.js + Prisma", IconGlyph: "LogoExpress"},
		{Name: "expressjs", DisplayName: "Express.js", IconGlyph: "LogoExpress"},
		{Name: "fastapi", DisplayName: "FastAPI"},
		{Name: "fastify", DisplayName: "Fastify"},
		{Name: "flask", DisplayName: "Flask", IconGlyph: "LogoFlask"},
		{Name: "ghost", DisplayName: "Ghost", IconGlyph: "LogoGhost"},
		{Name: "gin", DisplayName: "Gin"},
		{Name: "hapi", DisplayName: "Hapi"},
		{Name: "hasura", DisplayName: "Hasura", IconGlyph: "LogoHasura"},
		{Name: "http-nodejs", DisplayName: "Node.js HTTP", IconGlyph: "LogoNode"},
		{Name: "laravel", DisplayName: "Laravel", IconGlyph: "LogoLaravel"},
		{Name: "n8n", DisplayName: "n8n.io"},
		{Name: "next-notion-blog", DisplayName: "Next.js Notion blog", IconGlyph: "LogoNext"},
		{Name: "nextjs-auth-mongo", DisplayName: "Next.js authentication with MongoDB", IconGlyph: "LogoNext"},
		{Name: "nextjs-prisma", DisplayName: "Next.js with Prisma", IconGlyph: "LogoNext"},
		{Name: "nocodb", DisplayName: "NocoDB"},
		{Name: "nuxtjs", DisplayName: "Nuxt.js", IconGlyph: "LogoNuxt"},
		{Name: "rails-blog", DisplayName: "Rails blog", IconGlyph: "LogoRuby"},
		{Name: "rails-starter", DisplayName: "Rails starter", IconGlyph: "LogoRuby"},
		{Name: "rocket", DisplayName: "Rocket"},
		{Name: "rust-wasm", DisplayName: "Rust with Web Assembly", IconGlyph: "LogoRust"},
		{Name: "shiori", DisplayName: "Shiori"},
		{Name: "sinatra", DisplayName: "Sinatra"},
		{Name: "slack-webhook", DisplayName: "Slack webhook", IconGlyph: "LogoSlack"},
		{Name: "starlette", DisplayName: "Starlette"},
		{Name: "strapi", DisplayName: "Strapi", IconGlyph: "LogoStrapi"},
		{Name: "svelte-kit", DisplayName: "Svelte kit", IconGlyph: "LogoSvelte"},
		{Name: "telegram-bot", DisplayName: "Telegram bot", IconGlyph: "LogoTelegram"},
		{Name: "umami", DisplayName: "Umami"},
	}
	for i, t := range templates {
		templates[i].Url = fmt.Sprintf("%s/%s", bucketURL, t.Name)
	}
	return templates
}
