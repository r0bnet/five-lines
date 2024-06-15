package main

type Canvas struct {
	Width  int
	Height int
}

func GetElementById(id string) Canvas {
	return Canvas{}
}

type CanvasRenderingContext2D struct {
	FillStyle string
}

func (c *Canvas) GetContext(context string) CanvasRenderingContext2D {
	return CanvasRenderingContext2D{
		FillStyle: "#ffffff",
	}
}

func (c *CanvasRenderingContext2D) ClearRect(x, y, width, height int) {}

func (c *CanvasRenderingContext2D) FillRect(x, y, width, height int) {}
