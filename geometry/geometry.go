package geometry

import (
	"fmt"
	"math" // 👈 Нужен для Pi (π)
)

// Shape (Фигура) — наш интерфейс
type Shape interface {
	Area() (float64, error)
	Perimeter() (float64, error)
}

// --- Прямоугольник ---

type Rectangle struct {
	Width  float64
	Height float64
}

// Area для Rectangle (используем ЗНАЧЕНИЕ `r Rectangle`)
func (r Rectangle) Area() (float64, error) {
	if r.Width <= 0 || r.Height <= 0 {
		return 0, fmt.Errorf("прямоугольник: стороны должны быть положительными (ширина: %f, высота: %f)", r.Width, r.Height)
	}
	return r.Width * r.Height, nil
}

// Perimeter для Rectangle (тоже используем ЗНАЧЕНИЕ `r Rectangle`)
func (r Rectangle) Perimeter() (float64, error) {
	if r.Width <= 0 || r.Height <= 0 {
		return 0, fmt.Errorf("прямоугольник: стороны должны быть положительными (ширина: %f, высота: %f)", r.Width, r.Height)
	}
	return 2 * (r.Width + r.Height), nil
}

// --- Круг ---

type Circle struct {
	Radius float64
}

// Area для Circle (используем ЗНАЧЕНИЕ `c Circle`)
func (c Circle) Area() (float64, error) {
	if c.Radius <= 0 {
		return 0, fmt.Errorf("круг: радиус должен быть положительным (радиус: %f)", c.Radius)
	}
	// Площадь = π * r^2
	return math.Pi * c.Radius * c.Radius, nil
}

// Perimeter для Circle (используем ЗНАЧЕНИЕ `c Circle`)
func (c Circle) Perimeter() (float64, error) {
	if c.Radius <= 0 {
		return 0, fmt.Errorf("круг: радиус должен быть положительным (радиус: %f)", c.Radius)
	}
	// Длина окружности = 2 * π * r
	return 2 * math.Pi * c.Radius, nil
}
