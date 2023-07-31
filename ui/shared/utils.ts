export function isDescendant(el: HTMLElement, parent: HTMLElement): boolean {
	let isChild = false

	if (el === parent) { //is this the element itself?
		isChild = true
	}

	while (el.parentNode && (el = el.parentNode as HTMLElement)) {
		if (el == parent) {
			isChild = true
		}
	}

	return isChild
}
