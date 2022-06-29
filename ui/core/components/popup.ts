import { CloseButton } from './close_button.js';
import { Component } from './component.js';

declare var $: any;

export class Popup extends Component {

	constructor(parent: HTMLElement) {
		super(parent, 'popup');

		const computedStyles = window.getComputedStyle(parent);
		this.rootElem.style.setProperty('--main-text-color', computedStyles.getPropertyValue('--main-text-color').trim());
		this.rootElem.style.setProperty('--theme-color-primary', computedStyles.getPropertyValue('--theme-color-primary').trim());
		this.rootElem.style.setProperty('--theme-color-background', computedStyles.getPropertyValue('--theme-color-background').trim());
		this.rootElem.style.setProperty('--theme-color-background-raw', computedStyles.getPropertyValue('--theme-color-background-raw').trim());

		if (parent.closest('.hide-threat-metrics')) {
			this.rootElem.classList.add('hide-threat-metrics');
		}

		$(this.rootElem).bPopup({
			onClose: () => this.rootElem.remove(),
		});
	}

	protected addCloseButton() {
		new CloseButton(this.rootElem, () => this.close());
	}

	close() {
		$(this.rootElem).bPopup().close();
		this.rootElem.remove();
	}
}
