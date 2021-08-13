package car

type Body struct {
	Element

	FrontWindow *Window
	BackWindow  *Window

	LeftDoor  *Door
	RightDoor *Door
	BackDoor  *Door

	FrontLeftWhile  *While
	FrontRightWhile *While
	BackLeftWhile   *While
	BackRightWhile  *While
}

func (inst *Body) Start() error {
	return nil
}

func (inst *Body) Stop() error {
	return nil
}
