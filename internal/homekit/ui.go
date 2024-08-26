package homekit

import (
	"net/http"
	"strconv"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
)

func (m *HMManager) startUIHandler() {
	m.server.ServeMux().HandleFunc("/ui", func(res http.ResponseWriter, req *http.Request) {
		page(m).Render((res))
	})
}

func page(m *HMManager) g.Node {
	st := m.getAllStat()
	return Doctype(
		HTML(
			Lang("en"),
			Head(
				TitleEl(g.Text(st.Name)),
			),
			Body(
				Div(
					H1(g.Text(st.Name)),
					Div(g.Text(st.Now.String())),
					characteristicsList(st.Characteristics),
					automationsList(st.Automations),
				),
				Script(Type("text/javascript"), g.Raw(utilsScript)),
				Script(Type("text/javascript"), g.Raw(eventScript)),
			),
		),
	)
}

func characteristicsList(cstats []*stat.CharacteristicsStat) g.Node {
	items := []g.Node{}
	for _, cst := range cstats {
		el := Li(
			Dl(
				Dt(g.Text("ID")),
				Dd(g.Text(strconv.Itoa(cst.Id))),
				Dt(g.Text("Name")),
				Dd(g.Text(cst.Name)),
				Dt(g.Text("Type")),
				Dd(g.Text(cst.Type)),
				Dt(g.Text("Value")),
				Dd(characteristics.BuildHtmlEl(
					cst.Name,
					cst.Value,
					strconv.Itoa(cst.Id),
					cst,
				)),
			),
		)
		items = append(items, el)
	}
	return Div(
		H2(g.Text("Characteristics")),
		Ul(items...),
	)
}

func automationsList(astats []*stat.AutomationStat) g.Node {
	items := []g.Node{}
	for _, ast := range astats {
		el := Li(
			Dl(
				Dt(g.Text("ID")),
				Dd(g.Text(strconv.Itoa(ast.Id))),
				Dt(g.Text("Name")),
				Dd(g.Text(ast.Name)),
				Dt(g.Text("Cmd:")),
				Dd(g.Text(ast.Cmd)),
				Dt(g.Text("Cron:")),
				Dd(g.Text(ast.Cron)),
				Dt(g.Text("Margin:")),
				Dd(g.Text(strconv.Itoa(ast.Margin))),
				Dt(g.Text("Last Run:")),
				Dd(g.Text(ast.LastRun.String())),
				Dt(g.Text("Last Error:")),
				Dd(g.Text(ast.LastError)),
				Dt(g.Text("Next Run:")),
				Dd(g.Text(ast.NextRun.String())),
			),
		)
		items = append(items, el)
	}
	return Div(
		H2(g.Text("Characteristics")),
		Ul(items...),
	)
}

const (
	eventScript string = `
	document.querySelectorAll("input[type='radio']").forEach((input) => {
        input.addEventListener('change', async (e) => {
        	const el = e.target
         	const id = el.dataset.id
         	const val = el.dataset.value
          	await update_c(id, val)
        })
    })
    document.querySelectorAll("input[type='range']").forEach((input) => {
           input.addEventListener('change', async (e) => {
           	const el = e.target
            	const id = el.dataset.id
            	const val = el.value
	            await update_c(id, val)
           })
       })
	`
	utilsScript string = `
	async function update_c(id, value) {
		return await _fetch("/s/c/" + id, "POST", {"value": value})
	}
	async function _fetch(url, method, data) {
		const resp = await fetch(url, {
			method: method || "GET",
			headers: { "Content-Type": "application/json" },
			body: data ? JSON.stringify(data) : null,
		})
		if (!resp.ok) {
			const json = await resp.json()
			throw Error(json.message)
		}
		return await resp.json()
	}
	`
)

/*
const (
	initScript string = `
	document.addEventListener('alpine:init', () => {
		Alpine.magic('json', () => async (url, method, data) => {
			return _fetch(url, method, data)
		})
		Alpine.magic('update_c', () => async (id, value) => {
			return _fetch("/s/c/" + id, "POST", {"value": value})
		})
	})

	async function _fetch(url, method, data) {
		const resp = await fetch(url, {
			method: method || "GET",
			headers: { "Content-Type": "application/json" },
			body: data ? JSON.stringify(data) : null,
		})
		if (!resp.ok) {
			const json = await resp.json()
			throw Error(json.message)
		}
		return await resp.json()
	}
	`
)
*/
