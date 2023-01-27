import { CloseButton } from './close_button.js';
import { Component } from './component.js';

declare var $: any;

export class Popup extends Component {
	constructor(parent: HTMLElement) {
		super(parent, 'popup');

		if (parent.closest('.hide-damage-metrics')) {
			this.rootElem.classList.add('hide-damage-metrics');
		}
		if (parent.closest('.hide-threat-metrics')) {
			this.rootElem.classList.add('hide-threat-metrics');
		}
		if (parent.closest('.hide-healing-metrics')) {
			this.rootElem.classList.add('hide-healing-metrics');
		}

		$(this.rootElem).bPopup({
			onClose: () => {
				this.rootElem.remove();
				this.dispose();
			},
		});
	}

	protected addCloseButton() {
		new CloseButton(this.rootElem, () => this.close());
	}

	close() {
		$(this.rootElem).bPopup().close();
		this.rootElem.remove();
		this.dispose();
	}
}
