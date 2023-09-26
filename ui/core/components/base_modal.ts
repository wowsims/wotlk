import { CloseButton, CloseButtonConfig } from './close_button';
import { Component } from './component';

import { Modal } from 'bootstrap';

type ModalSize = 'sm' | 'md' | 'lg' | 'xl';

type BaseModalConfig = {
	closeButton?: CloseButtonConfig,
	// Whether or not to add a modal-footer element
	footer?: boolean,
	// Whether or not to add a modal-header element
	header?: boolean,
	// Whether or not to allow modal contents to extend past the screen height.
	// When true, the modal is fixed to the screen height and body contents will scroll.
	scrollContents?: boolean,
	// Specify the size of the modal
	size?: ModalSize,
	// A title for the modal
	title?: string | null,
};

const DEFAULT_CONFIG = {
	closeButton: {},
	footer: false,
	header: true,
	scrollContents: false,
	size: 'lg' as ModalSize,
	title: null,
}

export class BaseModal extends Component {
	readonly modalConfig: BaseModalConfig;

	readonly modal: Modal;
	readonly dialog: HTMLElement;
	readonly header: HTMLElement | undefined;
	readonly body: HTMLElement;
	readonly footer: HTMLElement | undefined;

	constructor(parent: HTMLElement, cssClass: string, config: BaseModalConfig = {}) {
		super(parent, 'modal');
		this.modalConfig = { ...DEFAULT_CONFIG, ...config };

		const modalSizeKlass = this.modalConfig.size && this.modalConfig.size != 'md' ? `modal-${this.modalConfig.size}` : '';

		this.rootElem.classList.add('fade');
		this.rootElem.innerHTML = `
			<div class="modal-dialog ${cssClass} ${modalSizeKlass}">
				<div class="modal-content"></div>
			</div>
		`;

		this.dialog = this.rootElem.querySelector('.modal-dialog') as HTMLElement;

		if (this.modalConfig.scrollContents) {
			this.dialog.classList.add('modal-overflow-scroll');
		}

		const container = this.rootElem.querySelector('.modal-content') as HTMLElement;

		if (this.modalConfig.header) {
			this.header = document.createElement('div');
			this.header.classList.add('modal-header');
			container.appendChild(this.header);

			if (this.modalConfig.title) {
				this.header.insertAdjacentHTML('afterbegin', `<h5 class="modal-title">${this.modalConfig.title}</h5>`);
			}
		}

		this.body = document.createElement('div');
		this.body.classList.add('modal-body');
		container.appendChild(this.body);

		this.addCloseButton();

		if (this.modalConfig.footer) {
			this.footer = document.createElement('div');
			this.footer.classList.add('modal-footer');
			container.appendChild(this.footer);
		}

		this.modal = new Modal(this.rootElem);
		this.open();

		this.rootElem.addEventListener('hidden.bs.modal', (event) => {
			this.rootElem.remove();
			this.dispose();
		})
	}

	private addCloseButton() {
		new CloseButton(this.header ? this.header : this.body, () => this.close(), this.modalConfig.closeButton);
	}

	protected onShow(e: Event) {}

	open() {
		// Hacks for better looking multi modals
		this.rootElem.addEventListener('show.bs.modal', async event => {
			// Prevent the event from bubbling up to parent modals
			event.stopImmediatePropagation();

			// Wait for the backdrop to be injected into the DOM
			const backdrop = await new Promise((resolve) => {
				setTimeout(() => {
					// @ts-ignore
					if (this.modal._backdrop._element)
						// @ts-ignore
						resolve(this.modal._backdrop._element)
				}, 100);
			}) as HTMLElement;
			// Then move it from <body> to the parent element
			this.rootElem.insertAdjacentElement('afterend', backdrop);
			this.onShow(event);
		});

		this.rootElem.addEventListener('hide.bs.modal', (event) => {
			// Prevent the event from bubbling up to parent modals
			event.stopImmediatePropagation();
		});

		this.rootElem.addEventListener('hidden.bs.modal', (event) => {
			// Prevent the event from bubbling up to parent modals
			// Do not use stopImmediatePropagation here. It prevents Bootstrap from removing the modal,
			// leading to other issues
			event.stopPropagation();
		})

		this.modal.show();
	}

	close() {
		this.modal.hide();
	}
}
