import { ResultComponent, ResultComponentConfig, SimResultData } from './result_component.js';

const layoutHTML = `
	<div class="log-runner-logs">
	</div>
`
export class LogRunner extends ResultComponent {
	private logsContainer: HTMLElement;

	constructor(config: ResultComponentConfig) {
		config.rootCssClass = 'log-runner-root';
		super(config)

		this.rootElem.innerHTML = layoutHTML
		this.logsContainer = this.rootElem.querySelector('.log-runner-logs') as HTMLElement;
	}

	onSimResult(resultData: SimResultData): void {
		const logs = resultData.result.logs
		this.logsContainer.innerHTML = '';
		logs
			.filter(log => {
				return !log.isCastCompleted();
			})
			.forEach(log => {
				const lineElem = document.createElement('span');
				lineElem.textContent = log.toString();
				this.logsContainer.appendChild(lineElem);
			});
	}
}
