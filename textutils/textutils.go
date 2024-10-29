package textutils

type TextUtilsInterface interface {
}

type TextUtils struct {
}

func NewTextUtils() TextUtilsInterface {
	return &TextUtils{}
}

func (u *TextUtils) GetWordVariants(word string) string {
	return word
}
