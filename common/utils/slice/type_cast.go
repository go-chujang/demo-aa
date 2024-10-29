package slice

func TypeCast[O any, N any](old []O, converter func(O) N) (new []N) {
	if old != nil {
		new = make([]N, 0, len(old))
		for _, o := range old {
			new = append(new, converter(o))
		}
	}
	return new
}

func TypeCastWithIdx[O any, N any](old []O, converter func(int, O) N) (new []N) {
	if old != nil {
		new = make([]N, len(old))
		for i, o := range old {
			new[i] = converter(i, o)
		}
	}
	return new
}

func TypeCastWithErr[O any, N any](old []O, converter func(O) (N, error)) (new []N, err error) {
	if old != nil {
		new = make([]N, 0, len(old))
		for _, o := range old {
			n, err := converter(o)
			if err != nil {
				return nil, err
			}
			new = append(new, n)
		}
	}
	return new, nil
}

func TypeCastWithIdxErr[O any, N any](old []O, converter func(int, O) (N, error)) (new []N, err error) {
	if old != nil {
		new = make([]N, len(old))
		for i, o := range old {
			n, err := converter(i, o)
			if err != nil {
				return nil, err
			}
			new = append(new, n)
		}
	}
	return new, nil
}
