package errorplus

import "errors"

// GetDetailStack returns the details registered by errors from this error down. Errors without details are skipped.
func GetDetailStack(err error) []detail {
	if err == nil {
		return []detail{}
	}

	var stack []detail
	for cur := err; cur != nil; cur = errors.Unwrap(cur) {
		if len(cur.Error()) == 0 {
			continue
		}

		switch cur.(type) {
		case *wrapError:
			w := cur.(*wrapError)
			m := map[string]interface{}{}
			for k, v := range w.detail {
				m[k] = v
			}
			if len(w.Error()) > 0 {
				m["error"] = w.Error()
			}
			stack = append(stack, m)
		default:
			stack = append(stack, detail{"error": cur.Error()})
		}
	}

	return stack
}

// GetStatusHint returns the closest status hint set from this error. If no status is found, 0 is returned.
func GetStatusHint(err error) int {
	if err == nil {
		return 0
	}

	for cur := err; cur != nil; cur = errors.Unwrap(cur) {
		if w, ok := cur.(*wrapError); ok && w.statusHint > 0 {
			return w.statusHint
		}
	}

	return 0
}

// GetMessage returns the closest error message to the error. If no message is defined, returns empty string.
func GetMessage(err error) string {
	if err == nil {
		return ""
	}

	for cur := err; cur != nil; cur = errors.Unwrap(cur) {
		if msg := cur.Error(); len(msg) > 0 {
			return msg
		}
	}

	return ""
}

type detail map[string]interface{}
