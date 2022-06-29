import { Component } from '/tbc/core/components/component.js';

export class ResultsViewer extends Component {
	readonly pendingElem: HTMLElement;
	readonly contentElem: HTMLElement;

	constructor(parentElem: HTMLElement) {
		super(parentElem, 'results-viewer');
		this.rootElem.innerHTML = `
      <div class="results-pending">
        <div class="loader"></div>
      </div>
      <div class="results-content">
      </div>
		`;

		this.pendingElem = this.rootElem.getElementsByClassName('results-pending')[0] as HTMLElement;
		this.contentElem = this.rootElem.getElementsByClassName('results-content')[0] as HTMLElement;
		this.hideAll();
	}

	hideAll() {
		this.contentElem.style.display = 'none';
		this.pendingElem.style.display = 'none';
	}

	setPending() {
		this.contentElem.style.display = 'none';
		this.pendingElem.style.display = 'initial';
	}

	setContent(innerHTML: string) {
		this.contentElem.innerHTML = innerHTML;
		this.contentElem.style.display = 'initial';
		this.pendingElem.style.display = 'none';
	}
}
