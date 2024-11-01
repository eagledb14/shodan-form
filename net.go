package main

import (
	createform "github.com/eagledb14/form-scanner/create-form"
	"github.com/eagledb14/form-scanner/alerts"
	"github.com/eagledb14/form-scanner/types"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	t "github.com/eagledb14/form-scanner/templates"
)

func serv(port string, state *types.State) {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak(), state))
	})

	servCredLeak(app, state)
	servOpenPort(app, state)
	servActor(app, state)
	servEvents(app, state)
	servEvents(app, state)
	servMarkdown(app, state)

	app.Static("/style.css", "./resources/style.css")

	app.Listen(port)
}

func servCredLeak(app *fiber.App, state *types.State) {
	app.Get("/credleak", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.CredLeak(), state))
	})

	app.Post("/credleak", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		form := createform.CredLeak{
			OrgName:    c.FormValue("orgName"),
			FormNumber: c.FormValue("formNumber"),
			VictimOrg:  c.FormValue("victimOrg"),
			Password:   c.FormValue("password"),
			UserPass:   c.FormValue("userPass"),
			AddInfo:    c.FormValue("addInfo"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Name = strings.Clone(form.OrgName)
		state.Title = "Threat Intel Summary"
		state.Tlp = form.Tlp

		return c.Redirect("/preview")
	})
}

func servOpenPort(app *fiber.App, state *types.State) {
	app.Get("/openport", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		if len(state.Events) > 0 {
			return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
		}

		return c.SendString(t.BuildPage(t.OpenPortDownload(), state))
	})

	//makes a new file
	app.Post("/openport", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		form := createform.OpenPort{
			OrgName:    state.Name,
			FormNumber: c.FormValue("formNumber"),
			Threat:     c.FormValue("threat"),
			Summary:    c.FormValue("summary"),
			Body:       c.FormValue("body"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
			Events:     state.Events,
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Title = "Threat Intel Summary"
		state.Tlp = form.Tlp

		return c.Redirect("/preview")
	})

	// clears the state of the selected event
	app.Put("/openport", func(c *fiber.Ctx) error {
		state.Events = []*alerts.Event{}
		state.Name = ""
		return c.SendString(t.BuildPage(t.OpenPortDownload(), state))
	})

	app.Post("/openport/form", func(c *fiber.Ctx) error {
		name := c.FormValue("orgName")
		ips := c.FormValue("ipAddress")

		events := alerts.DownloadIpList(name, ips)
		events = alerts.FilterEvents(events)

		state.Events = events
		state.Name = strings.Clone(name)

		return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
	})

	app.Get("/openport/port", func(c *fiber.Ctx) error {
		return c.SendString(t.BuildPage(t.OpenPortForm(types.Open, state.Name, state.Events), state))
	})

	app.Get("/openport/eol", func(c *fiber.Ctx) error {
		return c.SendString(t.BuildPage(t.OpenPortForm(types.EOL, state.Name, state.Events), state))
	})

	app.Get("/openport/login", func(c *fiber.Ctx) error {
		return c.SendString(t.BuildPage(t.OpenPortForm(types.Login, state.Name, state.Events), state))
	})
}

func servActor(app *fiber.App, state *types.State) {
	app.Get("/actor", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")
		return c.SendString(t.BuildPage(t.Actors(), state))
	})

	app.Post("/actor", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		form := createform.Actor{
			Name:         c.FormValue("name"),
			Alias:        c.FormValue("alias"),
			Date:         c.FormValue("date"),
			Country:      c.FormValue("country"),
			Motivation:   c.FormValue("motivation"),
			Target:       c.FormValue("target"),
			Malware:      c.FormValue("malware"),
			Reporter:     c.FormValue("report"),
			Confidence:   c.FormValue("confidence"),
			Exploits:     c.FormValue("exploits"),
			Summary:      c.FormValue("summary"),
			Capabilities: c.FormValue("capabilities"),
			Detection:    c.FormValue("detection"),
			Ttps:         c.FormValue("ttps"),
			Infra:        c.FormValue("infra"),
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Name = strings.Clone(form.Name)
		state.Title = "Threat Actor Profile"
		state.Tlp = false

		return c.Redirect("/preview")
	})
}

func servEvents(app *fiber.App, state *types.State) {
	app.Get("/event/page/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 {
			index = 0
		}

		state.EventIndex = index

		return c.SendString(t.BuildPage(t.EventList(state.FeedEvents, index), state))
	})

	app.Get("/event/open/:index", func(c *fiber.Ctx) error {
		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.Open, state.EventIndex), state))
	})

	app.Get("/event/eol/:index", func(c *fiber.Ctx) error {
		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.EOL, state.EventIndex), state))
	})

	app.Get("/event/login/:index", func(c *fiber.Ctx) error {
		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.Login, state.EventIndex), state))
	})

	app.Get("/event/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		event := state.GetFeedEvent(index)
		return c.SendString(t.BuildPage(t.EventView(event, index, types.Open, state.EventIndex), state))
	})

	app.Post("/event/:index", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		indexParam := c.Params("index")

		index, err := strconv.Atoi(indexParam)
		if err != nil || index < 0 || index >= len(state.FeedEvents) {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		form := createform.OpenPort{
			OrgName:    state.FeedEvents[index].Name,
			FormNumber: c.FormValue("formNumber"),
			Threat:     c.FormValue("threat"),
			Summary:    c.FormValue("summary"),
			Body:       c.FormValue("body"),
			Reference:  c.FormValue("reference"),
			Tlp:        c.FormValue("tlp") == "amber",
			Events:     []*alerts.Event{state.FeedEvents[index]},
		}
		state.Markdown = form.CreateMarkdown(state)
		state.Name = strings.Clone(form.OrgName)
		state.Title = "Threat Intel Summary"
		state.Tlp = form.Tlp

		return c.Redirect("/preview")
	})
}

func servMarkdown(app *fiber.App, state *types.State) {
	app.Get("/preview", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/html")

		return c.SendString(t.BuildPage(t.MarkdownViewer(state), state))
	})

	app.Post("/preview", func(c *fiber.Ctx) error {
		md := c.FormValue("markdown")
		state.Markdown = md

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/create", func(c *fiber.Ctx) error {
		c.Set("Content-Disposition", "attachment; filename=\""+state.Name+"-"+state.AlertId+".html\"")

		form := createform.CreateHtml(state.Markdown, state.Title, state.Tlp)
		return c.SendString(form)
	})

}
