package main

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/mobile/exp/sensor"
)

const INDICE_JOGADOR = 0
const INDICE_BOLA = 1

type Entidade struct {
	posicao                           fyne.Position
	largura, altura, raio, velx, vely float32
	cor                               color.RGBA
}

type Jogo struct {
	valor                float64
	janela               fyne.Window
	entidades            []Entidade
	tempoDoFrameAnterior time.Time
	bolaEstaParada       bool
}

func gerarTela(j Jogo) *fyne.Container {
	containerTela := container.NewWithoutLayout()
	for _, entidade := range j.entidades {
		if entidade.raio == 0 {
			c := canvas.NewRectangle(entidade.cor)
			c.Resize(fyne.NewSize(entidade.largura, entidade.altura))
			c.Move(entidade.posicao)
			containerTela.Add(c)
		} else {
			c := canvas.NewCircle(entidade.cor)
			c.Resize(fyne.NewSize(entidade.raio, entidade.raio))
			c.Move(entidade.posicao)
			containerTela.Add(c)
		}
	}

	bola := j.entidades[INDICE_BOLA]
	c := canvas.NewText(
		fmt.Sprintf("X: %.2f Y: %.2f VELX: %.2f VELY: %.2f", bola.posicao.X, bola.posicao.Y, bola.velx, bola.vely),
		color.White,
	)
	containerTela.Add(c)
	return containerTela
}

func (j *Jogo) Send(evento interface{}) {
	delta := 60
	j.tempoDoFrameAnterior = time.Now()

	switch evento.(type) {
	case sensor.Event:
		e := evento.(sensor.Event)

		if e.Data[0] > 1 {
			j.entidades[INDICE_JOGADOR].velx -= 1.5 / float32(delta)
			if j.bolaEstaParada {
				j.bolaEstaParada = false
				j.entidades[INDICE_BOLA].vely = -15
			}
		} else if e.Data[0] <= -1 {
			j.entidades[INDICE_JOGADOR].velx += 1.5 / float32(delta)
			if j.bolaEstaParada {
				j.bolaEstaParada = false
				j.entidades[INDICE_BOLA].vely = -15
			}
		} else {
			j.entidades[INDICE_JOGADOR].velx = 0
		}

		j.entidades[INDICE_JOGADOR].posicao.X += j.entidades[INDICE_JOGADOR].velx
		j.entidades[INDICE_BOLA].posicao.X += j.entidades[INDICE_BOLA].velx / float32(delta)
		j.entidades[INDICE_BOLA].posicao.Y += j.entidades[INDICE_BOLA].vely / float32(delta)

		if j.entidades[INDICE_JOGADOR].posicao.X < 0 {
			j.entidades[INDICE_JOGADOR].posicao.X = 0
		}

		if j.entidades[INDICE_JOGADOR].posicao.X+j.entidades[INDICE_JOGADOR].largura > j.janela.Canvas().Size().Width {
			j.entidades[INDICE_JOGADOR].posicao.X = j.janela.Canvas().Size().Width - j.entidades[INDICE_JOGADOR].largura
		}

		if j.entidades[INDICE_BOLA].posicao.X < 0 {
			j.entidades[INDICE_BOLA].posicao.X = 0
			j.entidades[INDICE_BOLA].velx *= -1
		}

		if j.entidades[INDICE_BOLA].posicao.X+j.entidades[INDICE_BOLA].raio > j.janela.Canvas().Size().Width {
			j.entidades[INDICE_BOLA].posicao.X = j.janela.Canvas().Size().Width - j.entidades[INDICE_BOLA].raio
			j.entidades[INDICE_BOLA].velx *= -1
		}

		if j.entidades[INDICE_BOLA].posicao.Y < 0 {
			j.entidades[INDICE_BOLA].posicao.Y = 0
			j.entidades[INDICE_BOLA].vely *= -1
		}

		// bola := j.entidades[INDICE_BOLA]
		// for indice, entidade := range j.entidades {
		// 	if indice == INDICE_BOLA {
		// 		continue
		// 	}

		// 	if (entidade.posicao.X < bola.posicao.X && entidade.posicao.X+entidade.largura > bola.posicao.X) ||
		// 		(entidade.posicao.X < bola.posicao.X+bola.largura && entidade.posicao.X+entidade.largura > bola.posicao.X+bola.largura) ||
		// 		(entidade.posicao.Y > bola.posicao.Y && entidade.posicao.Y+entidade.altura < bola.posicao.Y) ||
		// 		(entidade.posicao.Y > bola.posicao.Y+bola.altura && entidade.posicao.Y+entidade.altura < bola.posicao.Y+bola.altura) {
		// 		j.entidades[INDICE_BOLA].velx *= -1
		// 		j.entidades[INDICE_BOLA].vely *= -1

		// 	}

		// }

		c := canvas.NewRectangle(color.White)
		c.Resize(fyne.NewSize(100, 100))
		j.janela.SetContent(
			gerarTela(*j),
		)

	}

}

func main() {

	a := app.New()
	janela := a.NewWindow("breakout")
	janela.Resize(fyne.NewSize(400, 732))

	janela.SetContent(widget.NewLabel("Olá, mundo!"))

	var j Jogo
	j.janela = janela
	j.valor = 0
	j.tempoDoFrameAnterior = time.Now()
	j.bolaEstaParada = true
	j.entidades = []Entidade{
		Entidade{
			cor:     color.RGBA{255, 255, 255, 255},
			posicao: fyne.NewPos(150, 650),
			largura: 100,
			altura:  10,
			raio:    0,
			velx:    0,
			vely:    0,
		},
		Entidade{
			cor:     color.RGBA{255, 255, 255, 255},
			posicao: fyne.NewPos(190, 640),
			raio:    10,
			largura: 10, // hit box
			altura:  10, // hit box
		},
	}
	sensor.Notify(&j)
	e := sensor.Enable(sensor.Accelerometer, time.Duration(time.Second)/60)
	if e != nil {
		janela.SetContent(widget.NewLabel(e.Error()))
	}

	janela.ShowAndRun()
}
