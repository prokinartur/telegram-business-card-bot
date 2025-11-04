package concert

import (
	"fmt" // 👈 Он нужен для fmt.Println, fmt.Printf и fmt.Errorf
)

// Интерфейс, структуры и методы
type Animal interface {
	Speak()
	Running() error
}

type Bird struct {
	Name       string
	CanSpeak   string
	CanRunning bool
}

type Dog struct {
	Name       string
	CanSpeak   string
	CanRunning bool
}

type Cat struct {
	Name       string
	CanSpeak   string
	CanRunning bool
}

func (c *Cat) Speak() {
	fmt.Println("Кошечка", c.Name, "Гоаорит", c.CanSpeak)
}

func (b *Bird) Speak() {
	fmt.Println("Птичка", b.Name, "Говорит", b.CanSpeak)
	if b.Name == "Гоша" {
		panic("AAAAAAAA!!!!!! ГОША КРУШИИИТЬ!!!!!")
	}
}

func (b *Bird) Running() error {
	if !b.CanRunning {
		// Ошибка в логике (копипаста): тут должно быть "Птичка", а не "Кошечка"
		return fmt.Errorf("Птичка %s не умеет бегать.\n", b.Name)
	}
	fmt.Printf("Птичка %s побежала.\n", b.Name) // И тут
	return nil
}

func (c *Cat) Running() error {
	if !c.CanRunning {
		return fmt.Errorf("Кошечка %s не умеет бегать.\n", c.Name)
	}
	fmt.Printf("Кошечка %s побежала.\n", c.Name)
	return nil
}

func (d Dog) Speak() {
	fmt.Println("Собачка", d.Name, "говорит", d.CanSpeak)
}

func (d *Dog) Running() error {
	if !d.CanRunning {
		return fmt.Errorf("Собачка %s не умеет бегать.\n", d.Name)
	}
	fmt.Printf("🐕 Собачка %s побежала!\n", d.Name)
	return nil
}

// Сама функция концерта
func AnimalConcert(animals []Animal) {
	fmt.Println("\n🎪 НАЧИНАЕТСЯ КОНЦЕРТ ЖИВОТНЫХ!")
	for i, animal := range animals {
		fmt.Printf("\n%d. ", i+1)
		var didPanic bool
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("СРОЧНО! Выступление сорвано (паника: %v)\n", r)
					didPanic = true
				}
			}()
			animal.Speak()
		}()
		if didPanic {
			fmt.Println("...животное уводят со сцены.")
			continue
		}
		if err := animal.Running(); err != nil {
			fmt.Println("Ошибка", err)
		}
	}
	fmt.Println("\n🎉 Концерт завершен!")
}
