import { title } from 'process';
import { Component } from './component.js';

export interface ContentBlockTitleConfig {
  text: string,
  classes?: Array<string>,
  tag?: string,
}

export interface ContentBlockConfig {
	bodyClasses?: Array<string>,
  classes?: Array<string>,
  rootElem?: HTMLElement,
  title?: ContentBlockTitleConfig,
}

export class ContentBlock extends Component {
  readonly titleElement: HTMLElement|null;
  readonly bodyElement: HTMLElement;

  readonly config: ContentBlockConfig;

	constructor(parent: HTMLElement, cssClass: string, config: ContentBlockConfig) {
		super(parent, 'content-block', config.rootElem);
    this.config = config;
		this.rootElem.classList.add(cssClass);

		if (config.classes) {
			this.rootElem.classList.add(...config.classes);
    }

    this.titleElement = this.buildHeader();
    this.bodyElement = this.buildBody();
	}

  private buildHeader(): HTMLElement|null {
    if (this.config.title && Object.keys(this.config.title).length) {
      let titleTag = this.config.title.tag || 'h6';
      let headerFragment = document.createElement('fragment');
      headerFragment.innerHTML = `
        <div class="content-block-header">
          <${titleTag} class="content-block-title">${this.config.title.text}</${titleTag}>
        </div>
      `;

      let header = headerFragment.children[0] as HTMLElement;
      
      if (this.config.title.classes) {
        header.classList.add(...this.config.title.classes);
      }

      this.rootElem.appendChild(header);

      return header;
    } else {
      return null;
    }
  }

  private buildBody(): HTMLElement {
    let bodyElem = document.createElement('div');
    bodyElem.classList.add('content-block-body');

    this.rootElem.appendChild(bodyElem);

    return bodyElem;
  }
}
