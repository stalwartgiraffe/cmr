package gitlab

import "fmt"

func runErrors(errors []<-chan error) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		isErrorsOpen := true
		for isErrorsOpen {
			isErrorsOpen = false
			for i := range errors {
				if errors[i] == nil {
					continue
				}
				select {
				case err, ok := <-errors[i]:
					if !ok { // this chan is closed
						errors[i] = nil
						continue
					}
					fmt.Println(err)
				default:
					// no error, still open
				}
				isErrorsOpen = true
			}
		}
	}()
	return done
}

func FanIn[T any](fan []<-chan T) <-chan T {
	out := make(chan T, 2*len(fan))
	go func() {
		defer close(out)
		isOpen := true
		for isOpen {
			isOpen = false
			for i := range fan {
				if fan[i] == nil {
					continue
				}
				select {
				case v, ok := <-fan[i]:
					if !ok { // this chan is closed
						fan[i] = nil
						continue
					}
					out <- v
				default:
					// nope, still open
				}
				isOpen = true
			}
		}
	}()
	return out
}

func Transform[S any, D any](
	src <-chan S,
	destCap int,
	makeDest func(s S) D,
) <-chan D {
	dest := make(chan D, destCap)
	go func() {
		defer close(dest)
		for s := range src {
			dest <- makeDest(s)
		}
	}()
	return dest
}

func TransformToOne[S any, D any](
	src <-chan S,
	destCap int,
	makeDest func(s S) []D,
) <-chan D {
	dest := make(chan D, destCap)
	go func() {
		defer close(dest)
		for s := range src {
			for _, d := range makeDest(s) {
				dest <- d
			}
		}
	}()
	return dest
}
