import { Component } from './component.js';

export type CloseButtonConfig = {
	fixed?: boolean,
}

const DEFAULT_CONFIG = {
	fixed: false,
}

export class CloseButton extends Component {
	private readonly config: CloseButtonConfig;

	constructor(parent: HTMLElement, onClick: () => void, config: CloseButtonConfig = {}) {
		super(parent, 'close-button', document.createElement('a'));
		this.config = { ...DEFAULT_CONFIG, ...config };

		if (this.config.fixed)
			this.rootElem.classList.add('position-fixed');

		this.rootElem.setAttribute('href', 'javascript:void(0)');
		this.rootElem.setAttribute('role', 'button');
		this.rootElem.addEventListener('click', () => onClick());

		this.rootElem.insertAdjacentHTML('beforeend', '<i class="fas fa-times fa-2xl ms-1"></i>');
	}
}
