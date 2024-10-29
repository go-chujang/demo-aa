package response

type body struct {
	code    errorCode
	message *string
	data    interface{}
	err     error
}

const defaultData = "success"

// default data: success
func Body(data ...interface{}) *body {
	b := body{}
	if data != nil {
		return b.Data(data[0])
	}
	return b.Data(defaultData)
}

func Err(err error) *body {
	b := &body{}
	return b.Err(err)
}

func (b *body) Err(err error) *body {
	if err != nil {
		b.err = err
	}
	return b
}

func (b *body) Data(data interface{}) *body {
	b.data = data
	return b
}

func (b *body) Code(code errorCode) *body {
	b.code = code
	return b
}

func (b *body) Message(msg string) *body {
	b.message = &msg
	return b
}
