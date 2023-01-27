import { SimUI } from '../sim_ui.js';
import { EventID, TypedEvent } from '../typed_event.js';

import { Component } from './component.js';

export class LogRunner extends Component {
	constructor(parent: HTMLElement, simUI: SimUI) {
		super(parent, 'log-runner-root');

		const controlBar = document.createElement('div');
		controlBar.classList.add('log-runner-control-bar');
		this.rootElem.appendChild(controlBar);

		const simButton = document.createElement('button');
		simButton.classList.add('log-runner-button', 'btn');
		simButton.textContent = 'Sim 1 Iteration';
		controlBar.appendChild(simButton);

		const logsDiv = document.createElement('div');
		logsDiv.classList.add('log-runner-logs');
		this.rootElem.appendChild(logsDiv);

		simButton.addEventListener('click', async () => simUI.runSimOnce());

		simUI.sim.simResultEmitter.on((eventID, simResult) => {
			const logs = simResult.logs;
			logsDiv.textContent = '';
			logs
				.filter(log => {
					return !log.isCastCompleted();
				})
				.forEach(log => {
					const lineElem = document.createElement('span');
					lineElem.textContent = log.toString();
					logsDiv.appendChild(lineElem);
				});
		});
	}
}
