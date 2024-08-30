package html

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
	"github.com/waynezhang/homekit-proxy/internal/utils"
)

func RadioGroup(name string, options []string, currentVal string, id string, extraType string) g.Node {
	radios := []g.Node{}

	for _, v := range options {
		radio := Input(
			Class("mr1"+" "),
			Type("radio"),
			Name(name),
			_if(currentVal == v, Checked()),
			Data("id", id),
			Data("value", v),
			Data("type", extraType),
		)
		label := Label(
			Class("mr3"),
			g.Text(v),
		)
		radios = append(radios, radio, label)
	}

	return Span(radios...)
}

func Slider(min string, max string, step string, value string, id string) g.Node {
	return Span(
		g.Rawf("<output id=\"slider-text-%s\">%s</output>", id, utils.TruncateFloat(value)),
		Input(
			Class("ml2"),
			Type("range"),
			Min(min),
			Max(max),
			Step(step),
			Value(value),
			Data("id", id),
			g.Attr("oninput", "document.querySelector('#slider-text-"+id+"').value = this.value"),
		),
	)
}

func _if(cond bool, node g.Node) g.Node {
	if cond {
		return node
	}
	return g.El("")
}
