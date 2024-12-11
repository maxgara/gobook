package main

type TextField struct {
	f *func(string) //callback function for field change

}

// func (tf *TextField) String () string{
// 	const
// 	return
// }

func NewTextField(f func(string)) *TextField {

	return &TextField{&f}
}
