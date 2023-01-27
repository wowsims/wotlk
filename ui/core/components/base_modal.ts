import { CloseButton, CloseButtonConfig } from './close_button';
import { Component } from './component';

import { Modal } from 'bootstrap';

type BaseModalConfig = {
	closeButton?: CloseButtonConfig,
	header?: boolean,
};

const DEFAULT_CONFIG = {
	closeButton: {},
	header: true,
}

export class BaseModal extends Component {
	readonly config: BaseModalConfig;

	readonly modal: Modal;
	readonly header: HTMLElement | undefined;
	readonly body: HTMLElement;

	constructor(parent: HTMLElement, cssClass: string, config: BaseModalConfig = {}) {
		super(parent, 'modal');
		this.config = {...DEFAULT_CONFIG, ...config};

		this.rootElem.classList.add('fade');
		this.rootElem.innerHTML = `
			<div class="modal-dialog modal-lg ${cssClass}">
				<div class="modal-content">
				</div>
			</div>
		`;

		const container = this.rootElem.querySelector('.modal-content') as HTMLElement;

		if (this.config.header) {
			this.header = document.createElement('div');
			this.header.classList.add('modal-header');
			container.appendChild(this.header);
		}

		this.body = document.createElement('div');
		this.body.classList.add('modal-body');
		container.appendChild(this.body);

		this.addCloseButton();

		this.modal = new Modal(this.rootElem);
		this.open();
		
		this.rootElem.addEventListener('hidden.bs.modal', (event) => {
			this.rootElem.remove();
		})
	}

	private addCloseButton() {
		new CloseButton(this.header ? this.header : this.body, () => this.close(), this.config.closeButton);
	}

	open() {
		// Hacks for better looking multi modals
		this.rootElem.addEventListener('show.bs.modal', () => {
			const modals = this.rootElem.parentElement?.querySelectorAll('.modal') as NodeListOf<HTMLElement>;
			const siblingModals = Array.from(modals).filter((e) => e != this.rootElem);
			siblingModals.forEach((element, index) => element.style.zIndex = '1049');
		});

		this.rootElem.addEventListener('hide.bs.modal', () => {
			const modals = this.rootElem.parentElement?.querySelectorAll('.modal') as NodeListOf<HTMLElement>;
			const siblingModals = Array.from(modals).filter((e) => e != this.rootElem);
			const modalIndex = siblingModals.length - 1 < 0 ? 0 : siblingModals.length - 1;
			siblingModals[modalIndex].style.zIndex = '1055';
		});

		this.modal.show();
	}

	close() {
		this.modal.hide();
	}
}
