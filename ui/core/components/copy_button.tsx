import { Tooltip } from 'bootstrap';
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import { element, ref } from 'tsx-vanilla';

import { Component } from "./component";

export interface CopyButtonConfig {
	getContent: () => string,
	extraCssClasses?: string[],
	text?: string,
	tooltip?: string,
}

export class CopyButton extends Component {
  private readonly config: CopyButtonConfig;

  constructor(parent: HTMLElement, config: CopyButtonConfig) {
		const btnRef = ref<HTMLButtonElement>();
		const buttonElem = (
			<button
				className={`btn ${config.extraCssClasses?.join(' ') ?? ''}`}
				ref={btnRef}
			>
				<i className="fas fa-copy me-1" />{config.text ?? 'Copy to Clipboard'}
			</button>
		)

    super(parent, 'copy-button', buttonElem as HTMLElement);
    this.config = config;

		const button = btnRef.value!
		let clicked = false
		button.addEventListener('click', _event => {
			if (clicked) return

			const data = this.config.getContent()
			if (navigator.clipboard == undefined) {
				alert(data);
			} else {
				clicked = true
				navigator.clipboard.writeText(data);
				const originalContent = button.innerHTML;
				button.innerHTML = '<i class="fas fa-check me-1"></i>Copied';
				setTimeout(() => {
					button.innerHTML = originalContent;
					clicked = false
				}, 1500);
			}
		});

		if (config.tooltip) {
			Tooltip.getOrCreateInstance(button, {title: config.tooltip});
		}
  }
};
