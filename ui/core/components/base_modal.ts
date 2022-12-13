import { Modal } from 'bootstrap';
import { CloseButton } from './close_button.js';
import { Component } from './component.js';

export class BaseModal extends Component {
	readonly header: HTMLElement;
	readonly body: HTMLElement;

	readonly modal: Modal;

	constructor(rootElem: HTMLElement, cssClass: string) {
		super(rootElem, 'modal');

		this.rootElem.classList.add('fade');
		this.rootElem.innerHTML = `
			<div class="modal-dialog modal-lg ${cssClass}">
				<div class="modal-content">
				</div>
			</div>
		`;

		const container = this.rootElem.querySelector('.modal-content') as HTMLElement;

		this.header = document.createElement('div');
		this.header.classList.add('modal-header');
		container.appendChild(this.header);

		this.body = document.createElement('div');
		this.body.classList.add('modal-body');
		container.appendChild(this.body);

		this.modal = new Modal(this.rootElem);
		this.open();
	}

	protected addCloseButton() {
		new CloseButton(this.header, () => this.close());
	}

	open() {
		this.modal.show();
	}

	close() {
		this.modal.hide();
	}
}
