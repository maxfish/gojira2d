package input

var (
	GameControllerMappings []*GameControllerMapping

	MappingXBox360  GameControllerMapping
	MappingPS4      GameControllerMapping
	MappingKeyboard GameControllerMapping
)

type GameControllerMapping struct {
	nameRegEx string
	buttons   []int
	axes      []int
}

func (g *GameControllerMapping) set(nameRegEx string, buttons []int, axes []int) {
	g.nameRegEx = nameRegEx
	g.buttons = buttons
	g.axes = axes
}

func init() {
	GameControllerMappings = make([]*GameControllerMapping, 0, 10)

	// PS4 controller with USB cable (MacOS)
	MappingPS4.set(".*Wireless Controller.*", []int{1, 2, 0, 3, 8, 12, 9, 10, 11, 4, 5, 14, 16, 17, 15, 13, 6, 7}, []int{0, 1, 2, 3, 4, 5})
	GameControllerMappings = append(GameControllerMappings, &MappingPS4)

	// XBox 360 wired controller (MacOS)
	MappingXBox360.set(".*Xbox 360.*", []int{11, 12, 13, 14, 5, 10, 4, 6, 7, 8, 9, 0, 1, 2, 3}, []int{0, 1, 2, 3, 4, 5})
	GameControllerMappings = append(GameControllerMappings, &MappingXBox360)
}
