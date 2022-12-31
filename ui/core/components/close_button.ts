import { Component } from './component.js';

export type CloseButtonConfig = {
	fixed?: boolean,
	text?: boolean,
}

const DEFAULT_CONFIG = {
	fixed: false,
	text: true,
}

export class CloseButton extends Component {
	private readonly config: CloseButtonConfig;

	constructor(parent: HTMLElement, onClick: () => void, config: CloseButtonConfig = {}) {
		super(parent, 'close-button', document.createElement('a'));
		this.config = {...DEFAULT_CONFIG, ...config};

		this.rootElem.classList.add(this.config.fixed ? 'position-fixed' : 'position-absolute');
		this.rootElem.setAttribute('href', 'javascript:void(0)');
		this.rootElem.setAttribute('role', 'button');
		this.rootElem.addEventListener('click', () => onClick());

		if (this.config.text)
			this.rootElem.insertAdjacentHTML('beforeend', '<span>Close</span>');

		this.rootElem.insertAdjacentHTML('beforeend', '<i class="fas fa-times fa-xl ms-1"></i>');
	}
}
