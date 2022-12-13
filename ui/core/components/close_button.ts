import { Component } from './component.js';

export class CloseButton extends Component {
	constructor(parent: HTMLElement, onClick: () => void) {
		super(parent, 'close-button', document.createElement('a'));
		this.rootElem.setAttribute('href', 'javascript:void(0)');
		this.rootElem.setAttribute('role', 'button');
		this.rootElem.addEventListener('click', () => onClick());
		this.rootElem.innerHTML = `
			<span>Close</span><i class="fas fa-times fa-xl ms-1"></i>
		`;
	}
}
