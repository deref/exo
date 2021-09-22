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
		{Name: "deno", DisplayName: "Deno"},
		{Name: "discordjs-typescript", DisplayName: "discord.js Typescript"},
		{Name: "discordjs", DisplayName: "discord.js"},
		{Name: "discordpy", DisplayName: "discord.py"},
		{Name: "djangopy", DisplayName: "DjangoPy"},
		{Name: "elixir-phoenix", DisplayName: "Elixir/Phoenix"},
		{Name: "eris", DisplayName: "Eris"},
		{Name: "expressjs-mongoose", DisplayName: "Express.js + Mongoose"},
		{Name: "expressjs-postgres", DisplayName: "Express.js + Postgres"},
		{Name: "expressjs-prisma", DisplayName: "Express.js + Prisma"},
		{Name: "expressjs", DisplayName: "Express.js"},
		{Name: "fastapi", DisplayName: "FastAPI"},
		{Name: "fastify", DisplayName: "Fastify"},
		{Name: "flask", DisplayName: "Flask"},
		{Name: "ghost", DisplayName: "Ghost"},
		{Name: "gin", DisplayName: "Gin"},
		{Name: "hapi", DisplayName: "Hapi"},
		{Name: "hasura", DisplayName: "Hasura"},
		{Name: "http-nodejs", DisplayName: "Node.js HTTP"},
		{Name: "laravel", DisplayName: "Laravel"},
		{Name: "n8n", DisplayName: "n8n.io"},
		{Name: "next-notion-blog", DisplayName: "Next.js Notion blog"},
		{Name: "nextjs-auth-mongo", DisplayName: "Next.js authentication with MongoDB"},
		{Name: "nextjs-prisma", DisplayName: "Next.js with Prisma"},
		{Name: "nocodb", DisplayName: "NocoDB"},
		{Name: "nuxtjs", DisplayName: "Nuxt.js"},
		{Name: "rails-blog", DisplayName: "Rails blog"},
		{Name: "rails-starter", DisplayName: "Rails starter"},
		{Name: "rocket", DisplayName: "Rocket"},
		{Name: "rust-wasm", DisplayName: "Rust with Web Assembly"},
		{Name: "shiori", DisplayName: "Shiori"},
		{Name: "sinatra", DisplayName: "Sinatra"},
		{Name: "slack-webhook", DisplayName: "Slack webhook"},
		{Name: "starlette", DisplayName: "Starlette"},
		{Name: "strapi", DisplayName: "Strapi"},
		{Name: "svelte-kit", DisplayName: "Svelte kit"},
		{Name: "telegram-bot", DisplayName: "Telegram bot"},
		{Name: "umami", DisplayName: "Umami"},
	}
	for i, t := range templates {
		templates[i].Url = fmt.Sprintf("%s/%s", bucketURL, t.Name)
	}
	return templates
}
