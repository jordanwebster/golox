package interpreter

type Return struct {
    Value interface {}
}

func (r *Return) Error() string {
    return "Return statement."
}
