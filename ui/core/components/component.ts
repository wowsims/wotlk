export abstract class Component {
	readonly rootElem: HTMLElement;

	constructor(parentElem: HTMLElement | null, rootCssClass: string, rootElem?: HTMLElement) {
		this.rootElem = rootElem || document.createElement('div');
		this.rootElem.classList.add(rootCssClass);
		if (parentElem) {
			parentElem.appendChild(this.rootElem);
		}
	}
}
