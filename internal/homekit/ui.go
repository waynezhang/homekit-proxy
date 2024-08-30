package homekit

import (
	"net/http"
	"strconv"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/waynezhang/homekit-proxy/internal/homekit/characteristics"
	"github.com/waynezhang/homekit-proxy/internal/homekit/stat"
	"github.com/waynezhang/homekit-proxy/internal/html"
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
				Link(Rel("stylesheet"), Href("https://unpkg.com/tachyons/css/tachyons.min.css")),
			),
			Body(
				Class("dark-gray pa4 bg-black-025"),
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
	items := []g.Node{
		Class("list pa0"),
	}
	for _, cst := range cstats {
		el := Li(
			Class("mb4 pa3 bg-black-05"),
			Div(
				Class("b"),
				g.Text("# "),
				g.Text(strconv.Itoa(cst.Id)),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Name"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(cst.Name),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Type"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(cst.Type),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Value"),
				),
				Dd(
					Class("pa0 ma0"),
					characteristics.BuildHtmlEl(
						cst.Name,
						cst.Value,
						strconv.Itoa(cst.Id),
						cst,
					),
				),
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
	items := []g.Node{
		Class("list pa0"),
	}
	for _, ast := range astats {
		el := Li(
			Class("mb4 pa3 bg-black-05"),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("ID"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(strconv.Itoa(ast.Id)),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Name"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(ast.Name),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Cmd:"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(ast.Cmd)),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Cron:"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(ast.Cron)),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Tolerance:"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(strconv.Itoa(ast.Tolerance)),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Last Run:"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(ast.LastRun.String()),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Last Error:"),
				),
				Dd(
					Class("pa0 ma0 dib"),
					g.Text(ast.LastError),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Next Run:"),
				),
				Dd(
					Class("pa0 ma0"),
					g.Text(ast.NextRun.String()),
				),
			),
			Dl(
				Class("dib mr4 mt2"),
				Dt(
					Class("b mb1 gray"),
					g.Text("Enabled:"),
				),
				Dd(
					Class("pa0 ma0"),
					html.RadioGroup(
						ast.Name,
						[]string{"true", "false"},
						strconv.FormatBool(ast.Enabled),
						strconv.Itoa(ast.Id),
						characteristics.ExtraTypeAutomation,
					),
				),
			),
		)
		items = append(items, el)
	}
	return Div(
		H2(g.Text("Automations")),
		Ul(items...),
	)
}

const (
	eventScript string = `
	document.querySelectorAll("input[type='radio'][data-type='C']").forEach((input) => {
        input.addEventListener('change', async (e) => {
        	const el = e.target
          	await update("/s/c/" + el.dataset.id, el.dataset.value)
        })
    })
    document.querySelectorAll("input[type='range']").forEach((input) => {
        input.addEventListener('change', async (e) => {
			const el = e.target
          	await update("/s/c/" + el.dataset.id, el.value)
        })
    })
    document.querySelectorAll("input[type='radio'][data-type='A']").forEach((input) => {
           input.addEventListener('change', async (e) => {
           	const el = e.target
             	await update("/s/a/" + el.dataset.id, el.dataset.value)
           })
       })
	`
	utilsScript string = `
	async function update(url, value) {
		return await _fetch(url, "POST", {"value": value})
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
